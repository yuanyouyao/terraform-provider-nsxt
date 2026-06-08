// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

var cliRouteControllerStateClient = infra.NewRouteControllerStateClient

func dataSourceNsxtPolicyRouteControllerState() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerStateRead,

		Schema: map[string]*schema.Schema{
			"route_controller_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"route_controller_id", "route_controller_path"},
				Description:  "The ID of the route controller.",
			},
			"route_controller_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"route_controller_id", "route_controller_path"},
				Description:  "The policy path of the route controller.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source for which state is retrieved.",
			},
			"last_update_timestamp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the status was last updated.",
			},
			"logical_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the route controller logical gateway.",
			},
			"virtual_network_appliance_cluster_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy path of virtual network appliance cluster.",
			},
			"per_node_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Per node status of the route controller.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"high_availability_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "High availability status on virtual network appliance node.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of node.",
						},
						"service_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the service gateway where the status is retrieved.",
						},
						"virtual_network_appliance_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy path of virtual network appliance where the node status is retrieved.",
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerStateRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliRouteControllerStateClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	routeControllerID := d.Get("route_controller_id").(string)
	if routeControllerID == "" {
		path := d.Get("route_controller_path").(string)
		routeControllerID = getResourceIDFromResourcePath(path, "route-controllers")
		if routeControllerID == "" {
			segs := strings.Split(path, "/")
			routeControllerID = segs[len(segs)-1]
		}
	}

	var source *string
	if src, ok := d.GetOk("source"); ok {
		s := src.(string)
		source = &s
	}

	state, err := client.Get(routeControllerID, source)
	if err != nil {
		return fmt.Errorf("error getting route controller state for ID %s: %v", routeControllerID, err)
	}

	if state.LastUpdateTimestamp != nil {
		d.Set("last_update_timestamp", *state.LastUpdateTimestamp)
	}
	if state.LogicalGatewayId != nil {
		d.Set("logical_gateway_id", *state.LogicalGatewayId)
	}
	if state.VirtualNetworkApplianceClusterPath != nil {
		d.Set("virtual_network_appliance_cluster_path", *state.VirtualNetworkApplianceClusterPath)
	}

	nodeStatuses := make([]map[string]interface{}, 0, len(state.PerNodeStatus))
	for _, ns := range state.PerNodeStatus {
		item := make(map[string]interface{})
		if ns.HighAvailabilityStatus != nil {
			item["high_availability_status"] = *ns.HighAvailabilityStatus
		}
		if ns.NodeType != nil {
			item["node_type"] = *ns.NodeType
		}
		if ns.ServiceGatewayId != nil {
			item["service_gateway_id"] = *ns.ServiceGatewayId
		}
		if ns.VirtualNetworkAppliancePath != nil {
			item["virtual_network_appliance_path"] = *ns.VirtualNetworkAppliancePath
		}
		nodeStatuses = append(nodeStatuses, item)
	}
	d.Set("per_node_status", nodeStatuses)
	d.SetId(routeControllerID)

	return nil
}
