package config

import (
	"fmt"
	"os"

	"github.com/zinrai/sevalet/internal/models"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server         ServerConfig `yaml:"server"`
	CommandsFile   string       `yaml:"commands_file"`
	DefaultTimeout int          `yaml:"default_timeout"`
	commandList    *models.CommandList
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Loads configuration from the specified file path
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default values
	if config.DefaultTimeout <= 0 {
		config.DefaultTimeout = 30
	}

	// Load command list
	if config.CommandsFile != "" {
		if err := config.loadCommands(); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// Loads command list from the configuration file
func (c *Config) loadCommands() error {
	data, err := os.ReadFile(c.CommandsFile)
	if err != nil {
		return fmt.Errorf("failed to read commands file: %w", err)
	}

	var commandList models.CommandList
	if err := yaml.Unmarshal(data, &commandList); err != nil {
		return fmt.Errorf("failed to parse commands file: %w", err)
	}

	c.commandList = &commandList
	return nil
}

// Returns the loaded command list
func (c *Config) GetCommandList() *models.CommandList {
	return c.commandList
}

// Returns the server address in "host:port" format
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// SetCommandList allows setting a command list directly (useful for testing)
func (c *Config) SetCommandList(commandList *models.CommandList) {
	c.commandList = commandList
}
