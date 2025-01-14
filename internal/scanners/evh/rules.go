// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package evh

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub"
	"github.com/cmendible/azqr/internal/scanners"
)

// GetRules - Returns the rules for the EventHubScanner
func (a *EventHubScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "evh-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "Event Hub Namespace should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armeventhub.EHNamespace)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/event-hubs/monitor-event-hubs#collection-and-routing",
		},
		"AvailabilityZones": {
			Id:          "evh-002",
			Category:    "High Availability and Resiliency",
			Subcategory: "Availability Zones",
			Description: "Event Hub Namespace should have availability zones enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armeventhub.EHNamespace)
				zones := *i.Properties.ZoneRedundant
				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/event-hubs/event-hubs-premium-overview#high-availability-with-availability-zones",
		},
		"SLA": {
			Id:          "evh-003",
			Category:    "High Availability and Resiliency",
			Subcategory: "SLA",
			Description: "Event Hub Namespace should have a SLA",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armeventhub.EHNamespace)
				sku := string(*i.SKU.Name)
				sla := "99.95%"
				if !strings.Contains(sku, "Basic") && !strings.Contains(sku, "Standard") {
					sla = "99.99%"
				}
				return false, sla
			},
			Url: "https://www.azure.cn/en-us/support/sla/event-hubs/",
		},
		"Private": {
			Id:          "evh-004",
			Category:    "Security",
			Subcategory: "Networking",
			Description: "Event Hub Namespace should have private endpoints enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armeventhub.EHNamespace)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/event-hubs/network-security",
		},
		"SKU": {
			Id:          "evh-005",
			Category:    "High Availability and Resiliency",
			Subcategory: "SKU",
			Description: "Event Hub Namespace SKU",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armeventhub.EHNamespace)
				return false, string(*i.SKU.Name)
			},
			Url: "https://learn.microsoft.com/en-us/azure/event-hubs/compare-tiers",
		},
		"CAF": {
			Id:          "evh-006",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "Event Hub Namespace Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armeventhub.EHNamespace)
				caf := strings.HasPrefix(*c.Name, "evh")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"evh-007": {
			Id:          "evh-007",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "Event Hub should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armeventhub.EHNamespace)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
		"evh-008": {
			Id:          "evh-008",
			Category:    "Security",
			Subcategory: "Identity and Access Control",
			Description: "Event Hub should have local authentication disabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armeventhub.EHNamespace)
				return c.Properties.DisableLocalAuth != nil && !*c.Properties.DisableLocalAuth, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/event-hubs/authorize-access-event-hubs#shared-access-signatures",
		},
	}
}
