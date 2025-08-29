package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zinrai/sevalet/internal/config"
	"github.com/zinrai/sevalet/internal/daemon"
)

var (
	daemonConfigFile string
	daemonSocketPath string
	daemonLogLevel   string
	daemonLogFile    string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start mediator daemon",
	Long: `Start the sevalet mediator daemon that executes commands via Unix domain socket.
The daemon validates and executes commands based on the allowed list in the configuration file.`,
	Example: `  sevalet daemon --config /etc/sevalet/daemon.yaml
  sevalet daemon --socket /var/run/sevalet.sock --log-level debug`,
	PreRunE: validateDaemonFlags,
	RunE:    runDaemon,
}

func init() {
	daemonCmd.Flags().StringVarP(&daemonConfigFile, "config", "c", "/etc/sevalet/daemon.yaml", "Configuration file path")
	daemonCmd.Flags().StringVarP(&daemonSocketPath, "socket", "s", "", "Unix domain socket path (overrides config)")
	daemonCmd.Flags().StringVar(&daemonLogLevel, "log-level", "info", "Log level (debug|info|warn|error)")
	daemonCmd.Flags().StringVar(&daemonLogFile, "log-file", "", "Log file path (default: stderr)")
}

func validateDaemonFlags(cmd *cobra.Command, args []string) error {
	// Check if config file exists
	if _, err := os.Stat(daemonConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", daemonConfigFile)
	}

	// Validate log level
	switch daemonLogLevel {
	case "debug", "info", "warn", "error":
		// valid
	default:
		return fmt.Errorf("invalid log level: %s", daemonLogLevel)
	}

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) error {
	// Setup logging
	if daemonLogFile != "" {
		f, err := os.OpenFile(daemonLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// Load configuration
	cfg, err := config.LoadDaemonConfig(daemonConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override socket path if specified
	if daemonSocketPath != "" {
		cfg.SocketPath = daemonSocketPath
	}

	// Set log level
	cfg.LogLevel = daemonLogLevel

	log.Printf("Starting sevalet daemon (version: %s)", version)
	log.Printf("Socket path: %s", cfg.SocketPath)
	log.Printf("Log level: %s", cfg.LogLevel)

	// Start daemon
	d := daemon.New(cfg)
	return d.Start()
}
