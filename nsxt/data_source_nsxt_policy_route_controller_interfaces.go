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

var cliRouteControllerInterfacesClient = infra.NewRouteControllerInterfacesClient

func dataSourceNsxtPolicyRouteControllerInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerInterfacesRead,

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
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of route controller interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revision": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mtu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"urpf_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip_subnets": {
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
						"interface_address": {
							Type:     schema.TypeList,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerInterfacesRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerInterfacesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	listResult, err := client.List(routeControllerID, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error listing interfaces for Route Controller %s: %v", routeControllerID, err)
	}

	d.SetId(routeControllerID)

	results := make([]interface{}, 0, len(listResult.Results))
	for _, iface := range listResult.Results {
		item := make(map[string]interface{})

		if iface.Id != nil {
			item["id"] = *iface.Id
		}
		if iface.DisplayName != nil {
			item["display_name"] = *iface.DisplayName
		}
		if iface.Description != nil {
			item["description"] = *iface.Description
		}
		if iface.Path != nil {
			item["path"] = *iface.Path
		}
		if iface.Revision != nil {
			item["revision"] = int(*iface.Revision)
		}
		if iface.Mtu != nil {
			item["mtu"] = int(*iface.Mtu)
		}
		if iface.UrpfMode != nil {
			item["urpf_mode"] = *iface.UrpfMode
		}

		if iface.FloatingIpSubnets != nil {
			subnets := make([]interface{}, 0, len(iface.FloatingIpSubnets))
			for _, sub := range iface.FloatingIpSubnets {
				subItem := make(map[string]interface{})
				if sub.PrefixLen != nil {
					subItem["prefix_len"] = int(*sub.PrefixLen)
				}
				if sub.IpAddresses != nil {
					ips := make([]interface{}, 0, len(sub.IpAddresses))
					for _, ip := range sub.IpAddresses {
						ips = append(ips, ip)
					}
					subItem["ip_addresses"] = ips
				}
				subnets = append(subnets, subItem)
			}
			item["floating_ip_subnets"] = subnets
		}

		if iface.InterfaceAddress != nil {
			addrs := make([]interface{}, 0, len(iface.InterfaceAddress))
			for _, addr := range iface.InterfaceAddress {
				addrItem := make(map[string]interface{})
				if addr.PortgroupId != nil {
					addrItem["portgroup_id"] = *addr.PortgroupId
				}
				if addr.VirtualNetworkAppliancePath != nil {
					addrItem["virtual_network_appliance_path"] = *addr.VirtualNetworkAppliancePath
				}

				if addr.InterfaceSubnet != nil {
					subnets := make([]interface{}, 0, len(addr.InterfaceSubnet))
					for _, sub := range addr.InterfaceSubnet {
						subItem := make(map[string]interface{})
						if sub.PrefixLen != nil {
							subItem["prefix_len"] = int(*sub.PrefixLen)
						}
						if sub.IpAddresses != nil {
							ips := make([]interface{}, 0, len(sub.IpAddresses))
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
			item["interface_address"] = addrs
		}

		results = append(results, item)
	}

	if err := d.Set("results", results); err != nil {
		return err
	}

	return nil
}
