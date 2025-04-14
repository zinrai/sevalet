package executor

import (
	"testing"

	"github.com/zinrai/sevalet/internal/models"
)

func TestValidateCommand(t *testing.T) {
	// Set up allowed command list
	commandList := &models.CommandList{
		Commands: []models.Command{
			{
				Name:        "ls",
				Description: "List directory contents",
				AllowedArgs: []string{"-l", "-a", "/tmp"},
			},
			{
				Name:        "echo",
				Description: "Display text",
				AllowedArgs: []string{"hello", "world"},
			},
		},
	}

	// Test cases
	tests := []struct {
		name      string
		cmd       string
		args      []string
		expectErr bool
	}{
		{"Allowed command with allowed args", "ls", []string{"-l", "/tmp"}, false},
		{"Disallowed command", "rm", []string{}, true},
		{"Allowed command with disallowed args", "ls", []string{"-f"}, true},
		{"Empty command", "", []string{}, true},
		{"Command with no args", "echo", []string{}, false},
		{"Mixed allowed and disallowed args", "ls", []string{"-l", "-r"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCommand(tc.cmd, tc.args, commandList)
			if (err != nil) != tc.expectErr {
				t.Errorf("ValidateCommand(%q, %v) error = %v, expected %v",
					tc.cmd, tc.args, err, tc.expectErr)
			}
		})
	}

	// Test with nil command list
	t.Run("Nil command list", func(t *testing.T) {
		err := ValidateCommand("ls", []string{"-l"}, nil)
		if err == nil {
			t.Error("Expected error with nil command list, got nil")
		}
	})
}
