package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ExecutionResult struct {
	ExitCode      int
	Stdout        string
	Stderr        string
	ExecutionTime string
	Error         error
}

// Executes the specified command
func ExecuteCommand(ctx context.Context, command string, args []string, timeout int) *ExecutionResult {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, command, args...)

	// Create buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start measuring execution time
	startTime := time.Now()

	// Execute command
	err := cmd.Run()

	// Calculate execution time
	executionTime := time.Since(startTime).String()

	// Create result object
	result := &ExecutionResult{
		Stdout:        strings.TrimSpace(stdout.String()),
		Stderr:        strings.TrimSpace(stderr.String()),
		ExecutionTime: executionTime,
	}

	// Error handling and status code retrieval
	if err != nil {
		// Special handling for timeout errors
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("command execution timed out: %s", err)
			result.ExitCode = -1 // Special exit code for timeout
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			// Command executed successfully but exited with non-zero code
			result.ExitCode = exitErr.ExitCode()
			result.Error = nil // This is a normal execution result, so error is null
		} else {
			// Other errors
			result.Error = fmt.Errorf("error occurred during command execution: %w", err)
			result.ExitCode = -2 // Special exit code for general errors
		}
	} else {
		// Success handling
		result.ExitCode = 0
	}

	return result
}
