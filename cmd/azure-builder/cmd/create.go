/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Provision an Azure resource stack",
	Long: fmt.Sprintf(`Provision an Azure resource stack.
%s`, supportedResourceStacks),
}

func init() {
	rootCmd.AddCommand(createCmd)
}
