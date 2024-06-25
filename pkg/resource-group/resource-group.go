package resourcegroup

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nukleros/azure-builder/pkg/config"
)

func CreateResourceGroup(
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) (*armresources.ResourceGroup, error) {
	resourceGroupClient, err := credentialsConfig.CreateAzureResourceGroupsClient()
	if err != nil {
		return nil, fmt.Errorf("could not create resource groups client from credentials config: %w", err)
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		*aksConfig.ResourceGroup,
		armresources.ResourceGroup{
			Location: aksConfig.Region,
		},
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to run CreateOrUpdate: %w", err)
	}

	return &resourceGroupResp.ResourceGroup, nil
}

func CleanupResourceGroup(
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) error {
	resourceGroupClient, err := credentialsConfig.CreateAzureResourceGroupsClient()
	if err != nil {
		return fmt.Errorf("could not create resource groups client from credentials config: %w", err)
	}

	log.Println("deleting associated resource groups...")
	pollerResp, err := resourceGroupClient.BeginDelete(ctx, *aksConfig.ResourceGroup, nil)
	if err != nil {
		return fmt.Errorf("failed to run being deletion of resourceGroup %s: %w", *aksConfig.ResourceGroup, err)
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll for completion response on delete resource group %s: %w", *aksConfig.ResourceGroup, err)
	}

	log.Println(fmt.Sprintf("deleted resource group %s", *aksConfig.ResourceGroup))

	return nil
}
