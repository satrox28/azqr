// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package aks

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cmendible/azqr/internal/scanners"
)

func TestAKSScanner_Rules(t *testing.T) {
	type fields struct {
		rule                string
		target              interface{}
		scanContext         *scanners.ScanContext
		diagnosticsSettings scanners.DiagnosticsSettings
	}
	type want struct {
		broken bool
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "AKSScanner DiagnosticSettings",
			fields: fields{
				rule: "DiagnosticSettings",
				target: &armcontainerservice.ManagedCluster{
					ID: to.StringPtr("test"),
				},
				scanContext: &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{
					HasDiagnosticsFunc: func(resourceId string) (bool, error) {
						return true, nil
					},
				},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner AvailabilityZones",
			fields: fields{
				rule: "AvailabilityZones",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierPaid(),
					},
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								AvailabilityZones: []*string{to.StringPtr("1"), to.StringPtr("2"), to.StringPtr("3")},
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner Private Cluster",
			fields: fields{
				rule: "Private",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierPaid(),
					},
					Properties: &armcontainerservice.ManagedClusterProperties{
						APIServerAccessProfile: &armcontainerservice.ManagedClusterAPIServerAccessProfile{
							EnablePrivateCluster: to.BoolPtr(true),
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner SLA Free",
			fields: fields{
				rule: "SLA",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierFree(),
					},
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								AvailabilityZones: []*string{},
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "None",
			},
		},
		{
			name: "AKSScanner SLA Paid",
			fields: fields{
				rule: "SLA",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierPaid(),
					},
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								AvailabilityZones: []*string{},
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "99.9%",
			},
		},
		{
			name: "AKSScanner SLA Paid with AZ",
			fields: fields{
				rule: "SLA",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierPaid(),
					},
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								AvailabilityZones: []*string{to.StringPtr("1"), to.StringPtr("2"), to.StringPtr("3")},
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "99.95%",
			},
		},
		{
			name: "AKSScanner SKU",
			fields: fields{
				rule: "SKU",
				target: &armcontainerservice.ManagedCluster{
					SKU: &armcontainerservice.ManagedClusterSKU{
						Tier: getSKUTierFree(),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "Free",
			},
		},
		{
			name: "AKSScanner CAF",
			fields: fields{
				rule: "CAF",
				target: &armcontainerservice.ManagedCluster{
					Name: to.StringPtr("aks-test"),
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner AADProfile present",
			fields: fields{
				rule: "aks-007",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AADProfile: &armcontainerservice.ManagedClusterAADProfile{},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner AADProfile not present",
			fields: fields{
				rule: "aks-007",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AADProfile: nil,
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner Enable RBAC",
			fields: fields{
				rule: "aks-008",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						EnableRBAC: to.BoolPtr(true),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner Disable RBAC",
			fields: fields{
				rule: "aks-008",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						EnableRBAC: to.BoolPtr(false),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner DisableLocalAccounts",
			fields: fields{
				rule: "aks-009",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						DisableLocalAccounts: to.BoolPtr(true),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner DisableLocalAccounts not present",
			fields: fields{
				rule: "aks-009",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						DisableLocalAccounts: nil,
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner httpApplicationRouting enabled",
			fields: fields{
				rule: "aks-010",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{
							"httpApplicationRouting": {
								Enabled: to.BoolPtr(true),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner httpApplicationRouting disabled",
			fields: fields{
				rule: "aks-010",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{
							"httpApplicationRouting": {
								Enabled: to.BoolPtr(false),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner omsAgent enabled",
			fields: fields{
				rule: "aks-011",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{
							"omsagent": {
								Enabled: to.BoolPtr(true),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner omsAgent disabled",
			fields: fields{
				rule: "aks-011",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{
							"omsagent": {
								Enabled: to.BoolPtr(false),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner omsAgent not present",
			fields: fields{
				rule: "aks-011",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner httpApplicationRouting not present",
			fields: fields{
				rule: "aks-010",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AddonProfiles: map[string]*armcontainerservice.ManagedClusterAddonProfile{},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner OutboundType UserDefinedRouting",
			fields: fields{
				rule: "aks-012",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						NetworkProfile: &armcontainerservice.NetworkProfile{
							OutboundType: getOutboundTypeUserDefinedRouting(),
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "AKSScanner OutboundType not UserDefinedRouting",
			fields: fields{
				rule: "aks-012",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						NetworkProfile: &armcontainerservice.NetworkProfile{
							OutboundType: getOutboundTypeLoadBalancer(),
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner kubenet",
			fields: fields{
				rule: "aks-013",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						NetworkProfile: &armcontainerservice.NetworkProfile{
							NetworkPlugin: getNetworkPluginKubenet(),
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner autoscaling AgentPoolProfiles not present",
			fields: fields{
				rule: "aks-014",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: nil,
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner autoscaling EnableAutoScaling not present",
			fields: fields{
				rule: "aks-014",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner autoscaling false",
			fields: fields{
				rule: "aks-014",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								EnableAutoScaling: to.BoolPtr(false),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "AKSScanner autoscaling",
			fields: fields{
				rule: "aks-014",
				target: &armcontainerservice.ManagedCluster{
					Properties: &armcontainerservice.ManagedClusterProperties{
						AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
							{
								EnableAutoScaling: to.BoolPtr(true),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AKSScanner{
				diagnosticsSettings: tt.fields.diagnosticsSettings,
			}
			rules := s.GetRules()
			b, w := rules[tt.fields.rule].Eval(tt.fields.target, tt.fields.scanContext)
			got := want{
				broken: b,
				result: w,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AKSScanner Rule.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getNetworkPluginKubenet() *armcontainerservice.NetworkPlugin {
	s := armcontainerservice.NetworkPluginKubenet
	return &s
}

func getSKUTierFree() *armcontainerservice.ManagedClusterSKUTier {
	s := armcontainerservice.ManagedClusterSKUTierFree
	return &s
}

func getSKUTierPaid() *armcontainerservice.ManagedClusterSKUTier {
	s := armcontainerservice.ManagedClusterSKUTierPaid
	return &s
}

func getOutboundTypeUserDefinedRouting() *armcontainerservice.OutboundType {
	s := armcontainerservice.OutboundTypeUserDefinedRouting
	return &s
}

func getOutboundTypeLoadBalancer() *armcontainerservice.OutboundType {
	s := armcontainerservice.OutboundTypeLoadBalancer
	return &s
}
