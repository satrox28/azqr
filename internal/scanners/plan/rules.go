// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package plan

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/cmendible/azqr/internal/scanners"
)

// GetRules - Returns the rules for the AppServiceScanner
func (a *AppServiceScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "plan-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "Plan should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armappservice.Plan)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
		},
		"AvailabilityZones": {
			Id:          "plan-002",
			Category:    "High Availability and Resiliency",
			Subcategory: "Availability Zones",
			Description: "Plan should have availability zones enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armappservice.Plan)
				zones := *i.Properties.ZoneRedundant
				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/reliability/migrate-app-service",
		},
		"SLA": {
			Id:          "plan-003",
			Category:    "High Availability and Resiliency",
			Subcategory: "SLA",
			Description: "Plan should have a SLA",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armappservice.Plan)
				sku := string(*i.SKU.Tier)
				sla := "None"
				if sku != "Free" && sku != "Shared" {
					sla = "99.95%"
				}
				return sla == "None", sla
			},
			Url: "https://www.azure.cn/en-us/support/sla/app-service/",
		},
		"SKU": {
			Id:          "plan-005",
			Category:    "High Availability and Resiliency",
			Subcategory: "SKU",
			Description: "Plan SKU",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armappservice.Plan)
				return false, string(*i.SKU.Name)
			},
			Url: "https://learn.microsoft.com/en-us/azure/app-service/overview-hosting-plans",
		},
		"CAF": {
			Id:          "plan-006",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "Plan Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Plan)
				caf := strings.HasPrefix(*c.Name, "asp")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"plan-007": {
			Id:          "plan-007",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "Plan should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Plan)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}

// GetAppRules - Returns the rules for the AppServiceScanner
func (a *AppServiceScanner) GetAppRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "app-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "App Service should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armappservice.Site)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/app-service/troubleshoot-diagnostic-logs#send-logs-to-azure-monitor",
		},
		"Private": {
			Id:          "app-004",
			Category:    "Security",
			Subcategory: "Networking",
			Description: "App Service should have private endpoints enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armappservice.Site)
				_, pe := scanContext.PrivateEndpoints[*i.ID]
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/app-service/networking/private-endpoint",
		},
		"CAF": {
			Id:          "app-006",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "App Service Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				caf := strings.HasPrefix(*c.Name, "app")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"app-007": {
			Id:          "app-007",
			Category:    "Security",
			Subcategory: "Network Security",
			Description: "App Service should use HTTPS only",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				h := *c.Properties.HTTPSOnly
				return !h, ""
			},
			Url: "https://learn.microsoft.com/azure/app-service/configure-ssl-bindings#enforce-https",
		},
		"app-008": {
			Id:          "app-008",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "App Service should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}

// GetFunctionRules - Returns the rules for the AppServiceScanner
func (a *AppServiceScanner) GetFunctionRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "func-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "Function should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armappservice.Site)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-functions/functions-monitor-log-analytics?tabs=csharp",
		},
		"Private": {
			Id:          "func-004",
			Category:    "Security",
			Subcategory: "Networking",
			Description: "Function should have private endpoints enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armappservice.Site)
				_, pe := scanContext.PrivateEndpoints[*i.ID]
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-functions/functions-create-vnet",
		},
		"CAF": {
			Id:          "func-006",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "Function Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				caf := strings.HasPrefix(*c.Name, "func")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"func-007": {
			Id:          "func-007",
			Category:    "Security",
			Subcategory: "Network Security",
			Description: "Function should use HTTPS only",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				h := c.Properties.HTTPSOnly != nil && *c.Properties.HTTPSOnly
				return !h, ""
			},
			Url: "https://learn.microsoft.com/azure/app-service/configure-ssl-bindings#enforce-https",
		},
		"func-008": {
			Id:          "func-008",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "Function should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappservice.Site)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}
