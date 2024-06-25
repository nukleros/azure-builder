package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/nukleros/azure-builder/pkg/config"
	resourcegroup "github.com/nukleros/azure-builder/pkg/resource-group"
)

func CreateBlobStore(
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armstorage.Account, error) {

	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	ctx := context.Background()

	resourceGroup, err := resourcegroup.CreateResourceGroup(aksConfig, credentialsConfig, ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create the resource group: %w", err)
	}

	log.Println("created resource group id:", *resourceGroup.ID)

	storageAccount, err := createStorageAccount(ctx, aksConfig, credentialsConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create the blob storage account: %w", err)
	}

	log.Println("created blob storage account:", *storageAccount.ID)
	return storageAccount, nil
}

func createStorageAccount(
	ctx context.Context,
	storageConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armstorage.Account, error) {

	accountsClient, err := credentialsConfig.CreateStorageAccountsClient()
	if err != nil {
		return nil, fmt.Errorf("could not validate storage account: %w", err)
	}

	pollerResp, err := accountsClient.BeginCreate(
		ctx,
		*storageConfig.ResourceGroup,
		*storageConfig.Name,
		armstorage.AccountCreateParameters{
			Kind: to.Ptr(armstorage.KindStorageV2),
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS),
			},
			Location: storageConfig.Region,
			Properties: &armstorage.AccountPropertiesCreateParameters{
				AccessTier: to.Ptr(armstorage.AccessTierCool),
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
					},
					KeySource: to.Ptr(armstorage.KeySourceMicrosoftStorage),
				},
			},
		}, nil)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Account, nil
}
