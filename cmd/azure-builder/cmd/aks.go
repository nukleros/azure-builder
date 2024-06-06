/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/nukleros/azure-builder/pkg/aks"
	"github.com/nukleros/azure-builder/pkg/config"
	"github.com/spf13/cobra"
)

var aksConfigPath string

// createAksCmd represents the create command for an AKS cluster.
var createAksCmd = &cobra.Command{
	Use:   "aks",
	Short: "Provision an AKS cluster",
	Long:  fmt.Sprintf(`Provision an AKS cluster`),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load credentials used to connect to Azure
		credsFilePath, err := os.ReadFile(azureCredentialsPath)
		if err != nil {
			return fmt.Errorf("could not read credentials file: %w", err)
		}

		var credentialsConfig config.AzureCredentialsConfig
		if err = json.Unmarshal(credsFilePath, &credentialsConfig); err != nil {
			return fmt.Errorf("could not JSON unmarshal credentials config: %w", err)
		}

		if err = credentialsConfig.ValidateNotNull(); err != nil {
			return fmt.Errorf("could not validate credentials config: %w", err)
		}

		// Load aks config file used to create the cluster
		aksConfigBytes, err := os.ReadFile(aksConfigPath)
		if err != nil {
			return fmt.Errorf("could not read aks config file: %w", err)
		}

		var aksConfig config.AksConfig
		if err = yaml.Unmarshal(aksConfigBytes, &aksConfig); err != nil {
			return fmt.Errorf("could not YAML unmarshal aks config: %w", err)
		}

		if err = aksConfig.ValidateNotNull(); err != nil {
			return fmt.Errorf("could not validate aks config: %w", err)
		}

		// Create cluster
		if _, err = aks.CreateAksCluster(&aksConfig, &credentialsConfig); err != nil {
			return fmt.Errorf("could not create aks cluster: %w", err)
		}

		return nil
	},
}

func init() {
	createCmd.AddCommand(createAksCmd)
	createAksCmd.Flags().StringVarP(&aksConfigPath, "aks-config", "c", "",
		"Location to aks config used to create the resource")
	createAksCmd.Flags().StringVarP(&azureCredentialsPath, "creds-path", "p", "",
		"Location to JSON file containing Azure credentials. To generate one, use the Azure CLI and refer to command 'az ad sp create-for-rbac'")

	createAksCmd.MarkFlagRequired("creds-path")
	createAksCmd.MarkFlagRequired("aks-config")
}

// createAKSCmd represents the create command for an AKS cluster.
var deleteAksCmd = &cobra.Command{
	Use:   "aks",
	Short: "Delete an AKS cluster",
	Long:  fmt.Sprintf(`Delete an AKS cluster`),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load credentials used to connect to Azure
		credsFilePath, err := os.ReadFile(azureCredentialsPath)
		if err != nil {
			return fmt.Errorf("could not read credentials file: %w", err)
		}

		var credentialsConfig config.AzureCredentialsConfig
		if err = json.Unmarshal(credsFilePath, &credentialsConfig); err != nil {
			return fmt.Errorf("could not JSON unmarshall credentials config: %w", err)
		}

		if err = credentialsConfig.ValidateNotNull(); err != nil {
			return fmt.Errorf("could not validate credentials config: %w", err)
		}

		// Load aks config file used to create the cluster
		aksConfigBytes, err := os.ReadFile(aksConfigPath)
		if err != nil {
			return fmt.Errorf("could not read aks config file: %w", err)
		}

		var aksConfig config.AksConfig
		if err = yaml.Unmarshal(aksConfigBytes, &aksConfig); err != nil {
			return fmt.Errorf("could not YAML unmarshal aks config: %w", err)
		}

		if err = aksConfig.ValidateNotNull(); err != nil {
			return fmt.Errorf("could not validate aks config: %w", err)
		}

		// Delete aks cluster
		if err = aks.DeleteAksCluster(&aksConfig, &credentialsConfig); err != nil {
			return fmt.Errorf("could not delete aks cluster: %w", err)
		}

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteAksCmd)
	deleteAksCmd.Flags().StringVarP(&aksConfigPath, "aks-config", "c", "",
		"Location to aks config used to create the resource")
	deleteAksCmd.Flags().StringVarP(&azureCredentialsPath, "creds-path", "p", "",
		"Location to JSON file containing Azure credentials. To generate one, use the Azure CLI and refer to command 'az ad sp create-for-rbac'")

	deleteAksCmd.MarkFlagRequired("creds-path")
	deleteAksCmd.MarkFlagRequired("aks-config")
}
