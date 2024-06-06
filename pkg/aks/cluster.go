// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package aks

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nukleros/azure-builder/pkg/config"
)

func CreateAksCluster(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armcontainerservice.ManagedCluster, error) {
	ctx := context.Background()

	resourceGroup, err := createResourceGroup(aksConfig, credentialsConfig, ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create the resource group: %w", err)
	}

	log.Println("created resource group id:", *resourceGroup.ID)

	managedCluster, err := createManagedCluster(aksConfig, credentialsConfig, ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create managed aks cluster: %w", err)
	}

	log.Println("created aks cluster id:", *managedCluster.ID)
	return managedCluster, nil
}

func createManagedCluster(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) (*armcontainerservice.ManagedCluster, error) {
	resourceGroupName := getResourceGroupNameForCluster(*aksConfig.Name)
	managedClustersClient, err := credentialsConfig.CreateAzureManagedClustersClient(resourceGroupName)
	if err != nil {
		fmt.Errorf("could not create managed clusters client from credentials config: %w", err)
	}

	pollerResp, err := managedClustersClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		*aksConfig.Name,
		armcontainerservice.ManagedCluster{
			Location: aksConfig.Region,
			Properties: &armcontainerservice.ManagedClusterProperties{
				DNSPrefix: to.Ptr("aksgosdk"),
				AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
					{
						Name:              to.Ptr("askagent"),
						Count:             to.Ptr[int32](1),
						VMSize:            to.Ptr("Standard_DS2_v2"),
						MaxPods:           to.Ptr[int32](110),
						MinCount:          to.Ptr[int32](1),
						MaxCount:          to.Ptr[int32](100),
						OSType:            to.Ptr(armcontainerservice.OSTypeLinux),
						Type:              to.Ptr(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets),
						EnableAutoScaling: to.Ptr(true),
						Mode:              to.Ptr(armcontainerservice.AgentPoolModeSystem),
					},
				},
				ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
					ClientID: credentialsConfig.ClientID,
					Secret:   credentialsConfig.ClientSecret,
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run BeginCreateOrUpdate for aks cluster: %w", err)
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for completion response for create aks cluster: %w", err)
	}

	return &resp.ManagedCluster, nil
}

func GetKubeConfigForCluster(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) ([]byte, error) {
	resourceGroupName := getResourceGroupNameForCluster(*aksConfig.Name)
	managedClustersClient, err := credentialsConfig.CreateAzureManagedClustersClient(resourceGroupName)
	if err != nil {
		fmt.Errorf("could not create managed clusters client from credentials config: %w", err)
	}

	// get kubeconfig for the cluster
	adminClusterCredentials, err := managedClustersClient.ListClusterAdminCredentials(ctx, resourceGroupName, *aksConfig.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("could not list cluster credentials: %w", err)
	}

	if len(adminClusterCredentials.Kubeconfigs) == 0 {
		return nil, fmt.Errorf("could not retrieve any kube config for created cluster")
	}

	return adminClusterCredentials.Kubeconfigs[0].Value, nil
}

func createResourceGroup(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) (*armresources.ResourceGroup, error) {
	resourceGroupClient, err := credentialsConfig.CreateAzureResourceGroupsClient()
	if err != nil {
		return nil, fmt.Errorf("could not create resource groups client from credentials config: %w", err)
	}

	resourceGroupName := getResourceGroupNameForCluster(*aksConfig.Name)
	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: aksConfig.Region,
		},
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to run CreateOrUpdate: %w", err)
	}

	return &resourceGroupResp.ResourceGroup, nil
}

func DeleteAksCluster(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) error {
	ctx := context.Background()

	// delete the entire resource group that was provisioned for the cluster, this ensures that azure handles all the
	// individual resources the correspond the to the aks cluster deployment
	if err := cleanupResourceGroup(aksConfig, credentialsConfig, ctx); err != nil {
		return fmt.Errorf("could not clean up resource group for the aks cluster: %w", err)
	}

	return nil
}

func cleanupResourceGroup(
	aksConfig *config.AksConfig,
	credentialsConfig *config.AzureCredentialsConfig,
	ctx context.Context,
) error {
	resourceGroupClient, err := credentialsConfig.CreateAzureResourceGroupsClient()
	if err != nil {
		return fmt.Errorf("could not create resource groups client from credentials config: %w", err)
	}

	resourceGroupName := getResourceGroupNameForCluster(*aksConfig.Name)
	log.Println("deleting associated resource groups...")
	pollerResp, err := resourceGroupClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return fmt.Errorf("failed to run being deletion of resourceGroup %s: %w", resourceGroupName, err)
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll for completion response on delete resource group %s: %w", resourceGroupName, err)
	}

	log.Println(fmt.Sprintf("deleted resource group %s", resourceGroupName))

	return nil
}
