package executor

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestExecuteCommand(t *testing.T) {
	ctx := context.Background()

	// Test echo command
	t.Run("Echo command", func(t *testing.T) {
		result := ExecuteCommand(ctx, "echo", []string{"hello"}, 10)

		if result.Error != nil {
			t.Errorf("Failed to execute echo command: %v", result.Error)
		}

		if result.ExitCode != 0 {
			t.Errorf("Unexpected exit code: %d", result.ExitCode)
		}

		if !strings.Contains(result.Stdout, "hello") {
			t.Errorf("Stdout does not contain expected output: %s", result.Stdout)
		}

		// Check execution time
		if result.ExecutionTime == "" {
			t.Error("Execution time not set")
		}
	})

	// Test nonexistent command
	t.Run("Nonexistent command", func(t *testing.T) {
		result := ExecuteCommand(ctx, "nonexistentcommand", []string{}, 10)

		if result.Error == nil {
			t.Error("No error returned for nonexistent command")
		}

		if result.ExitCode != -2 { // General error code
			t.Errorf("Expected exit code -2, got: %d", result.ExitCode)
		}
	})

	// Test command with non-zero exit code
	t.Run("Non-zero exit code", func(t *testing.T) {
		// Using sh -c to create a command that exits with code 1
		result := ExecuteCommand(ctx, "sh", []string{"-c", "exit 1"}, 10)

		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}

		if result.ExitCode != 1 {
			t.Errorf("Expected exit code 1, got: %d", result.ExitCode)
		}
	})

	// Test timeout
	t.Run("Command timeout", func(t *testing.T) {
		// Skip in short test mode as it takes time
		if testing.Short() {
			t.Skip("Skipping timeout test in short mode")
		}

		// Command will sleep for 2 seconds but timeout is set to 1 second
		result := ExecuteCommand(ctx, "sleep", []string{"2"}, 1)

		if result.Error == nil {
			t.Error("No error returned for timeout")
		}

		if result.ExitCode != -1 { // Timeout error code
			t.Errorf("Expected exit code -1 for timeout, got: %d", result.ExitCode)
		}
	})

	// Test standard error output
	t.Run("Standard error output", func(t *testing.T) {
		// Use ls with a nonexistent directory to generate stderr
		result := ExecuteCommand(ctx, "ls", []string{"/nonexistentdirectory12345"}, 10)

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for nonexistent directory")
		}

		if result.Stderr == "" {
			t.Error("Expected error output in stderr")
		}
	})

	// Test execution time measurement
	t.Run("Execution time measurement", func(t *testing.T) {
		startTime := time.Now()
		result := ExecuteCommand(ctx, "sleep", []string{"0.5"}, 10)
		elapsed := time.Since(startTime)

		if result.Error != nil {
			t.Errorf("Failed to execute sleep command: %v", result.Error)
		}

		// Parse the execution time
		executionDuration, err := time.ParseDuration(result.ExecutionTime)
		if err != nil {
			t.Errorf("Failed to parse execution time: %v", err)
		}

		// Verify it's in a reasonable range (at least 0.5s but not too much more)
		if executionDuration < 500*time.Millisecond || executionDuration > elapsed*2 {
			t.Errorf("Execution time measurement seems off: %v", executionDuration)
		}
	})
}
