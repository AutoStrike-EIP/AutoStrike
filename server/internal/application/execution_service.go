package application

import (
	"context"
	"fmt"
	"time"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"
	"autostrike/internal/domain/service"

	"github.com/google/uuid"
)

// ExecutionService handles execution-related business logic
type ExecutionService struct {
	resultRepo    repository.ResultRepository
	scenarioRepo  repository.ScenarioRepository
	techniqueRepo repository.TechniqueRepository
	agentRepo     repository.AgentRepository
	orchestrator  *service.AttackOrchestrator
	calculator    *service.ScoreCalculator
}

// NewExecutionService creates a new execution service
func NewExecutionService(
	resultRepo repository.ResultRepository,
	scenarioRepo repository.ScenarioRepository,
	techniqueRepo repository.TechniqueRepository,
	agentRepo repository.AgentRepository,
	orchestrator *service.AttackOrchestrator,
	calculator *service.ScoreCalculator,
) *ExecutionService {
	return &ExecutionService{
		resultRepo:    resultRepo,
		scenarioRepo:  scenarioRepo,
		techniqueRepo: techniqueRepo,
		agentRepo:     agentRepo,
		orchestrator:  orchestrator,
		calculator:    calculator,
	}
}

// TaskDispatchInfo contains information needed to dispatch a task to an agent
type TaskDispatchInfo struct {
	ResultID    string
	AgentPaw    string
	TechniqueID string
	Command     string
	Executor    string
	Timeout     int
	Cleanup     string
}

// ExecutionWithTasks contains the execution and tasks to dispatch
type ExecutionWithTasks struct {
	Execution *entity.Execution
	Tasks     []TaskDispatchInfo
}

// StartExecution starts a new scenario execution
func (s *ExecutionService) StartExecution(
	ctx context.Context,
	scenarioID string,
	agentPaws []string,
	safeMode bool,
) (*ExecutionWithTasks, error) {
	// Load scenario
	scenario, err := s.scenarioRepo.FindByID(ctx, scenarioID)
	if err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}

	// Load agents in batch (single query instead of N+1)
	agents, err := s.agentRepo.FindByPaws(ctx, agentPaws)
	if err != nil {
		return nil, fmt.Errorf("failed to load agents: %w", err)
	}

	// Build a map for quick lookup and validate all agents exist
	agentMap := make(map[string]*entity.Agent, len(agents))
	for _, agent := range agents {
		agentMap[agent.Paw] = agent
	}

	// Verify all requested agents were found and are online
	for _, paw := range agentPaws {
		agent, found := agentMap[paw]
		if !found {
			return nil, fmt.Errorf("agent %s not found", paw)
		}
		if agent.Status != entity.AgentOnline {
			return nil, fmt.Errorf("agent %s is not online", paw)
		}
	}

	// Create execution plan
	plan, err := s.orchestrator.PlanExecution(ctx, scenario, agents, safeMode)
	if err != nil {
		return nil, fmt.Errorf("failed to plan execution: %w", err)
	}

	// Create execution record
	execution := &entity.Execution{
		ID:         uuid.New().String(),
		ScenarioID: scenarioID,
		AgentPaws:  agentPaws,
		Status:     entity.ExecutionRunning,
		StartedAt:  time.Now(),
		SafeMode:   safeMode,
	}

	if err := s.resultRepo.CreateExecution(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Create pending results for each task and collect dispatch info
	tasks := make([]TaskDispatchInfo, 0, len(plan.Tasks))
	for _, task := range plan.Tasks {
		result := &entity.ExecutionResult{
			ID:          uuid.New().String(),
			ExecutionID: execution.ID,
			TechniqueID: task.TechniqueID,
			AgentPaw:    task.AgentPaw,
			Status:      entity.StatusPending,
			StartedAt:   time.Now(),
		}

		if err := s.resultRepo.CreateResult(ctx, result); err != nil {
			return nil, fmt.Errorf("failed to create result: %w", err)
		}

		// Get the technique to determine the executor
		technique, _ := s.techniqueRepo.FindByID(ctx, task.TechniqueID)
		executor := "sh" // default
		if technique != nil {
			agent := agentMap[task.AgentPaw]
			if agent != nil {
				if exec := technique.GetExecutorForPlatform(agent.Platform, agent.Executors); exec != nil {
					executor = exec.Type
				}
			}
		}

		tasks = append(tasks, TaskDispatchInfo{
			ResultID:    result.ID,
			AgentPaw:    task.AgentPaw,
			TechniqueID: task.TechniqueID,
			Command:     task.Command,
			Executor:    executor,
			Timeout:     task.Timeout,
			Cleanup:     task.Cleanup,
		})
	}

	return &ExecutionWithTasks{
		Execution: execution,
		Tasks:     tasks,
	}, nil
}

// UpdateResult updates an execution result
func (s *ExecutionService) UpdateResult(
	ctx context.Context,
	resultID string,
	status entity.ResultStatus,
	output string,
	detected bool,
) error {
	result, err := s.resultRepo.FindResultByID(ctx, resultID)
	if err != nil {
		return err
	}

	now := time.Now()
	result.Status = status
	result.Output = output
	result.Detected = detected
	result.CompletedAt = &now
	return s.resultRepo.UpdateResult(ctx, result)
}

// UpdateResultByID updates a result by its ID with exit code
func (s *ExecutionService) UpdateResultByID(
	ctx context.Context,
	resultID string,
	status entity.ResultStatus,
	output string,
	exitCode int,
) error {
	result, err := s.resultRepo.FindResultByID(ctx, resultID)
	if err != nil {
		return fmt.Errorf("result not found: %w", err)
	}

	executionID := result.ExecutionID

	now := time.Now()
	result.Status = status
	result.Output = output
	result.ExitCode = exitCode
	result.CompletedAt = &now

	if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
		return err
	}

	// Check if all results are completed and auto-complete execution
	return s.checkAndCompleteExecution(ctx, executionID)
}

// checkAndCompleteExecution checks if all results are done and completes the execution
func (s *ExecutionService) checkAndCompleteExecution(ctx context.Context, executionID string) error {
	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return nil // Don't fail the result update if we can't check
	}

	// Check if all results are completed (not pending or running)
	allDone := true
	for _, r := range results {
		if r.Status == entity.StatusPending || r.Status == entity.StatusRunning {
			allDone = false
			break
		}
	}

	if allDone && len(results) > 0 {
		// All tasks completed, mark execution as completed
		return s.CompleteExecution(ctx, executionID)
	}

	return nil
}

// CompleteExecution marks an execution as completed and calculates score
func (s *ExecutionService) CompleteExecution(ctx context.Context, executionID string) error {
	execution, err := s.resultRepo.FindExecutionByID(ctx, executionID)
	if err != nil {
		return err
	}

	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return err
	}

	// Calculate score
	now := time.Now()
	score := s.calculator.CalculateScore(results)
	execution.Score = score
	execution.Status = entity.ExecutionCompleted
	execution.CompletedAt = &now

	return s.resultRepo.UpdateExecution(ctx, execution)
}

// GetExecution retrieves an execution by ID
func (s *ExecutionService) GetExecution(ctx context.Context, id string) (*entity.Execution, error) {
	return s.resultRepo.FindExecutionByID(ctx, id)
}

// GetExecutionResults retrieves results for an execution
func (s *ExecutionService) GetExecutionResults(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return s.resultRepo.FindResultsByExecution(ctx, executionID)
}

// GetRecentExecutions retrieves recent executions
func (s *ExecutionService) GetRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	return s.resultRepo.FindRecentExecutions(ctx, limit)
}

// CancelExecution stops a running execution
func (s *ExecutionService) CancelExecution(ctx context.Context, executionID string) error {
	execution, err := s.resultRepo.FindExecutionByID(ctx, executionID)
	if err != nil {
		return fmt.Errorf("execution not found: %w", err)
	}

	// Only running or pending executions can be cancelled
	if execution.Status != entity.ExecutionRunning && execution.Status != entity.ExecutionPending {
		return fmt.Errorf("execution cannot be cancelled: status is %s", execution.Status)
	}

	// Update all pending results to skipped
	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return fmt.Errorf("failed to get results: %w", err)
	}

	now := time.Now()
	for _, result := range results {
		if result.Status == entity.StatusPending || result.Status == entity.StatusRunning {
			result.Status = entity.StatusSkipped
			result.CompletedAt = &now
			if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
				return fmt.Errorf("failed to update result %s: %w", result.ID, err)
			}
		}
	}

	// Mark execution as cancelled
	execution.Status = entity.ExecutionCancelled
	execution.CompletedAt = &now

	return s.resultRepo.UpdateExecution(ctx, execution)
}
