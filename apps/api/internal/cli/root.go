package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "billing-cli",
	Short: "Orvex Billing CLI",
	Long:  `A CLI tool for managing the Orvex billing system, including user management and database operations.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(dbCmd)
}