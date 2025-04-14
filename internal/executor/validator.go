package executor

import (
	"fmt"

	"github.com/zinrai/sevalet/internal/models"
)

// Verifies if the command and arguments are allowed
func ValidateCommand(cmd string, args []string, commandList *models.CommandList) error {
	// Check if command list is nil
	if commandList == nil {
		return fmt.Errorf("command list is not available")
	}

	// Check if command exists in the allowed list
	command := commandList.FindCommand(cmd)
	if command == nil {
		return fmt.Errorf("command '%s' is not allowed", cmd)
	}

	// Check if arguments are allowed
	for _, arg := range args {
		if !command.IsArgAllowed(arg) {
			return fmt.Errorf("argument '%s' is not allowed for command '%s'", arg, cmd)
		}
	}

	return nil
}
