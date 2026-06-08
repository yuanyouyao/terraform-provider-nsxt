// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

var cliRouteControllerInterfaceStatisticsClient = infra.NewRouteControllerInterfaceStatisticsClient

func dataSourceNsxtPolicyRouteControllerInterfaceStatistics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerInterfaceStatisticsRead,

		Schema: map[string]*schema.Schema{
			"route_controller_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"route_controller_path", "route_controller_id"},
				Description:  "The policy path of the route controller.",
			},
			"route_controller_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"route_controller_path", "route_controller_id"},
				Description:  "The ID of the route controller.",
			},
			"interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the route controller interface.",
			},
			"interface_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the interface.",
			},
			"per_node_statistics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interface statistics per node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transport_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_router_port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_appliance_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rx": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getLogicalRouterPortCountersSchema(),
							},
						},
						"tx": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getLogicalRouterPortCountersSchema(),
							},
						},
					},
				},
			},
		},
	}
}

func getLogicalRouterPortCountersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"total_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_bytes": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"blocked_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"firewall_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ipv6_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"no_arp_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"no_route_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rpf_check_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ttl_exceeded_dropped_packets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func flattenLogicalRouterPortCounters(counters *model.LogicalRouterPortCounters) []interface{} {
	if counters == nil {
		return nil
	}

	m := make(map[string]interface{})
	if counters.TotalPackets != nil {
		m["total_packets"] = int(*counters.TotalPackets)
	}
	if counters.TotalBytes != nil {
		m["total_bytes"] = int(*counters.TotalBytes)
	}
	if counters.DroppedPackets != nil {
		m["dropped_packets"] = int(*counters.DroppedPackets)
	}
	if counters.BlockedPackets != nil {
		m["blocked_packets"] = int(*counters.BlockedPackets)
	}
	if counters.FirewallDroppedPackets != nil {
		m["firewall_dropped_packets"] = int(*counters.FirewallDroppedPackets)
	}
	if counters.Ipv6DroppedPackets != nil {
		m["ipv6_dropped_packets"] = int(*counters.Ipv6DroppedPackets)
	}
	if counters.NoArpDroppedPackets != nil {
		m["no_arp_dropped_packets"] = int(*counters.NoArpDroppedPackets)
	}
	if counters.NoRouteDroppedPackets != nil {
		m["no_route_dropped_packets"] = int(*counters.NoRouteDroppedPackets)
	}
	if counters.RpfCheckDroppedPackets != nil {
		m["rpf_check_dropped_packets"] = int(*counters.RpfCheckDroppedPackets)
	}
	if counters.TtlExceededDroppedPackets != nil {
		m["ttl_exceeded_dropped_packets"] = int(*counters.TtlExceededDroppedPackets)
	}

	return []interface{}{m}
}

func dataSourceNsxtPolicyRouteControllerInterfaceStatisticsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerID := d.Get("route_controller_id").(string)
	if routeControllerID == "" {
		path := d.Get("route_controller_path").(string)
		routeControllerID = getResourceIDFromResourcePath(path, "route-controllers")
		if routeControllerID == "" {
			segs := strings.Split(path, "/")
			routeControllerID = segs[len(segs)-1]
		}
	}

	interfaceID := d.Get("interface_id").(string)

	client := cliRouteControllerInterfaceStatisticsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, interfaceID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error getting route controller interface statistics for Route Controller %s and Interface %s: %v", routeControllerID, interfaceID, err)
	}

	d.SetId(routeControllerID + "/" + interfaceID + "/statistics")

	if obj.InterfacePath != nil {
		d.Set("interface_path", *obj.InterfacePath)
	}

	if len(obj.PerNodeStatistics) > 0 {
		var nodes []interface{}
		for _, node := range obj.PerNodeStatistics {
			nodeMap := make(map[string]interface{})
			if node.TransportNodeId != nil {
				nodeMap["transport_node_id"] = *node.TransportNodeId
			}
			if node.SubClusterId != nil {
				nodeMap["sub_cluster_id"] = *node.SubClusterId
			}
			if node.LogicalRouterPortId != nil {
				nodeMap["logical_router_port_id"] = *node.LogicalRouterPortId
			}
			if node.VirtualNetworkAppliancePath != nil {
				nodeMap["virtual_network_appliance_path"] = *node.VirtualNetworkAppliancePath
			}
			if node.LastUpdateTimestamp != nil {
				nodeMap["last_update_timestamp"] = int(*node.LastUpdateTimestamp)
			}
			if node.Rx != nil {
				nodeMap["rx"] = flattenLogicalRouterPortCounters(node.Rx)
			}
			if node.Tx != nil {
				nodeMap["tx"] = flattenLogicalRouterPortCounters(node.Tx)
			}
			nodes = append(nodes, nodeMap)
		}
		d.Set("per_node_statistics", nodes)
	} else {
		d.Set("per_node_statistics", nil)
	}

	return nil
}
