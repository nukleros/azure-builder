package config

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
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
	if err := config.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

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
	if err := config.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

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

func (config *AzureCredentialsConfig) CreateAzureSqlDatabaseClient() (*armsql.DatabasesClient, error) {
	if err := config.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create default azure credentials: %w", err)
	}

	sqlClientFactory, err := armsql.NewClientFactory(*config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new sql client: %w", err)
	}

	databasesClient := sqlClientFactory.NewDatabasesClient()
	return databasesClient, nil
}

func (config *AzureCredentialsConfig) CreateAzureSqlServersClient() (*armsql.ServersClient, error) {
	if err := config.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create default azure credentials: %w", err)
	}

	sqlClientFactory, err := armsql.NewClientFactory(*config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new client factory: %w", err)
	}

	serversClient := sqlClientFactory.NewServersClient()
	return serversClient, nil
}

func (config *AzureCredentialsConfig) CreateStorageAccountsClient() (*armstorage.AccountsClient, error) {
	if err := config.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create default azure credentials: %w", err)
	}

	storageClientFactory, err := armstorage.NewClientFactory(*config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new client factory: %w", err)
	}

	accountsClient := storageClientFactory.NewAccountsClient()

	return accountsClient, nil
}
