package models

import (
	"strings"
	"testing"
)

func TestExecuteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *ExecuteRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    []string{"-l"},
				Timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "empty command",
			request: &ExecuteRequest{
				Command: "",
				Args:    []string{"-l"},
				Timeout: 30,
			},
			wantErr: true,
			errMsg:  "command is not specified",
		},
		{
			name: "whitespace only command",
			request: &ExecuteRequest{
				Command: "   ",
				Args:    []string{},
				Timeout: 30,
			},
			wantErr: false, // Note: Validate doesn't trim, so this passes
		},
		{
			name: "timeout exceeds maximum",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    []string{},
				Timeout: 301,
			},
			wantErr: true,
			errMsg:  "timeout must be 300 seconds or less",
		},
		{
			name: "timeout at maximum",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    []string{},
				Timeout: 300,
			},
			wantErr: false,
		},
		{
			name: "negative timeout",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    []string{},
				Timeout: -1,
			},
			wantErr: false, // Negative timeouts are handled by setting default
		},
		{
			name: "zero timeout",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    []string{},
				Timeout: 0,
			},
			wantErr: false, // Zero timeout is handled by setting default
		},
		{
			name: "nil args",
			request: &ExecuteRequest{
				Command: "ls",
				Args:    nil,
				Timeout: 30,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestNewExecuteRequestFromJSON(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		wantCommand   string
		wantTimeout   int
		wantArgsCount int
		wantErr       bool
	}{
		{
			name:          "valid json with all fields",
			json:          `{"command":"ls","args":["-l","/tmp"],"timeout":60}`,
			wantCommand:   "ls",
			wantTimeout:   60,
			wantArgsCount: 2,
			wantErr:       false,
		},
		{
			name:          "json without timeout (should use default)",
			json:          `{"command":"echo","args":["hello"]}`,
			wantCommand:   "echo",
			wantTimeout:   30, // default
			wantArgsCount: 1,
			wantErr:       false,
		},
		{
			name:          "json with null args (should initialize to empty slice)",
			json:          `{"command":"uptime","args":null,"timeout":10}`,
			wantCommand:   "uptime",
			wantTimeout:   10,
			wantArgsCount: 0,
			wantErr:       false,
		},
		{
			name:          "json without args field",
			json:          `{"command":"date"}`,
			wantCommand:   "date",
			wantTimeout:   30,
			wantArgsCount: 0,
			wantErr:       false,
		},
		{
			name:    "invalid json",
			json:    `{"command":}`,
			wantErr: true,
		},
		{
			name:          "empty json",
			json:          `{}`,
			wantCommand:   "",
			wantTimeout:   30,
			wantArgsCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			req, err := NewExecuteRequestFromJSON(reader)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewExecuteRequestFromJSON() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewExecuteRequestFromJSON() unexpected error = %v", err)
				return
			}

			if req.Command != tt.wantCommand {
				t.Errorf("Command = %v, want %v", req.Command, tt.wantCommand)
			}

			if req.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", req.Timeout, tt.wantTimeout)
			}

			if len(req.Args) != tt.wantArgsCount {
				t.Errorf("Args length = %v, want %v", len(req.Args), tt.wantArgsCount)
			}

			// Ensure Args is not nil
			if req.Args == nil {
				t.Error("Args should not be nil, expected empty slice")
			}
		})
	}
}
