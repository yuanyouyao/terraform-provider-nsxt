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

var cliRouteControllerBgpNeighborsStatusClient = infra.NewRouteControllerBgpNeighborsStatusClient

func dataSourceNsxtPolicyRouteControllerBgpNeighborsStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpNeighborsStatusRead,

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
			"bgp_neighbor_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BGP neighbor type: INTER_SR or USER.",
			},
			"enforcement_point_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enforcement point path.",
			},
			"stats_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Stats type.",
			},
			"transport_node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Transport node ID.",
			},
			"virtual_network_appliance_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual network appliance path.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of BGP neighbor statuses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"neighbor_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_drop_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"established_connection_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"graceful_restart_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hold_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_dynamic": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"keep_alive_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"local_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"messages_received": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"messages_sent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"neighbor_edge_node": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"neighbor_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"neighbor_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_as_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remote_site_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_controller_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_since_established": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_in_prefix_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_out_prefix_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_appliance_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"announced_capabilities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"negotiated_capabilities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"address_families": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"in_prefix_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"out_prefix_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpNeighborsStatusRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerBgpNeighborsStatusClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	var bgpNeighborType *string
	if val, ok := d.GetOk("bgp_neighbor_type"); ok {
		s := val.(string)
		bgpNeighborType = &s
	}

	var enforcementPointPath *string
	if val, ok := d.GetOk("enforcement_point_path"); ok {
		s := val.(string)
		enforcementPointPath = &s
	}

	var statsType *string
	if val, ok := d.GetOk("stats_type"); ok {
		s := val.(string)
		statsType = &s
	}

	var transportNodeID *string
	if val, ok := d.GetOk("transport_node_id"); ok {
		s := val.(string)
		transportNodeID = &s
	}

	var virtualNetworkAppliancePath *string
	if val, ok := d.GetOk("virtual_network_appliance_path"); ok {
		s := val.(string)
		virtualNetworkAppliancePath = &s
	}

	listResult, err := client.List(routeControllerID, bgpNeighborType, nil, enforcementPointPath, nil, nil, nil, nil, nil, statsType, transportNodeID, virtualNetworkAppliancePath)
	if err != nil {
		return fmt.Errorf("error listing BGP neighbor status for Route Controller %s: %v", routeControllerID, err)
	}

	d.SetId(routeControllerID)

	results := make([]interface{}, 0, len(listResult.Results))
	for _, status := range listResult.Results {
		item := make(map[string]interface{})

		if status.NeighborAddress != nil {
			item["neighbor_address"] = *status.NeighborAddress
		}
		if status.ConnectionState != nil {
			item["connection_state"] = *status.ConnectionState
		}
		if status.ConnectionDropCount != nil {
			item["connection_drop_count"] = int(*status.ConnectionDropCount)
		}
		if status.EstablishedConnectionCount != nil {
			item["established_connection_count"] = int(*status.EstablishedConnectionCount)
		}
		if status.GracefulRestartMode != nil {
			item["graceful_restart_mode"] = *status.GracefulRestartMode
		}
		if status.HoldTime != nil {
			item["hold_time"] = int(*status.HoldTime)
		}
		if status.IsDynamic != nil {
			item["is_dynamic"] = *status.IsDynamic
		}
		if status.KeepAliveInterval != nil {
			item["keep_alive_interval"] = int(*status.KeepAliveInterval)
		}
		if status.LocalPort != nil {
			item["local_port"] = int(*status.LocalPort)
		}
		if status.MessagesReceived != nil {
			item["messages_received"] = int(*status.MessagesReceived)
		}
		if status.MessagesSent != nil {
			item["messages_sent"] = int(*status.MessagesSent)
		}
		if status.NeighborEdgeNode != nil {
			item["neighbor_edge_node"] = *status.NeighborEdgeNode
		}
		if status.NeighborPath != nil {
			item["neighbor_path"] = *status.NeighborPath
		}
		if status.NeighborRouterId != nil {
			item["neighbor_router_id"] = *status.NeighborRouterId
		}
		if status.RemoteAsNumber != nil {
			item["remote_as_number"] = *status.RemoteAsNumber
		}
		if status.RemotePort != nil {
			item["remote_port"] = int(*status.RemotePort)
		}
		if status.RemoteSitePath != nil {
			item["remote_site_path"] = *status.RemoteSitePath
		}
		if status.RouteControllerPath != nil {
			item["route_controller_path"] = *status.RouteControllerPath
		}
		if status.SourceAddress != nil {
			item["source_address"] = *status.SourceAddress
		}
		if status.TimeSinceEstablished != nil {
			item["time_since_established"] = int(*status.TimeSinceEstablished)
		}
		if status.TotalInPrefixCount != nil {
			item["total_in_prefix_count"] = int(*status.TotalInPrefixCount)
		}
		if status.TotalOutPrefixCount != nil {
			item["total_out_prefix_count"] = int(*status.TotalOutPrefixCount)
		}
		if status.Type_ != nil {
			item["type"] = *status.Type_
		}
		if status.VirtualNetworkAppliancePath != nil {
			item["virtual_network_appliance_path"] = *status.VirtualNetworkAppliancePath
		}

		if status.AnnouncedCapabilities != nil {
			item["announced_capabilities"] = status.AnnouncedCapabilities
		}
		if status.NegotiatedCapability != nil {
			item["negotiated_capabilities"] = status.NegotiatedCapability
		}

		if status.AddressFamilies != nil {
			addressFamilies := make([]interface{}, 0, len(status.AddressFamilies))
			for _, af := range status.AddressFamilies {
				afItem := make(map[string]interface{})
				if af.Type_ != nil {
					afItem["type"] = *af.Type_
				}
				if af.InPrefixCount != nil {
					afItem["in_prefix_count"] = int(*af.InPrefixCount)
				}
				if af.OutPrefixCount != nil {
					afItem["out_prefix_count"] = int(*af.OutPrefixCount)
				}
				addressFamilies = append(addressFamilies, afItem)
			}
			item["address_families"] = addressFamilies
		}

		results = append(results, item)
	}

	if err := d.Set("results", results); err != nil {
		return err
	}

	return nil
}
