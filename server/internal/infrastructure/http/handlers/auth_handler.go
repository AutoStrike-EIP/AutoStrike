package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service        *application.AuthService
	tokenBlacklist *application.TokenBlacklist
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// NewAuthHandlerWithBlacklist creates a new auth handler with token revocation support
func NewAuthHandlerWithBlacklist(service *application.AuthService, blacklist *application.TokenBlacklist) *AuthHandler {
	return &AuthHandler{service: service, tokenBlacklist: blacklist}
}

// RegisterRoutes registers auth routes (public - no auth middleware)
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}

// RegisterRoutesWithRateLimit registers auth routes with rate limiting
func (h *AuthHandler) RegisterRoutesWithRateLimit(r *gin.Engine, loginLimiter, refreshLimiter *middleware.RateLimiter) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", middleware.RateLimitMiddleware(loginLimiter), h.Login)
		auth.POST("/refresh", middleware.RateLimitMiddleware(refreshLimiter), h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}

// RegisterProtectedRoutes registers routes that require authentication
func (h *AuthHandler) RegisterProtectedRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.GET("/me", h.Me)
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,max=72"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshRequest represents the refresh token request body
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, application.ErrInvalidToken) || errors.Is(err, application.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		if errors.Is(err, application.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token refresh failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handles user logout by revoking the current access token
func (h *AuthHandler) Logout(c *gin.Context) {
	if h.tokenBlacklist != nil {
		h.revokeTokenFromHeader(c.GetHeader("Authorization"))
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// revokeTokenFromHeader extracts a Bearer token and adds it to the blacklist
func (h *AuthHandler) revokeTokenFromHeader(authHeader string) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return
	}

	tokenString := parts[1]
	expiry := getTokenExpiry(tokenString)
	h.tokenBlacklist.Revoke(tokenString, expiry)
}

// getTokenExpiry parses the token expiry without signature validation
func getTokenExpiry(tokenString string) time.Time {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Now().Add(24 * time.Hour)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Now().Add(24 * time.Hour)
	}

	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return time.Now().Add(24 * time.Hour)
	}

	return exp.Time
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	// User ID is set by the auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), userIDStr)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
