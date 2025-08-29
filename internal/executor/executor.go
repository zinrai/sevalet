package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Result contains the result of command execution
type Result struct {
	ExitCode      int
	Stdout        string
	Stderr        string
	ExecutionTime string
	Error         error
}

// ExecuteCommand executes the specified command with timeout
func ExecuteCommand(ctx context.Context, command string, args []string, timeout int) *Result {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, command, args...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start time measurement
	startTime := time.Now()

	// Execute command
	err := cmd.Run()

	// Calculate execution time
	executionTime := time.Since(startTime).String()

	// Create result
	result := &Result{
		Stdout:        strings.TrimSpace(stdout.String()),
		Stderr:        strings.TrimSpace(stderr.String()),
		ExecutionTime: executionTime,
	}

	// Handle errors and exit codes
	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("command execution timed out")
			result.ExitCode = -1 // Special code for timeout
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			// Command executed but returned non-zero exit code
			result.ExitCode = exitErr.ExitCode()
			// Don't set error for non-zero exit codes
			// as this is a normal execution result
		} else {
			// Other execution errors
			result.Error = fmt.Errorf("command execution failed: %w", err)
			result.ExitCode = -2 // Special code for general errors
		}
	} else {
		// Success
		result.ExitCode = 0
	}

	return result
}
