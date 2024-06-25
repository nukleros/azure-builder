// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package aks

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/nukleros/azure-builder/pkg/config"
	resourcegroup "github.com/nukleros/azure-builder/pkg/resource-group"
)

func CreateAksCluster(
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armcontainerservice.ManagedCluster, error) {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	ctx := context.Background()

	resourceGroup, err := resourcegroup.CreateResourceGroup(aksConfig, credentialsConfig, ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create the resource group: %w", err)
	}

	log.Println("created resource group id:", *resourceGroup.ID)

	managedCluster, err := createManagedCluster(ctx, aksConfig, credentialsConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create managed aks cluster: %w", err)
	}

	log.Println("created aks cluster id:", *managedCluster.ID)
	return managedCluster, nil
}

func GetAksCluster(
	ctx context.Context,
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armcontainerservice.ManagedCluster, error) {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	managedClustersClient, err := credentialsConfig.CreateAzureManagedClustersClient(*aksConfig.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("could not create managed clusters client from credentials config: %w", err)
	}

	clusterResponse, err := managedClustersClient.Get(ctx, *aksConfig.ResourceGroup, *aksConfig.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get managed cluster %s: %w", *aksConfig.Name, err)
	}

	return &clusterResponse.ManagedCluster, nil
}

func createManagedCluster(
	ctx context.Context,
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armcontainerservice.ManagedCluster, error) {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	managedClustersClient, err := credentialsConfig.CreateAzureManagedClustersClient(*aksConfig.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("could not create managed clusters client from credentials config: %w", err)
	}

	pollerResp, err := managedClustersClient.BeginCreateOrUpdate(
		ctx,
		*aksConfig.ResourceGroup,
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
	ctx context.Context,
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) ([]byte, error) {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	managedClustersClient, err := credentialsConfig.CreateAzureManagedClustersClient(*aksConfig.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("could not create managed clusters client from credentials config: %w", err)
	}

	// get kubeconfig for the cluster
	adminClusterCredentials, err := managedClustersClient.ListClusterAdminCredentials(ctx, *aksConfig.ResourceGroup, *aksConfig.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("could not list cluster credentials: %w", err)
	}

	if len(adminClusterCredentials.Kubeconfigs) == 0 {
		return nil, fmt.Errorf("could not retrieve any kube config for created cluster")
	}

	return adminClusterCredentials.Kubeconfigs[0].Value, nil
}

func DeleteAksCluster(
	aksConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) error {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return fmt.Errorf("could not validate credentials config: %w", err)
	}

	ctx := context.TODO()

	// delete the entire resource group that was provisioned for the cluster, this ensures that azure handles all the
	// individual resources the correspond the to the aks cluster deployment
	if err := resourcegroup.CleanupResourceGroup(aksConfig, credentialsConfig, ctx); err != nil {
		return fmt.Errorf("could not clean up resource group for the aks cluster: %w", err)
	}

	return nil
}
