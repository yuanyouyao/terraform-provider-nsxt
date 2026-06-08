// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyRouteControllerBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpNeighborRead,

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
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the BGP neighbor config.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the route controller BGP neighbor config.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the resource.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The revision of the resource.",
			},
			"tag": getTagsSchema(),
			"allow_as_in": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable allow_as_in option for BGP neighbor.",
			},
			"bfd_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "BFD configuration for failure detection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to enable/disable BFD configuration.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time interval between heartbeat packets in milliseconds.",
						},
						"multiple": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of times heartbeat packet is missed before BFD declares the neighbor is down.",
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable/disable BGP peering.",
			},
			"gateway_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Next hop gateway IP addresses to reach non-directly connected BGP peers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"graceful_restart_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BGP Graceful Restart Configuration Mode.",
			},
			"hold_down_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Wait time in seconds before declaring peer dead.",
			},
			"keep_alive_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Interval in seconds between keep alive messages sent to peer.",
			},
			"maximum_hop_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of hops allowed to reach BGP neighbor.",
			},
			"neighbor_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Neighbor IP Address.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password for BGP neighbor authentication.",
				Sensitive:   true,
			},
			"remote_as_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ASN of the neighbor in ASPLAIN format.",
			},
			"source_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source IP addresses for BGP peering.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"route_filtering": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Enable address families and route filtering in each direction.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address family type.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to enable/disable address family.",
						},
						"in_route_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"maximum_routes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of routes for the address family.",
						},
						"out_route_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
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

	id := d.Get("id").(string)

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, id)
	if err != nil {
		return fmt.Errorf("error getting route controller BGP neighbor for Route Controller %s and ID %s: %v", routeControllerID, id, err)
	}

	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)

	if obj.AllowAsIn != nil {
		d.Set("allow_as_in", *obj.AllowAsIn)
	}

	if obj.Bfd != nil {
		bfdMap := map[string]interface{}{
			"enabled":  *obj.Bfd.Enabled,
			"interval": int(*obj.Bfd.Interval),
			"multiple": int(*obj.Bfd.Multiple),
		}
		d.Set("bfd_config", []interface{}{bfdMap})
	} else {
		d.Set("bfd_config", nil)
	}

	if obj.Enabled != nil {
		d.Set("enabled", *obj.Enabled)
	}

	if len(obj.GatewayIps) > 0 {
		var ips []interface{}
		for _, ip := range obj.GatewayIps {
			ips = append(ips, ip)
		}
		d.Set("gateway_ips", ips)
	} else {
		d.Set("gateway_ips", nil)
	}

	if obj.GracefulRestartMode != nil {
		d.Set("graceful_restart_mode", *obj.GracefulRestartMode)
	}

	if obj.HoldDownTime != nil {
		d.Set("hold_down_time", int(*obj.HoldDownTime))
	}

	if obj.KeepAliveTime != nil {
		d.Set("keep_alive_time", int(*obj.KeepAliveTime))
	}

	if obj.MaximumHopLimit != nil {
		d.Set("maximum_hop_limit", int(*obj.MaximumHopLimit))
	}

	if obj.NeighborAddress != nil {
		d.Set("neighbor_address", *obj.NeighborAddress)
	}

	if obj.Password != nil {
		d.Set("password", *obj.Password)
	}

	if obj.RemoteAsNum != nil {
		d.Set("remote_as_num", *obj.RemoteAsNum)
	}

	if len(obj.SourceAddresses) > 0 {
		var ips []interface{}
		for _, ip := range obj.SourceAddresses {
			ips = append(ips, ip)
		}
		d.Set("source_addresses", ips)
	} else {
		d.Set("source_addresses", nil)
	}

	if len(obj.RouteFiltering) > 0 {
		var rfList []interface{}
		for _, rf := range obj.RouteFiltering {
			rfMap := map[string]interface{}{
				"address_family": *rf.AddressFamily,
				"enabled":        *rf.Enabled,
			}
			if rf.MaximumRoutes != nil {
				rfMap["maximum_routes"] = int(*rf.MaximumRoutes)
			}
			if len(rf.InRouteFilters) > 0 {
				var inFilters []interface{}
				for _, filter := range rf.InRouteFilters {
					inFilters = append(inFilters, filter)
				}
				rfMap["in_route_filters"] = inFilters
			}
			if len(rf.OutRouteFilters) > 0 {
				var outFilters []interface{}
				for _, filter := range rf.OutRouteFilters {
					outFilters = append(outFilters, filter)
				}
				rfMap["out_route_filters"] = outFilters
			}
			rfList = append(rfList, rfMap)
		}
		d.Set("route_filtering", rfList)
	} else {
		d.Set("route_filtering", nil)
	}

	d.SetId(routeControllerID + "/" + id)
	return nil
}
