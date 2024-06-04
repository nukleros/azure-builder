/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createInventoryFile string
	inputInventoryFile  string
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create <resource stack> <config file>",
	Short: "Provision an Azure resource stack",
	Long: fmt.Sprintf(`Provision an Azure resource stack.
%s`, supportedResourceStacks),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
