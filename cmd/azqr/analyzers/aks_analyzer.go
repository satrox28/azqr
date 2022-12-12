package analyzers

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
)

type AKSAnalyzer struct {
	diagnosticsSettings DiagnosticsSettings
	subscriptionId      string
	ctx                 context.Context
	cred                azcore.TokenCredential
	clustersClient      *armcontainerservice.ManagedClustersClient
}

func NewAKSAnalyzer(subscriptionId string, ctx context.Context, cred azcore.TokenCredential) *AKSAnalyzer {
	diagnosticsSettings, _ := NewDiagnosticsSettings(cred, ctx)
	clustersClient, err := armcontainerservice.NewManagedClustersClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	analyzer := AKSAnalyzer{
		diagnosticsSettings: *diagnosticsSettings,
		subscriptionId:      subscriptionId,
		ctx:                 ctx,
		cred:                cred,
		clustersClient:      clustersClient,
	}
	return &analyzer
}

func (a AKSAnalyzer) Review(resourceGroupName string) ([]AzureServiceResult, error) {
	log.Printf("Analyzing AKS Clusters in Resource Group %s", resourceGroupName)

	clusters, err := a.listClusters(resourceGroupName)
	if err != nil {
		return nil, err
	}
	results := []AzureServiceResult{}
	for _, c := range clusters {
		hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*c.ID)
		if err != nil {
			return nil, err
		}

		zones := true
		for _, profile := range c.Properties.AgentPoolProfiles {
			if profile.AvailabilityZones == nil || (profile.AvailabilityZones != nil && len(profile.AvailabilityZones) <= 1) {
				zones = false
			}
		}

		sku := string(*c.SKU.Tier)
		sla := "None"
		if sku == "Paid" {
			sla = "99.9%"
			if zones {
				sla = "99.95%"
			}
		}

		results = append(results, AzureServiceResult{
			SubscriptionId:     a.subscriptionId,
			ResourceGroup:      resourceGroupName,
			ServiceName:        *c.Name,
			Sku:                sku,
			Sla:                sla,
			Type:               *c.Type,
			AvailabilityZones:  zones,
			PrivateEndpoints:   *c.Properties.APIServerAccessProfile.EnablePrivateCluster,
			DiagnosticSettings: hasDiagnostics,
			CAFNaming:          strings.HasPrefix(*c.Name, "aks"),
		})
	}
	return results, nil
}

func (a AKSAnalyzer) listClusters(resourceGroupName string) ([]*armcontainerservice.ManagedCluster, error) {
	pager := a.clustersClient.NewListByResourceGroupPager(resourceGroupName, nil)

	clusters := make([]*armcontainerservice.ManagedCluster, 0)
	for pager.More() {
		resp, err := pager.NextPage(a.ctx)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, resp.Value...)
	}
	return clusters, nil
}