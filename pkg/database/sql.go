package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/nukleros/azure-builder/pkg/config"
	resourcegroup "github.com/nukleros/azure-builder/pkg/resource-group"
)

func CreateSqlDb(
	sqlConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armsql.Server, *armsql.Database, error) {
	if err := credentialsConfig.ValidateNotNull(); err != nil {
		return nil, nil, fmt.Errorf("could not validate credentials config: %w", err)
	}

	ctx := context.Background()

	resourceGroup, err := resourcegroup.CreateResourceGroup(sqlConfig, credentialsConfig, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create the resource group: %w", err)
	}
	log.Println("resources group:", *resourceGroup.ID)

	server, err := createSqlServer(ctx, sqlConfig, credentialsConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create sql server: %w", err)
	}
	log.Println("server:", *server.ID)

	database, err := createSqlDatabase(ctx, sqlConfig, credentialsConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create sql database: %w", err)
	}
	log.Println("database:", *database.ID)

	return server, database, nil
}

func createSqlServer(
	ctx context.Context,
	serverConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armsql.Server, error) {

	serversClient, err := credentialsConfig.CreateAzureSqlServersClient()
	if err != nil {
		return nil, fmt.Errorf("could not create servers client: %w", err)
	}

	pollerResp, err := serversClient.BeginCreateOrUpdate(
		ctx,
		*serverConfig.ResourceGroup,
		*serverConfig.Name,
		armsql.Server{
			Location: serverConfig.Region,
			Properties: &armsql.ServerProperties{
				AdministratorLogin:         to.Ptr("dummylogin"),
				AdministratorLoginPassword: to.Ptr("QWE123!@#"),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func createSqlDatabase(
	ctx context.Context,
	dbConfig *config.AzureResourceConfig,
	credentialsConfig *config.AzureCredentialsConfig,
) (*armsql.Database, error) {

	databaseClient, err := credentialsConfig.CreateAzureSqlDatabaseClient()
	if err != nil {
		return nil, fmt.Errorf("could not create database client: %w", err)
	}

	pollerResp, err := databaseClient.BeginCreateOrUpdate(
		ctx,
		*dbConfig.ResourceGroup,
		*dbConfig.Name,
		*dbConfig.Name,
		armsql.Database{
			Location: dbConfig.Region,
			Properties: &armsql.DatabaseProperties{
				ReadScale: to.Ptr(armsql.DatabaseReadScaleDisabled),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Database, nil
}
