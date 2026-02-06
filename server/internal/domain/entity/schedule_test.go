package entity

import (
	"testing"
	"time"
)

func TestScheduleStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		got      ScheduleStatus
		expected string
	}{
		{"Active", ScheduleStatusActive, "active"},
		{"Paused", ScheduleStatusPaused, "paused"},
		{"Disabled", ScheduleStatusDisabled, "disabled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("ScheduleStatus %s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestScheduleFrequency_Constants(t *testing.T) {
	tests := []struct {
		name     string
		got      ScheduleFrequency
		expected string
	}{
		{"Once", FrequencyOnce, "once"},
		{"Hourly", FrequencyHourly, "hourly"},
		{"Daily", FrequencyDaily, "daily"},
		{"Weekly", FrequencyWeekly, "weekly"},
		{"Monthly", FrequencyMonthly, "monthly"},
		{"Cron", FrequencyCron, "cron"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("ScheduleFrequency %s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestSchedule_CalculateNextRun(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		schedule  *Schedule
		from      time.Time
		wantNil   bool
		checkFunc func(t *testing.T, result *time.Time)
	}{
		{
			name: "Hourly - calculates next hour",
			schedule: &Schedule{
				Frequency: FrequencyHourly,
				Status:    ScheduleStatusActive,
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				expected := baseTime.Add(time.Hour)
				if !result.Equal(expected) {
					t.Errorf("got %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Daily - calculates next day",
			schedule: &Schedule{
				Frequency: FrequencyDaily,
				Status:    ScheduleStatusActive,
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				expected := baseTime.AddDate(0, 0, 1)
				if !result.Equal(expected) {
					t.Errorf("got %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Weekly - calculates next week",
			schedule: &Schedule{
				Frequency: FrequencyWeekly,
				Status:    ScheduleStatusActive,
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				expected := baseTime.AddDate(0, 0, 7)
				if !result.Equal(expected) {
					t.Errorf("got %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Monthly - calculates next month",
			schedule: &Schedule{
				Frequency: FrequencyMonthly,
				Status:    ScheduleStatusActive,
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				expected := baseTime.AddDate(0, 1, 0)
				if !result.Equal(expected) {
					t.Errorf("got %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Once with no last run and no NextRunAt - returns from time",
			schedule: &Schedule{
				Frequency: FrequencyOnce,
				Status:    ScheduleStatusActive,
				LastRunAt: nil,
				NextRunAt: nil,
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				if !result.Equal(baseTime) {
					t.Errorf("got %v, want %v", result, baseTime)
				}
			},
		},
		{
			name: "Once with NextRunAt set - returns NextRunAt",
			schedule: &Schedule{
				Frequency: FrequencyOnce,
				Status:    ScheduleStatusActive,
				LastRunAt: nil,
				NextRunAt: func() *time.Time { t := baseTime.Add(2 * time.Hour); return &t }(),
			},
			from:    baseTime,
			wantNil: false,
			checkFunc: func(t *testing.T, result *time.Time) {
				expected := baseTime.Add(2 * time.Hour)
				if !result.Equal(expected) {
					t.Errorf("got %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Once with last run - returns nil",
			schedule: &Schedule{
				Frequency: FrequencyOnce,
				Status:    ScheduleStatusActive,
				LastRunAt: &baseTime,
			},
			from:    baseTime,
			wantNil: true,
		},
		{
			name: "Paused status - returns nil",
			schedule: &Schedule{
				Frequency: FrequencyDaily,
				Status:    ScheduleStatusPaused,
			},
			from:    baseTime,
			wantNil: true,
		},
		{
			name: "Disabled status - returns nil",
			schedule: &Schedule{
				Frequency: FrequencyDaily,
				Status:    ScheduleStatusDisabled,
			},
			from:    baseTime,
			wantNil: true,
		},
		{
			name: "Cron frequency - returns next run from cron expression",
			schedule: &Schedule{
				Frequency: FrequencyCron,
				Status:    ScheduleStatusActive,
				CronExpr:  "0 * * * *", // Every hour at minute 0
			},
			from: baseTime,
			checkFunc: func(t *testing.T, result *time.Time) {
				// Next run should be at the next hour mark
				expected := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
				if !result.Equal(expected) {
					t.Errorf("CalculateNextRun() = %v, want %v", result, expected)
				}
			},
		},
		{
			name: "Cron frequency without expression - returns nil",
			schedule: &Schedule{
				Frequency: FrequencyCron,
				Status:    ScheduleStatusActive,
				CronExpr:  "",
			},
			from:    baseTime,
			wantNil: true,
		},
		{
			name: "Cron frequency with invalid expression - returns nil",
			schedule: &Schedule{
				Frequency: FrequencyCron,
				Status:    ScheduleStatusActive,
				CronExpr:  "invalid cron",
			},
			from:    baseTime,
			wantNil: true,
		},
		{
			name: "Unknown frequency - returns nil",
			schedule: &Schedule{
				Frequency: ScheduleFrequency("unknown"),
				Status:    ScheduleStatusActive,
			},
			from:    baseTime,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schedule.CalculateNextRun(tt.from)

			if tt.wantNil {
				if result != nil {
					t.Errorf("CalculateNextRun() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Fatal("CalculateNextRun() = nil, want non-nil")
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, result)
			}
		})
	}
}

func TestSchedule_IsReadyToRun(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)
	futureTime := now.Add(1 * time.Hour)

	tests := []struct {
		name     string
		schedule *Schedule
		now      time.Time
		want     bool
	}{
		{
			name: "Active with past NextRunAt - ready",
			schedule: &Schedule{
				Status:    ScheduleStatusActive,
				NextRunAt: &pastTime,
			},
			now:  now,
			want: true,
		},
		{
			name: "Active with NextRunAt equal to now - ready",
			schedule: &Schedule{
				Status:    ScheduleStatusActive,
				NextRunAt: &now,
			},
			now:  now,
			want: true,
		},
		{
			name: "Active with future NextRunAt - not ready",
			schedule: &Schedule{
				Status:    ScheduleStatusActive,
				NextRunAt: &futureTime,
			},
			now:  now,
			want: false,
		},
		{
			name: "Active with nil NextRunAt - not ready",
			schedule: &Schedule{
				Status:    ScheduleStatusActive,
				NextRunAt: nil,
			},
			now:  now,
			want: false,
		},
		{
			name: "Paused with past NextRunAt - not ready",
			schedule: &Schedule{
				Status:    ScheduleStatusPaused,
				NextRunAt: &pastTime,
			},
			now:  now,
			want: false,
		},
		{
			name: "Disabled with past NextRunAt - not ready",
			schedule: &Schedule{
				Status:    ScheduleStatusDisabled,
				NextRunAt: &pastTime,
			},
			now:  now,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schedule.IsReadyToRun(tt.now); got != tt.want {
				t.Errorf("IsReadyToRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_Struct(t *testing.T) {
	now := time.Now()
	nextRun := now.Add(24 * time.Hour)
	lastRun := now.Add(-24 * time.Hour)

	schedule := &Schedule{
		ID:          "sched-1",
		Name:        "Test Schedule",
		Description: "A test schedule",
		ScenarioID:  "scenario-1",
		AgentPaw:    "agent-1",
		Frequency:   FrequencyDaily,
		CronExpr:    "",
		SafeMode:    true,
		Status:      ScheduleStatusActive,
		NextRunAt:   &nextRun,
		LastRunAt:   &lastRun,
		LastRunID:   "exec-1",
		CreatedBy:   "user-1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if schedule.ID != "sched-1" {
		t.Errorf("ID = %q, want %q", schedule.ID, "sched-1")
	}
	if schedule.Name != "Test Schedule" {
		t.Errorf("Name = %q, want %q", schedule.Name, "Test Schedule")
	}
	if schedule.Frequency != FrequencyDaily {
		t.Errorf("Frequency = %q, want %q", schedule.Frequency, FrequencyDaily)
	}
	if !schedule.SafeMode {
		t.Error("SafeMode should be true")
	}
	if schedule.Status != ScheduleStatusActive {
		t.Errorf("Status = %q, want %q", schedule.Status, ScheduleStatusActive)
	}
}

func TestScheduleRun_Struct(t *testing.T) {
	startedAt := time.Now()
	completedAt := startedAt.Add(5 * time.Minute)

	run := &ScheduleRun{
		ID:          "run-1",
		ScheduleID:  "sched-1",
		ExecutionID: "exec-1",
		StartedAt:   startedAt,
		CompletedAt: &completedAt,
		Status:      "completed",
		Error:       "",
	}

	if run.ID != "run-1" {
		t.Errorf("ID = %q, want %q", run.ID, "run-1")
	}
	if run.ScheduleID != "sched-1" {
		t.Errorf("ScheduleID = %q, want %q", run.ScheduleID, "sched-1")
	}
	if run.Status != "completed" {
		t.Errorf("Status = %q, want %q", run.Status, "completed")
	}
	if run.Error != "" {
		t.Errorf("Error = %q, want empty", run.Error)
	}
}

func TestScheduleRun_Failed(t *testing.T) {
	startedAt := time.Now()
	completedAt := startedAt.Add(1 * time.Second)

	run := &ScheduleRun{
		ID:          "run-2",
		ScheduleID:  "sched-1",
		ExecutionID: "", // No execution ID for failed runs
		StartedAt:   startedAt,
		CompletedAt: &completedAt,
		Status:      "failed",
		Error:       "connection timeout",
	}

	if run.Status != "failed" {
		t.Errorf("Status = %q, want %q", run.Status, "failed")
	}
	if run.Error != "connection timeout" {
		t.Errorf("Error = %q, want %q", run.Error, "connection timeout")
	}
	if run.ExecutionID != "" {
		t.Errorf("ExecutionID = %q, want empty", run.ExecutionID)
	}
}

func TestValidateCronExpr_Valid(t *testing.T) {
	tests := []string{
		"0 0 * * *",       // Every day at midnight
		"*/5 * * * *",     // Every 5 minutes
		"0 9 * * 1-5",     // Weekdays at 9am
		"30 2 1 * *",      // 1st of month at 2:30
	}
	for _, expr := range tests {
		if err := ValidateCronExpr(expr); err != nil {
			t.Errorf("ValidateCronExpr(%q) = %v, want nil", expr, err)
		}
	}
}

func TestValidateCronExpr_Invalid(t *testing.T) {
	tests := []string{
		"not a cron",
		"* * *",           // Too few fields
		"60 * * * *",      // Invalid minute
		"* 25 * * *",      // Invalid hour
	}
	for _, expr := range tests {
		if err := ValidateCronExpr(expr); err == nil {
			t.Errorf("ValidateCronExpr(%q) = nil, want error", expr)
		}
	}
}

func TestValidateCronExpr_Empty(t *testing.T) {
	if err := ValidateCronExpr(""); err != nil {
		t.Errorf("ValidateCronExpr(\"\") = %v, want nil", err)
	}
}

func TestCalculateNextRun_CronEmptyExpr(t *testing.T) {
	s := &Schedule{
		Frequency: FrequencyCron,
		CronExpr:  "",
		Status:    ScheduleStatusActive,
	}
	result := s.CalculateNextRun(time.Now())
	if result != nil {
		t.Error("CalculateNextRun should return nil for empty cron expression")
	}
}

func TestCalculateNextRun_CronInvalidExpr(t *testing.T) {
	s := &Schedule{
		Frequency: FrequencyCron,
		CronExpr:  "invalid",
		Status:    ScheduleStatusActive,
	}
	result := s.CalculateNextRun(time.Now())
	if result != nil {
		t.Error("CalculateNextRun should return nil for invalid cron expression")
	}
}

func TestCalculateNextRun_CronValidExpr(t *testing.T) {
	s := &Schedule{
		Frequency: FrequencyCron,
		CronExpr:  "0 0 * * *",
		Status:    ScheduleStatusActive,
	}
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	result := s.CalculateNextRun(now)
	if result == nil {
		t.Fatal("CalculateNextRun should return next run time for valid cron")
	}
	if !result.After(now) {
		t.Error("Next run should be after 'now'")
	}
}

func TestCalculateNextRun_UnknownFrequency(t *testing.T) {
	s := &Schedule{
		Frequency: "unknown",
		Status:    ScheduleStatusActive,
	}
	result := s.CalculateNextRun(time.Now())
	if result != nil {
		t.Error("CalculateNextRun should return nil for unknown frequency")
	}
}

func TestCalculateNextRun_OnceWithExistingNextRun(t *testing.T) {
	nextRun := time.Now().Add(time.Hour)
	s := &Schedule{
		Frequency: FrequencyOnce,
		Status:    ScheduleStatusActive,
		NextRunAt: &nextRun,
	}
	result := s.CalculateNextRun(time.Now())
	if result == nil {
		t.Fatal("CalculateNextRun should return NextRunAt for once frequency")
	}
	if !result.Equal(nextRun) {
		t.Error("Should return the existing NextRunAt")
	}
}

func TestCalculateNextRun_OnceAlreadyRun(t *testing.T) {
	lastRun := time.Now()
	s := &Schedule{
		Frequency: FrequencyOnce,
		Status:    ScheduleStatusActive,
		LastRunAt: &lastRun,
	}
	result := s.CalculateNextRun(time.Now())
	if result != nil {
		t.Error("CalculateNextRun should return nil for once frequency that already ran")
	}
}

func TestCalculateNextRun_OnceNoStartAt(t *testing.T) {
	now := time.Now()
	s := &Schedule{
		Frequency: FrequencyOnce,
		Status:    ScheduleStatusActive,
	}
	result := s.CalculateNextRun(now)
	if result == nil {
		t.Fatal("CalculateNextRun should return 'from' time for once frequency with no start_at")
	}
	if !result.Equal(now) {
		t.Error("Should return the 'from' time")
	}
}
