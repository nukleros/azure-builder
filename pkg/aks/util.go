package aks

import (
	"fmt"
)

func getResourceGroupNameForCluster(clusterName string) string {
	return fmt.Sprintf("%s-resource-group", clusterName)
}
