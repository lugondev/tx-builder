package cmd

import (
	"github.com/lugondev/tx-builder/cmd/api"
	"github.com/spf13/cobra"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "tx builder",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	// Add Run command
	rootCmd.AddCommand(api.NewRootCommand())

	return rootCmd
}
