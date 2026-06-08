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

var cliRouteControllerRoutingTableClient = infra.NewRouteControllerRoutingTableClient

func dataSourceNsxtPolicyRouteControllerRoutingTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerRoutingTableRead,

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
			"network_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network address filter parameter.",
			},
			"route_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter routes based on the source from which route is learned.",
			},
			"virtual_network_appliance_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contains string path of virtual network appliance.",
			},
			"routing_table": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routing table entries from the route controller.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of route entries.",
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error message if any.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the routing table retrieval.",
						},
						"transport_node_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contains string path of transport node.",
						},
						"virtual_network_appliance_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contains string path of virtual network appliance.",
						},
						"route_entries": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of routing entries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"admin_distance": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Admin distance.",
									},
									"black_hole": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Value of this field will be true if given routes are null routes.",
									},
									"lr_component_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Logical router component ID.",
									},
									"lr_component_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Logical router component type.",
									},
									"network": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network CIDR.",
									},
									"next_hop": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Next hop address.",
									},
									"next_hop_gateway": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Next hop gateway path.",
									},
									"route_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route type in routing table.",
									},
								},
							},
						},
					},
				},
			},
			"routing_table_csv": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The routing table in CSV format downloaded from the route controller.",
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerRoutingTableRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	clients := m.(nsxtClients)

	username := clients.CommonConfig.Username
	password := clients.CommonConfig.Password
	token := clients.CommonConfig.BearerToken

	var cookie, xsrf string
	if clients.NsxtClientConfig != nil {
		if len(clients.NsxtClientConfig.DefaultHeader["Cookie"]) > 0 {
			cookie = clients.NsxtClientConfig.DefaultHeader["Cookie"]
		}
		if len(clients.NsxtClientConfig.DefaultHeader["X-XSRF-TOKEN"]) > 0 {
			xsrf = clients.NsxtClientConfig.DefaultHeader["X-XSRF-TOKEN"]
		}
	}

	client := cliRouteControllerRoutingTableClient(sessionContext, connector, clients.Host, clients.PolicyHTTPClient, username, password, token, cookie, xsrf)
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

	var networkPrefix *string
	if val, ok := d.GetOk("network_prefix"); ok {
		s := val.(string)
		networkPrefix = &s
	}

	var routeSource *string
	if val, ok := d.GetOk("route_source"); ok {
		s := val.(string)
		routeSource = &s
	}

	var vnaPath *string
	if val, ok := d.GetOk("virtual_network_appliance_path"); ok {
		s := val.(string)
		vnaPath = &s
	}

	res, err := client.Get(routeControllerID, networkPrefix, routeSource, vnaPath)
	if err != nil {
		return fmt.Errorf("error getting route controller routing table for ID %s: %v", routeControllerID, err)
	}

	routingTables := make([]map[string]interface{}, 0, len(res.Results))
	for _, rt := range res.Results {
		item := make(map[string]interface{})
		if rt.Count != nil {
			item["count"] = int(*rt.Count)
		}
		if rt.ErrorMessage != nil {
			item["error_message"] = *rt.ErrorMessage
		}
		if rt.Status != nil {
			item["status"] = *rt.Status
		}
		if rt.TransportNodePath != nil {
			item["transport_node_path"] = *rt.TransportNodePath
		}
		if rt.VirtualNetworkAppliancePath != nil {
			item["virtual_network_appliance_path"] = *rt.VirtualNetworkAppliancePath
		}

		routeEntries := make([]map[string]interface{}, 0, len(rt.RouteEntries))
		for _, re := range rt.RouteEntries {
			entry := make(map[string]interface{})
			if re.AdminDistance != nil {
				entry["admin_distance"] = int(*re.AdminDistance)
			}
			if re.BlackHole != nil {
				entry["black_hole"] = *re.BlackHole
			}
			if re.LrComponentId != nil {
				entry["lr_component_id"] = *re.LrComponentId
			}
			if re.LrComponentType != nil {
				entry["lr_component_type"] = *re.LrComponentType
			}
			if re.Network != nil {
				entry["network"] = *re.Network
			}
			if re.NextHop != nil {
				entry["next_hop"] = *re.NextHop
			}
			if re.NextHopGateway != nil {
				entry["next_hop_gateway"] = *re.NextHopGateway
			}
			if re.RouteType != nil {
				entry["route_type"] = *re.RouteType
			}
			routeEntries = append(routeEntries, entry)
		}
		item["route_entries"] = routeEntries
		routingTables = append(routingTables, item)
	}
	d.Set("routing_table", routingTables)

	csvData, err := client.Download(routeControllerID, networkPrefix, routeSource, vnaPath)
	if err != nil {
		return fmt.Errorf("error downloading route controller routing table for ID %s: %v", routeControllerID, err)
	}
	d.Set("routing_table_csv", csvData)

	d.SetId(routeControllerID)
	return nil
}
