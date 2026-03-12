package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danielbrodie/metabase-cli/internal/api"
	"github.com/danielbrodie/metabase-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	profileFlag string
	jsonFlag    bool
	verboseFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "metabase",
	Short: "CLI for Metabase API operations",
	Long:  "metabase — A zero-dependency CLI for Metabase. Query cards, dashboards, and collections from the command line.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// SetVersion sets the version string shown by --version.
func SetVersion(v string) {
	rootCmd.Version = v
}

func getClient() (*api.Client, *config.Config) {
	cfg, err := config.Load()
	if err != nil {
		die("failed to load config: %v", err)
	}
	return api.New(cfg.GetProfile(profileFlag)), cfg
}

// must exits with a JSON error if err != nil, otherwise returns data.
func must(data json.RawMessage, err error) json.RawMessage {
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"success\":false,\"error\":{\"code\":\"API_ERROR\",\"message\":%q}}\n", err.Error())
		os.Exit(1)
	}
	return data
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "p", "default", "Named profile")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "JSON output")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
}
