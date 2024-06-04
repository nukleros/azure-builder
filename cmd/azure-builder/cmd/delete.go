/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command.
var deleteCmd = &cobra.Command{
	Use:   "delete <resource stack> <inventory file>",
	Short: "Remove an Azure resource stack",
	Long: fmt.Sprintf(`Remove an Azure resource stack.
%s`, supportedResourceStacks),
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
