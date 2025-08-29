package models

import (
	"encoding/json"
	"fmt"
	"io"
)

// Command represents an allowed command with its arguments
type Command struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	AllowedArgs []string `yaml:"allowed_args" json:"allowed_args"`
}

// CommandList contains all allowed commands
type CommandList struct {
	Commands []Command `yaml:"commands" json:"commands"`
}

// FindCommand searches for a command by name
func (cl *CommandList) FindCommand(name string) *Command {
	for i := range cl.Commands {
		if cl.Commands[i].Name == name {
			return &cl.Commands[i]
		}
	}
	return nil
}

// IsArgAllowed checks if an argument is in the allowed list
func (c *Command) IsArgAllowed(arg string) bool {
	for _, allowedArg := range c.AllowedArgs {
		if allowedArg == arg {
			return true
		}
	}
	return false
}

// ExecuteRequest represents an HTTP request to execute a command
type ExecuteRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Timeout int      `json:"timeout"`
}

// NewExecuteRequestFromJSON creates a request from JSON body
func NewExecuteRequestFromJSON(body io.Reader) (*ExecuteRequest, error) {
	var req ExecuteRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON request: %w", err)
	}

	// Set default timeout
	if req.Timeout <= 0 {
		req.Timeout = 30
	}

	// Initialize args if nil
	if req.Args == nil {
		req.Args = []string{}
	}

	return &req, nil
}

// Validate performs basic validation on the request
func (r *ExecuteRequest) Validate() error {
	if r.Command == "" {
		return fmt.Errorf("command is not specified")
	}

	// Enforce maximum timeout
	if r.Timeout > 300 {
		return fmt.Errorf("timeout must be 300 seconds or less")
	}

	return nil
}

// HTTPResponse represents the API response
type HTTPResponse struct {
	Success       bool   `json:"success"`
	ExitCode      int    `json:"exit_code,omitempty"`
	Stdout        string `json:"stdout,omitempty"`
	Stderr        string `json:"stderr,omitempty"`
	ExecutionTime string `json:"execution_time,omitempty"`
	Error         string `json:"error,omitempty"`
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp     string   `json:"timestamp"`
	Level         string   `json:"level"`
	Mode          string   `json:"mode"`
	Event         string   `json:"event"`
	Command       string   `json:"command,omitempty"`
	Args          []string `json:"args,omitempty"`
	ExitCode      int      `json:"exit_code,omitempty"`
	ExecutionTime string   `json:"execution_time,omitempty"`
	Method        string   `json:"method,omitempty"`
	Path          string   `json:"path,omitempty"`
	RemoteAddr    string   `json:"remote_addr,omitempty"`
	Status        int      `json:"status,omitempty"`
	Latency       string   `json:"latency,omitempty"`
	Error         string   `json:"error,omitempty"`
}
