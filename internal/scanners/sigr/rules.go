// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package sigr

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
	"github.com/cmendible/azqr/internal/scanners"
)

// GetRules - Returns the rules for the SignalRScanner
func (a *SignalRScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "sigr-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "SignalR should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armsignalr.ResourceInfo)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/signalr-howto-diagnostic-logs",
		},
		"AvailabilityZones": {
			Id:          "sigr-002",
			Category:    "High Availability and Resiliency",
			Subcategory: "Availability Zones",
			Description: "SignalR should have availability zones enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				sku := string(*i.SKU.Name)
				zones := false
				if strings.Contains(sku, "Premium") {
					zones = true
				}
				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/availability-zones",
		},
		"SLA": {
			Id:          "sigr-003",
			Category:    "High Availability and Resiliency",
			Subcategory: "SLA",
			Description: "SignalR should have a SLA",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				return false, "99.9%"
			},
			Url: "https://www.azure.cn/en-us/support/sla/signalr-service/",
		},
		"Private": {
			Id:          "sigr-004",
			Category:    "Security",
			Subcategory: "Networking",
			Description: "SignalR should have private endpoints enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-signalr/howto-private-endpoints",
		},
		"SKU": {
			Id:          "sigr-005",
			Category:    "High Availability and Resiliency",
			Subcategory: "SKU",
			Description: "SignalR SKU",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armsignalr.ResourceInfo)
				return false, string(*i.SKU.Name)
			},
			Url: "https://azure.microsoft.com/en-us/pricing/details/signalr-service/",
		},
		"CAF": {
			Id:          "sigr-006",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "SignalR Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armsignalr.ResourceInfo)
				caf := strings.HasPrefix(*c.Name, "sigr")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"sigr-007": {
			Id:          "sigr-007",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "SignalR should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armsignalr.ResourceInfo)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}
