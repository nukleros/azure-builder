package config

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type AzureCredentialsConfig struct {
	ClientID       *string `json:"clientId"`
	ClientSecret   *string `json:"clientSecret"`
	SubscriptionID *string `json:"subscriptionId"`
}

func (config *AzureCredentialsConfig) ValidateNotNull() error {
	if config.ClientID == nil {
		return fmt.Errorf("could not find Client ID in credentials config")
	}

	if config.ClientSecret == nil {
		return fmt.Errorf("could not find Client Secret in credentials config")
	}

	if config.SubscriptionID == nil {
		return fmt.Errorf("could not find Subscription ID in credentials config")
	}

	return nil
}

func (config *AzureCredentialsConfig) CreateAzureResourceGroupsClient() (*armresources.ResourceGroupsClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create default azure credentials: %w", err)
	}

	resourcesClientFactory, err := armresources.NewClientFactory(*config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create arm resources client factory: %w", err)
	}

	resourceGroupClient := resourcesClientFactory.NewResourceGroupsClient()
	return resourceGroupClient, nil
}

func (config *AzureCredentialsConfig) CreateAzureManagedClustersClient(resourceGroupName string) (*armcontainerservice.ManagedClustersClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create default azure credentials: %w", err)
	}

	containerserviceClientFactory, err := armcontainerservice.NewClientFactory(*config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create arm container service client: %w", err)
	}
	managedClustersClient := containerserviceClientFactory.NewManagedClustersClient()
	return managedClustersClient, nil
}
