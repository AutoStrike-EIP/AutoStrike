package service

import (
	"context"
	"fmt"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AttackOrchestrator coordinates attack execution across agents
type AttackOrchestrator struct {
	agentRepo     repository.AgentRepository
	techniqueRepo repository.TechniqueRepository
	validator     *TechniqueValidator
	logger        *zap.Logger
}

// NewAttackOrchestrator creates a new orchestrator instance
func NewAttackOrchestrator(
	agentRepo repository.AgentRepository,
	techniqueRepo repository.TechniqueRepository,
	validator *TechniqueValidator,
	logger *zap.Logger,
) *AttackOrchestrator {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &AttackOrchestrator{
		agentRepo:     agentRepo,
		techniqueRepo: techniqueRepo,
		validator:     validator,
		logger:        logger,
	}
}

// ExecutionPlan represents a planned execution
type ExecutionPlan struct {
	ID    string
	Tasks []PlannedTask
}

// PlannedTask represents a single task in the execution plan
type PlannedTask struct {
	TechniqueID string
	AgentPaw    string
	Phase       string
	Order       int
	Command     string
	Cleanup     string
	Timeout     int
}

// PlanExecution creates an execution plan for a scenario
func (o *AttackOrchestrator) PlanExecution(
	ctx context.Context,
	scenario *entity.Scenario,
	targetAgents []*entity.Agent,
	safeMode bool,
) (*ExecutionPlan, error) {
	plan := &ExecutionPlan{
		ID:    uuid.New().String(),
		Tasks: make([]PlannedTask, 0),
	}

	taskOrder := 0

	for _, phase := range scenario.Phases {
		for _, techID := range phase.Techniques {
			technique, err := o.techniqueRepo.FindByID(ctx, techID)
			if err != nil {
				o.logger.Warn("Skipping technique: not found in repository", zap.String("technique_id", techID))
				continue
			}

			// Skip unsafe techniques in safe mode
			if safeMode && !technique.IsSafe {
				o.logger.Info("Skipping unsafe technique in safe mode", zap.String("technique_id", techID))
				continue
			}

			// Find compatible agents
			for _, agent := range targetAgents {
				if !agent.IsCompatible(technique) {
					continue
				}

				// Get the appropriate executor
				executor := technique.GetExecutorForPlatform(agent.Platform, agent.Executors)
				if executor == nil {
					continue
				}

				plan.Tasks = append(plan.Tasks, PlannedTask{
					TechniqueID: techID,
					AgentPaw:    agent.Paw,
					Phase:       phase.Name,
					Order:       taskOrder,
					Command:     executor.Command,
					Cleanup:     executor.Cleanup,
					Timeout:     executor.Timeout,
				})

				taskOrder++
			}
		}
	}

	if len(plan.Tasks) == 0 {
		return nil, fmt.Errorf("no executable tasks for the given scenario and agents")
	}

	return plan, nil
}

// ValidatePlan validates an execution plan
func (o *AttackOrchestrator) ValidatePlan(ctx context.Context, plan *ExecutionPlan) error {
	for _, task := range plan.Tasks {
		agent, err := o.agentRepo.FindByPaw(ctx, task.AgentPaw)
		if err != nil {
			return fmt.Errorf("agent %s not found", task.AgentPaw)
		}

		if agent.Status != entity.AgentOnline {
			return fmt.Errorf("agent %s is not online", task.AgentPaw)
		}

		technique, err := o.techniqueRepo.FindByID(ctx, task.TechniqueID)
		if err != nil {
			return fmt.Errorf("technique %s not found", task.TechniqueID)
		}

		if !agent.IsCompatible(technique) {
			return fmt.Errorf("agent %s is not compatible with technique %s", task.AgentPaw, task.TechniqueID)
		}
	}

	return nil
}
