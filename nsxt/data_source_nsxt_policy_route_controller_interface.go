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

var cliRouteControllerInterfacesClientForDataSource = infra.NewRouteControllerInterfacesClient

func dataSourceNsxtPolicyRouteControllerInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerInterfaceRead,

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
				Description: "The ID of the route controller interface.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the route controller interface.",
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
			"mtu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "MTU size.",
			},
			"urpf_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unicast Reverse Path Forwarding mode.",
			},
			"floating_ip_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP address and subnet specification for VIP floating IP address subnets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"prefix_len": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"interface_address": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Route Controller Interface Address parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"portgroup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_appliance_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface_subnet": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"prefix_len": {
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

func dataSourceNsxtPolicyRouteControllerInterfaceRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerInterfacesClientForDataSource(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, id)
	if err != nil {
		return fmt.Errorf("error getting route controller interface for Route Controller %s and ID %s: %v", routeControllerID, id, err)
	}

	d.SetId(routeControllerID + "/" + id)

	if obj.Path != nil {
		d.Set("path", *obj.Path)
	}
	if obj.DisplayName != nil {
		d.Set("display_name", *obj.DisplayName)
	}
	if obj.Description != nil {
		d.Set("description", *obj.Description)
	}
	if obj.Revision != nil {
		d.Set("revision", int(*obj.Revision))
	}
	setPolicyTagsInSchema(d, obj.Tags)

	if obj.Mtu != nil {
		d.Set("mtu", int(*obj.Mtu))
	}
	if obj.UrpfMode != nil {
		d.Set("urpf_mode", *obj.UrpfMode)
	}

	if len(obj.FloatingIpSubnets) > 0 {
		var subnets []interface{}
		for _, sub := range obj.FloatingIpSubnets {
			subItem := make(map[string]interface{})
			if sub.PrefixLen != nil {
				subItem["prefix_len"] = int(*sub.PrefixLen)
			}
			if len(sub.IpAddresses) > 0 {
				var ips []interface{}
				for _, ip := range sub.IpAddresses {
					ips = append(ips, ip)
				}
				subItem["ip_addresses"] = ips
			}
			subnets = append(subnets, subItem)
		}
		d.Set("floating_ip_subnets", subnets)
	} else {
		d.Set("floating_ip_subnets", nil)
	}

	if len(obj.InterfaceAddress) > 0 {
		var addrs []interface{}
		for _, addr := range obj.InterfaceAddress {
			addrItem := make(map[string]interface{})
			if addr.PortgroupId != nil {
				addrItem["portgroup_id"] = *addr.PortgroupId
			}
			if addr.VirtualNetworkAppliancePath != nil {
				addrItem["virtual_network_appliance_path"] = *addr.VirtualNetworkAppliancePath
			}

			if len(addr.InterfaceSubnet) > 0 {
				var subnets []interface{}
				for _, sub := range addr.InterfaceSubnet {
					subItem := make(map[string]interface{})
					if sub.PrefixLen != nil {
						subItem["prefix_len"] = int(*sub.PrefixLen)
					}
					if len(sub.IpAddresses) > 0 {
						var ips []interface{}
						for _, ip := range sub.IpAddresses {
							ips = append(ips, ip)
						}
						subItem["ip_addresses"] = ips
					}
					subnets = append(subnets, subItem)
				}
				addrItem["interface_subnet"] = subnets
			}
			addrs = append(addrs, addrItem)
		}
		d.Set("interface_address", addrs)
	} else {
		d.Set("interface_address", nil)
	}

	return nil
}
