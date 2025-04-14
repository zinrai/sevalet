package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zinrai/sevalet/internal/config"
	"github.com/zinrai/sevalet/internal/models"
)

func TestExecuteHandler(t *testing.T) {
	// Skip if running in short mode since this executes real commands
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Set up test command list
	cmdList := &models.CommandList{
		Commands: []models.Command{
			{
				Name:        "echo",
				Description: "Display text",
				AllowedArgs: []string{"hello", "test"},
			},
		},
	}

	// Prepare config and handler
	cfg := &config.Config{}
	cfg.SetCommandList(cmdList)

	handler := NewHandler(cfg)

	// Test valid request
	t.Run("Valid command", func(t *testing.T) {
		validReq := `{"command":"echo","args":["hello"],"timeout":10}`
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(validReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ExecuteHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Unexpected status code: %d", w.Code)
		}

		// Parse response
		var resp models.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if resp.API.Status != "success" {
			t.Errorf("Unexpected API status: %s", resp.API.Status)
		}

		if !resp.Command.Executed || resp.Command.ExitCode != 0 {
			t.Errorf("Command execution result incorrect. Executed: %v, Code: %d",
				resp.Command.Executed, resp.Command.ExitCode)
		}

		if !strings.Contains(resp.Command.Stdout, "hello") {
			t.Errorf("Command output incorrect: %s", resp.Command.Stdout)
		}
	})

	// Test disallowed command
	t.Run("Disallowed command", func(t *testing.T) {
		invalidReq := `{"command":"rm","args":["-rf"],"timeout":10}`
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(invalidReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ExecuteHandler(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Unexpected status code for disallowed command: %d", w.Code)
		}

		// Parse response
		var resp models.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if resp.API.Status != "error" {
			t.Errorf("Expected error status for disallowed command, got: %s", resp.API.Status)
		}

		if resp.API.Code != models.ErrCodeCommandNotAllowed {
			t.Errorf("Expected error code %s, got: %s", models.ErrCodeCommandNotAllowed, resp.API.Code)
		}
	})

	// Test invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"command":"echo",}`
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ExecuteHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Unexpected status code for invalid JSON: %d", w.Code)
		}
	})

	// Test allowed command with disallowed args
	t.Run("Disallowed args", func(t *testing.T) {
		invalidArgsReq := `{"command":"echo","args":["danger"],"timeout":10}`
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(invalidArgsReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ExecuteHandler(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Unexpected status code for disallowed args: %d", w.Code)
		}
	})
}

func TestCommandsHandler(t *testing.T) {
	// Set up test command list
	cmdList := &models.CommandList{
		Commands: []models.Command{
			{
				Name:        "echo",
				Description: "Display text",
				AllowedArgs: []string{"hello", "test"},
			},
			{
				Name:        "ls",
				Description: "List directory contents",
				AllowedArgs: []string{"-l", "/tmp"},
			},
		},
	}

	// Test with valid command list
	t.Run("List available commands", func(t *testing.T) {
		// Create config with command list
		cfg := &config.Config{}
		cfg.SetCommandList(cmdList)

		handler := NewHandler(cfg)

		// Create request
		req := httptest.NewRequest("GET", "/commands", nil)
		w := httptest.NewRecorder()

		handler.CommandsHandler(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Unexpected status code: %d", w.Code)
		}

		// Parse response
		var resp models.CommandsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Check API status
		if resp.API.Status != "success" {
			t.Errorf("Unexpected API status: %s", resp.API.Status)
		}

		// Check commands
		if len(resp.Commands) != 2 {
			t.Errorf("Expected 2 commands, got: %d", len(resp.Commands))
		}

		// Check specific command
		foundEcho := false
		for _, cmd := range resp.Commands {
			if cmd.Name == "echo" {
				foundEcho = true
				if len(cmd.AllowedArgs) != 2 {
					t.Errorf("Expected 2 allowed args for echo, got: %d", len(cmd.AllowedArgs))
				}
			}
		}

		if !foundEcho {
			t.Error("Echo command not found in response")
		}
	})

	// Test with nil command list
	t.Run("Nil command list", func(t *testing.T) {
		// Create config without command list
		cfg := &config.Config{}

		handler := NewHandler(cfg)

		// Create request
		req := httptest.NewRequest("GET", "/commands", nil)
		w := httptest.NewRecorder()

		handler.CommandsHandler(w, req)

		// Check status code
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500 for nil command list, got: %d", w.Code)
		}
	})
}
