package models

import (
	"encoding/json"
	"fmt"
	"io"
)

type ExecuteRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Timeout int      `json:"timeout"` // Timeout in seconds
}

func NewExecuteRequestFromJSON(body io.Reader) (*ExecuteRequest, error) {
	var req ExecuteRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON request: %w", err)
	}

	// Set default timeout value
	if req.Timeout <= 0 {
		req.Timeout = 30 // Default is 30 seconds
	}

	return &req, nil
}

func (r *ExecuteRequest) Validate() error {
	if r.Command == "" {
		return fmt.Errorf("command is not specified")
	}

	// Set timeout upper limit (optional)
	if r.Timeout > 300 {
		return fmt.Errorf("timeout must be 300 seconds or less")
	}

	// Initialize args if nil
	if r.Args == nil {
		r.Args = []string{}
	}

	return nil
}
