// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

var cliRouteControllerInterfacesClientForResource = infra.NewRouteControllerInterfacesClient

func resourceNsxtPolicyRouteControllerInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerInterfaceCreate,
		Read:   resourceNsxtPolicyRouteControllerInterfaceRead,
		Update: resourceNsxtPolicyRouteControllerInterfaceUpdate,
		Delete: resourceNsxtPolicyRouteControllerInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyRouteControllerInterfaceImporter,
		},

		Schema: map[string]*schema.Schema{
			"route_controller_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "The policy path of the route controller.",
			},
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(64, 9000),
				Description:  "MTU size.",
			},
			"urpf_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "STRICT"}, false),
				Description:  "Unicast Reverse Path Forwarding mode.",
			},
			"floating_ip_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "IP address and subnet specification for VIP floating IP address subnets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateSingleIP(),
							},
						},
						"prefix_len": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 128),
						},
					},
				},
			},
			"interface_address": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Route Controller Interface Address parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"portgroup_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DV port group identifier discovered from vCenter.",
						},
						"virtual_network_appliance_path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validatePolicyPath(),
							Description:  "Policy path for virtual network appliance.",
						},
						"interface_subnet": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "IP address and subnet specification for interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_addresses": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validateSingleIP(),
										},
									},
									"prefix_len": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 128),
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

func resourceNsxtPolicyRouteControllerInterfaceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	var routeControllerID, interfaceID string
	if strings.HasPrefix(importID, "/infra/route-controllers/") {
		routeControllerID = getResourceIDFromResourcePath(importID, "route-controllers")
		interfaceID = getResourceIDFromResourcePath(importID, "interfaces")
		if interfaceID == "" {
			segs := strings.Split(importID, "/")
			if len(segs) >= 6 {
				interfaceID = segs[5]
			}
		}
	} else {
		segs := strings.Split(importID, "/")
		if len(segs) == 2 {
			routeControllerID = segs[0]
			interfaceID = segs[1]
		}
	}

	if routeControllerID == "" || interfaceID == "" {
		return nil, fmt.Errorf("invalid import ID: %s; expected format: <route-controller-id>/<interface-id> or full policy path", importID)
	}

	d.Set("route_controller_path", "/infra/route-controllers/"+routeControllerID)
	d.Set("nsx_id", interfaceID)
	d.SetId(routeControllerID + "/" + interfaceID)
	return []*schema.ResourceData{d}, nil
}

func resourceNsxtPolicyRouteControllerInterfaceExists(routeControllerID string, interfaceID string, connector client.Connector, d *schema.ResourceData, m interface{}) (bool, error) {
	sessionContext := getSessionContext(d, m)
	client := cliRouteControllerInterfacesClientForResource(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := client.Get(routeControllerID, interfaceID)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyRouteControllerInterfaceToStruct(d *schema.ResourceData) model.RouteControllerInterface {
	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = d.Get("nsx_id").(string)
	}

	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	var mtu *int64
	if val, ok := d.GetOk("mtu"); ok {
		m := int64(val.(int))
		mtu = &m
	}

	urpfMode := d.Get("urpf_mode").(string)

	var floatingIpSubnets []model.InterfaceSubnet
	if val, ok := d.GetOk("floating_ip_subnets"); ok {
		rawList := val.([]interface{})
		for _, rawItem := range rawList {
			m := rawItem.(map[string]interface{})
			prefixLen := int64(m["prefix_len"].(int))
			var ips []string
			for _, ip := range m["ip_addresses"].([]interface{}) {
				ips = append(ips, ip.(string))
			}
			floatingIpSubnets = append(floatingIpSubnets, model.InterfaceSubnet{
				IpAddresses: ips,
				PrefixLen:   &prefixLen,
			})
		}
	}

	var interfaceAddress []model.RouteControllerInterfaceAddress
	if val, ok := d.GetOk("interface_address"); ok {
		rawList := val.([]interface{})
		for _, rawItem := range rawList {
			m := rawItem.(map[string]interface{})
			portgroupID := m["portgroup_id"].(string)
			vnaPath := m["virtual_network_appliance_path"].(string)

			var subnets []model.InterfaceSubnet
			if subVal, ok := m["interface_subnet"]; ok {
				subList := subVal.([]interface{})
				for _, subItem := range subList {
					sm := subItem.(map[string]interface{})
					prefixLen := int64(sm["prefix_len"].(int))
					var ips []string
					for _, ip := range sm["ip_addresses"].([]interface{}) {
						ips = append(ips, ip.(string))
					}
					subnets = append(subnets, model.InterfaceSubnet{
						IpAddresses: ips,
						PrefixLen:   &prefixLen,
					})
				}
			}

			interfaceAddress = append(interfaceAddress, model.RouteControllerInterfaceAddress{
				PortgroupId:                 &portgroupID,
				VirtualNetworkAppliancePath: &vnaPath,
				InterfaceSubnet:             subnets,
			})
		}
	}

	return model.RouteControllerInterface{
		DisplayName:       &displayName,
		Description:       &description,
		Tags:              tags,
		Mtu:               mtu,
		UrpfMode:          &urpfMode,
		FloatingIpSubnets: floatingIpSubnets,
		InterfaceAddress:  interfaceAddress,
	}
}

func resourceNsxtPolicyRouteControllerInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerPath := d.Get("route_controller_path").(string)
	routeControllerID := getResourceIDFromResourcePath(routeControllerPath, "route-controllers")
	if routeControllerID == "" {
		segs := strings.Split(routeControllerPath, "/")
		routeControllerID = segs[len(segs)-1]
	}

	nsxID, err := getOrGenerateID(d, m, func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicyRouteControllerInterfaceExists(routeControllerID, id, connector, d, m)
	})
	if err != nil {
		return err
	}

	client := cliRouteControllerInterfacesClientForResource(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj := resourceNsxtPolicyRouteControllerInterfaceToStruct(d)
	err = client.Patch(routeControllerID, nsxID, obj)
	if err != nil {
		return logAPIError("Error creating resource", err)
	}

	d.SetId(routeControllerID + "/" + nsxID)
	d.Set("nsx_id", nsxID)
	return resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
}

func resourceNsxtPolicyRouteControllerInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	segs := strings.Split(d.Id(), "/")
	if len(segs) != 2 {
		return fmt.Errorf("invalid resource ID: %s", d.Id())
	}
	routeControllerID := segs[0]
	interfaceID := segs[1]

	client := cliRouteControllerInterfacesClientForResource(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, interfaceID)
	if err != nil {
		if isNotFoundError(err) {
			log.Printf("[WARN] Route Controller Interface %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return logAPIError("Error reading resource", err)
	}

	d.Set("nsx_id", interfaceID)
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

func resourceNsxtPolicyRouteControllerInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	segs := strings.Split(d.Id(), "/")
	routeControllerID := segs[0]
	interfaceID := segs[1]

	client := cliRouteControllerInterfacesClientForResource(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj := resourceNsxtPolicyRouteControllerInterfaceToStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	_, err := client.Update(routeControllerID, interfaceID, obj)
	if err != nil {
		return logAPIError("Error updating resource", err)
	}

	return resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
}

func resourceNsxtPolicyRouteControllerInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	segs := strings.Split(d.Id(), "/")
	routeControllerID := segs[0]
	interfaceID := segs[1]

	client := cliRouteControllerInterfacesClientForResource(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Delete(routeControllerID, interfaceID)
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return logAPIError("Error deleting resource", err)
	}

	return nil
}
