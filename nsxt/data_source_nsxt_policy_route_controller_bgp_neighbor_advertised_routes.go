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

var cliRouteControllerBgpNeighborAdvertisedRoutesClient = infra.NewRouteControllerBgpNeighborAdvertisedRoutesClient

func dataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutesRead,

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
			"neighbor_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"neighbor_path", "neighbor_id"},
				Description:  "The policy path of the BGP neighbor.",
			},
			"neighbor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"neighbor_path", "neighbor_id"},
				Description:  "The ID of the BGP neighbor.",
			},
			"enforcement_point_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enforcement point path.",
			},
			"neighbor_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dynamically discovered BGP neighbor address.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of advertised routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enforcement_point_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"neighbor_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_appliance_routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"transport_node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"routes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"as_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"esi": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"eth_tag": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"evpn_route_type": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"local_pref": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"med": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"network": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"next_hop": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rd": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rmac": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rmac_len": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"weight": {
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
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutesRead(d *schema.ResourceData, m interface{}) error {
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

	neighborID := d.Get("neighbor_id").(string)
	if neighborID == "" {
		path := d.Get("neighbor_path").(string)
		neighborID = getResourceIDFromResourcePath(path, "neighbors")
		if neighborID == "" {
			segs := strings.Split(path, "/")
			neighborID = segs[len(segs)-1]
		}
	}

	client := cliRouteControllerBgpNeighborAdvertisedRoutesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	var enforcementPointPath *string
	if val, ok := d.GetOk("enforcement_point_path"); ok {
		s := val.(string)
		enforcementPointPath = &s
	}

	var neighborAddress *string
	if val, ok := d.GetOk("neighbor_address"); ok {
		s := val.(string)
		neighborAddress = &s
	}

	listResult, err := client.List(routeControllerID, neighborID, nil, nil, enforcementPointPath, nil, neighborAddress, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error listing BGP neighbor advertised routes for Route Controller %s neighbor %s: %v", routeControllerID, neighborID, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", routeControllerID, neighborID))

	results := make([]interface{}, 0, len(listResult.Results))
	for _, routes := range listResult.Results {
		item := make(map[string]interface{})

		if routes.EnforcementPointPath != nil {
			item["enforcement_point_path"] = *routes.EnforcementPointPath
		}
		if routes.NeighborPath != nil {
			item["neighbor_path"] = *routes.NeighborPath
		}

		if routes.VirtualNetworkApplianceRoutes != nil {
			vnaRoutes := make([]interface{}, 0, len(routes.VirtualNetworkApplianceRoutes))
			for _, vna := range routes.VirtualNetworkApplianceRoutes {
				vnaItem := make(map[string]interface{})
				if vna.SourceAddress != nil {
					vnaItem["source_address"] = *vna.SourceAddress
				}
				if vna.TransportNodeId != nil {
					vnaItem["transport_node_id"] = *vna.TransportNodeId
				}

				if vna.Routes != nil {
					routeDetails := make([]interface{}, 0, len(vna.Routes))
					for _, route := range vna.Routes {
						routeItem := make(map[string]interface{})
						if route.AsPath != nil {
							routeItem["as_path"] = *route.AsPath
						}
						if route.Esi != nil {
							routeItem["esi"] = *route.Esi
						}
						if route.EthTag != nil {
							routeItem["eth_tag"] = int(*route.EthTag)
						}
						if route.EvpnRouteType != nil {
							routeItem["evpn_route_type"] = int(*route.EvpnRouteType)
						}
						if route.LocalPref != nil {
							routeItem["local_pref"] = int(*route.LocalPref)
						}
						if route.Med != nil {
							routeItem["med"] = int(*route.Med)
						}
						if route.Network != nil {
							routeItem["network"] = *route.Network
						}
						if route.NextHop != nil {
							routeItem["next_hop"] = *route.NextHop
						}
						if route.Rd != nil {
							routeItem["rd"] = *route.Rd
						}
						if route.Rmac != nil {
							routeItem["rmac"] = *route.Rmac
						}
						if route.RmacLen != nil {
							routeItem["rmac_len"] = int(*route.RmacLen)
						}
						if route.Weight != nil {
							routeItem["weight"] = int(*route.Weight)
						}
						routeDetails = append(routeDetails, routeItem)
					}
					vnaItem["routes"] = routeDetails
				}
				vnaRoutes = append(vnaRoutes, vnaItem)
			}
			item["virtual_network_appliance_routes"] = vnaRoutes
		}

		results = append(results, item)
	}

	if err := d.Set("results", results); err != nil {
		return err
	}

	return nil
}
