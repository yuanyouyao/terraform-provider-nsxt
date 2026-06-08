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

var cliRouteControllerBgpRouteTableClient = infra.NewRouteControllerBgpRouteTableClient

func dataSourceNsxtPolicyRouteControllerBgpRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpRouteTableRead,

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
			"virtual_network_appliance_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy path of virtual network appliance to retrieve BGP routes from.",
			},
			"network_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network address filter parameter.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of BGP route table entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enforcement_point_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"transport_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transport_node_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"as_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bestpath": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"community": {
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
									"extended_community": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"large_community": {
										Type:     schema.TypeString,
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
									"multipath": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"network": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nexthops": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"afi": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"scope": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"used": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"path_from": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"peer_id": {
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
									"route_origin": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"stale": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"valid": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vni": {
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
	}
}

func dataSourceNsxtPolicyRouteControllerBgpRouteTableRead(d *schema.ResourceData, m interface{}) error {
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

	vnaPath := d.Get("virtual_network_appliance_path").(string)

	var networkPrefix *string
	if val, ok := d.GetOk("network_prefix"); ok {
		s := val.(string)
		networkPrefix = &s
	}

	client := cliRouteControllerBgpRouteTableClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	listResult, err := client.List(routeControllerID, vnaPath, networkPrefix)
	if err != nil {
		return fmt.Errorf("error listing BGP route table for Route Controller %s: %v", routeControllerID, err)
	}

	d.SetId(fmt.Sprintf("%s/bgp/route-table", routeControllerID))

	results := make([]interface{}, 0, len(listResult.Results))
	for _, bgpRoute := range listResult.Results {
		item := make(map[string]interface{})

		if bgpRoute.EnforcementPointPath != nil {
			item["enforcement_point_path"] = *bgpRoute.EnforcementPointPath
		}
		if bgpRoute.GatewayPath != nil {
			item["gateway_path"] = *bgpRoute.GatewayPath
		}
		if bgpRoute.LastUpdateTimestamp != nil {
			item["last_update_timestamp"] = int(*bgpRoute.LastUpdateTimestamp)
		}
		if bgpRoute.TransportNodeId != nil {
			item["transport_node_id"] = *bgpRoute.TransportNodeId
		}
		if bgpRoute.TransportNodePath != nil {
			item["transport_node_path"] = *bgpRoute.TransportNodePath
		}

		if bgpRoute.RouteDetails != nil {
			routeDetails := make([]interface{}, 0, len(bgpRoute.RouteDetails))
			for _, rd := range bgpRoute.RouteDetails {
				rdItem := make(map[string]interface{})

				if rd.AsPath != nil {
					rdItem["as_path"] = *rd.AsPath
				}
				if rd.Bestpath != nil {
					rdItem["bestpath"] = *rd.Bestpath
				}
				if rd.Community != nil {
					rdItem["community"] = *rd.Community
				}
				if rd.Esi != nil {
					rdItem["esi"] = *rd.Esi
				}
				if rd.EthTag != nil {
					rdItem["eth_tag"] = int(*rd.EthTag)
				}
				if rd.EvpnRouteType != nil {
					rdItem["evpn_route_type"] = int(*rd.EvpnRouteType)
				}
				if rd.ExtendedCommunity != nil {
					rdItem["extended_community"] = *rd.ExtendedCommunity
				}
				if rd.LargeCommunity != nil {
					rdItem["large_community"] = *rd.LargeCommunity
				}
				if rd.LocalPref != nil {
					rdItem["local_pref"] = int(*rd.LocalPref)
				}
				if rd.Med != nil {
					rdItem["med"] = int(*rd.Med)
				}
				if rd.Multipath != nil {
					rdItem["multipath"] = *rd.Multipath
				}
				if rd.Network != nil {
					rdItem["network"] = *rd.Network
				}
				if rd.PathFrom != nil {
					rdItem["path_from"] = *rd.PathFrom
				}
				if rd.PeerId != nil {
					rdItem["peer_id"] = *rd.PeerId
				}
				if rd.Rd != nil {
					rdItem["rd"] = *rd.Rd
				}
				if rd.Rmac != nil {
					rdItem["rmac"] = *rd.Rmac
				}
				if rd.RouteOrigin != nil {
					rdItem["route_origin"] = *rd.RouteOrigin
				}
				if rd.Stale != nil {
					rdItem["stale"] = *rd.Stale
				}
				if rd.Valid != nil {
					rdItem["valid"] = *rd.Valid
				}
				if rd.Vni != nil {
					rdItem["vni"] = int(*rd.Vni)
				}
				if rd.Weight != nil {
					rdItem["weight"] = int(*rd.Weight)
				}

				if rd.Nexthops != nil {
					nexthops := make([]interface{}, 0, len(rd.Nexthops))
					for _, nh := range rd.Nexthops {
						nhItem := make(map[string]interface{})
						if nh.Afi != nil {
							nhItem["afi"] = *nh.Afi
						}
						if nh.Ip != nil {
							nhItem["ip"] = *nh.Ip
						}
						if nh.Scope != nil {
							nhItem["scope"] = *nh.Scope
						}
						if nh.Used != nil {
							nhItem["used"] = *nh.Used
						}
						nexthops = append(nexthops, nhItem)
					}
					rdItem["nexthops"] = nexthops
				}

				routeDetails = append(routeDetails, rdItem)
			}
			item["route_details"] = routeDetails
		}

		results = append(results, item)
	}

	if err := d.Set("results", results); err != nil {
		return err
	}

	return nil
}
