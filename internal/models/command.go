package models

type Command struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	AllowedArgs []string `yaml:"allowed_args" json:"allowed_args"`
}

type CommandList struct {
	Commands []Command `yaml:"commands" json:"commands"`
}

// Searches for a command with the specified name
// Returns nil if not found
func (cl *CommandList) FindCommand(name string) *Command {
	for i := range cl.Commands {
		if cl.Commands[i].Name == name {
			return &cl.Commands[i]
		}
	}
	return nil
}

// Checks if the specified argument is allowed for this command
func (c *Command) IsArgAllowed(arg string) bool {
	for _, allowedArg := range c.AllowedArgs {
		if allowedArg == arg {
			return true
		}
	}
	return false
}

// Checks if all arguments in the list are allowed
func (c *Command) AreArgsAllowed(args []string) bool {
	for _, arg := range args {
		if !c.IsArgAllowed(arg) {
			return false
		}
	}
	return true
}
