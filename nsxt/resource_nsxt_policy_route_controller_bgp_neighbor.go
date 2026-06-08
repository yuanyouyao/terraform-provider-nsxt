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

var cliRouteControllerBgpNeighborClient = infra.NewRouteControllerBgpNeighborClient

func resourceNsxtPolicyRouteControllerBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerBgpNeighborCreate,
		Read:   resourceNsxtPolicyRouteControllerBgpNeighborRead,
		Update: resourceNsxtPolicyRouteControllerBgpNeighborUpdate,
		Delete: resourceNsxtPolicyRouteControllerBgpNeighborDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyRouteControllerBgpNeighborImporter,
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
			"allow_as_in": {
				Description: "Flag to enable allow_as_in option for BGP neighbor.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"bfd_config": {
				Type:        schema.TypeList,
				Description: "BFD configuration for failure detection.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to enable/disable BFD configuration.",
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      500,
							Description:  "Time interval between heartbeat packets in milliseconds.",
							ValidateFunc: validation.IntBetween(50, 60000),
						},
						"multiple": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							Description:  "Number of times heartbeat packet is missed before BFD declares the neighbor is down.",
							ValidateFunc: validation.IntBetween(2, 16),
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to enable/disable BGP peering.",
			},
			"gateway_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Next hop gateway IP addresses to reach non-directly connected BGP peers.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			"graceful_restart_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      model.RouteControllerBgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY,
				ValidateFunc: validation.StringInSlice([]string{model.RouteControllerBgpNeighborConfig_GRACEFUL_RESTART_MODE_DISABLE, model.RouteControllerBgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY}, false),
				Description:  "BGP Graceful Restart Configuration Mode.",
			},
			"hold_down_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      180,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Wait time in seconds before declaring peer dead.",
			},
			"keep_alive_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Interval in seconds between keep alive messages sent to peer.",
			},
			"maximum_hop_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 255),
				Description:  "Maximum number of hops allowed to reach BGP neighbor.",
			},
			"neighbor_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Neighbor IP Address.",
				ValidateFunc: validateSingleIP(),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Password for BGP neighbor authentication.",
				ValidateFunc: validation.StringLenBetween(0, 32),
				Sensitive:    true,
			},
			"remote_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ASN of the neighbor in ASPLAIN format.",
				ValidateFunc: validateASPlainOrDot,
			},
			"source_addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Source IP addresses for BGP peering.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			"route_filtering": {
				Type:        schema.TypeList,
				Description: "Enable address families and route filtering in each direction.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Address family type.",
							ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "L2VPN_EVPN"}, false),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Flag to enable/disable address family.",
						},
						"in_route_filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Prefix-list or route map paths for IN direction.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validatePolicyPath(),
							},
						},
						"maximum_routes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1000000),
							Description:  "Maximum number of routes for the address family.",
						},
						"out_route_filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Prefix-list or route map paths for OUT direction.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validatePolicyPath(),
							},
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyRouteControllerBgpNeighborImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	var routeControllerID, neighborID string
	if strings.HasPrefix(importID, "/infra/route-controllers/") {
		routeControllerID = getResourceIDFromResourcePath(importID, "route-controllers")
		neighborID = getResourceIDFromResourcePath(importID, "neighbors")
		if neighborID == "" {
			segs := strings.Split(importID, "/")
			if len(segs) >= 6 {
				neighborID = segs[5]
			}
		}
	} else {
		segs := strings.Split(importID, "/")
		if len(segs) == 2 {
			routeControllerID = segs[0]
			neighborID = segs[1]
		}
	}

	if routeControllerID == "" || neighborID == "" {
		return nil, fmt.Errorf("invalid import ID: %s; expected format: <route-controller-id>/<neighbor-id> or full policy path", importID)
	}

	d.Set("route_controller_path", "/infra/route-controllers/"+routeControllerID)
	d.Set("nsx_id", neighborID)
	d.SetId(routeControllerID + "/" + neighborID)
	return []*schema.ResourceData{d}, nil
}

func resourceNsxtPolicyRouteControllerBgpNeighborExists(routeControllerID string, neighborID string, connector client.Connector, d *schema.ResourceData, m interface{}) (bool, error) {
	sessionContext := getSessionContext(d, m)
	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := client.Get(routeControllerID, neighborID)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyRouteControllerBgpNeighborToStruct(d *schema.ResourceData) model.RouteControllerBgpNeighborConfig {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.RouteControllerBgpNeighborConfig{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	if val, ok := d.GetOkExists("allow_as_in"); ok {
		b := val.(bool)
		obj.AllowAsIn = &b
	}

	if bfdVal, ok := d.GetOk("bfd_config"); ok && len(bfdVal.([]interface{})) > 0 {
		bfdMap := bfdVal.([]interface{})[0].(map[string]interface{})
		enabled := bfdMap["enabled"].(bool)
		interval := int64(bfdMap["interval"].(int))
		multiple := int64(bfdMap["multiple"].(int))
		obj.Bfd = &model.BgpBfdConfig{
			Enabled:  &enabled,
			Interval: &interval,
			Multiple: &multiple,
		}
	}

	if val, ok := d.GetOkExists("enabled"); ok {
		b := val.(bool)
		obj.Enabled = &b
	}

	if val, ok := d.GetOk("gateway_ips"); ok {
		var ips []string
		for _, ip := range val.([]interface{}) {
			ips = append(ips, ip.(string))
		}
		obj.GatewayIps = ips
	}

	if val, ok := d.GetOkExists("graceful_restart_mode"); ok {
		s := val.(string)
		obj.GracefulRestartMode = &s
	}

	if val, ok := d.GetOkExists("hold_down_time"); ok {
		i := int64(val.(int))
		obj.HoldDownTime = &i
	}

	if val, ok := d.GetOkExists("keep_alive_time"); ok {
		i := int64(val.(int))
		obj.KeepAliveTime = &i
	}

	if val, ok := d.GetOkExists("maximum_hop_limit"); ok {
		i := int64(val.(int))
		obj.MaximumHopLimit = &i
	}

	if val, ok := d.GetOkExists("neighbor_address"); ok {
		s := val.(string)
		obj.NeighborAddress = &s
	}

	if val, ok := d.GetOkExists("password"); ok {
		s := val.(string)
		obj.Password = &s
	}

	if val, ok := d.GetOkExists("remote_as_num"); ok {
		s := val.(string)
		obj.RemoteAsNum = &s
	}

	if val, ok := d.GetOk("source_addresses"); ok {
		var ips []string
		for _, ip := range val.([]interface{}) {
			ips = append(ips, ip.(string))
		}
		obj.SourceAddresses = ips
	}

	if rfVal, ok := d.GetOk("route_filtering"); ok {
		var routeFiltering []model.BgpRouteFiltering
		for _, rfItem := range rfVal.([]interface{}) {
			rfMap := rfItem.(map[string]interface{})
			family := rfMap["address_family"].(string)
			enabled := rfMap["enabled"].(bool)

			var inFilters []string
			if inVal, ok := rfMap["in_route_filters"]; ok {
				for _, filter := range inVal.([]interface{}) {
					inFilters = append(inFilters, filter.(string))
				}
			}

			var outFilters []string
			if outVal, ok := rfMap["out_route_filters"]; ok {
				for _, filter := range outVal.([]interface{}) {
					outFilters = append(outFilters, filter.(string))
				}
			}

			rf := model.BgpRouteFiltering{
				AddressFamily:   &family,
				Enabled:         &enabled,
				InRouteFilters:  inFilters,
				OutRouteFilters: outFilters,
			}

			if maxRoutesVal, ok := rfMap["maximum_routes"]; ok && maxRoutesVal.(int) > 0 {
				maxRoutes := int64(maxRoutesVal.(int))
				rf.MaximumRoutes = &maxRoutes
			}

			routeFiltering = append(routeFiltering, rf)
		}
		obj.RouteFiltering = routeFiltering
	}

	return obj
}

func resourceNsxtPolicyRouteControllerBgpNeighborCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerPath := d.Get("route_controller_path").(string)
	routeControllerID := getResourceIDFromResourcePath(routeControllerPath, "route-controllers")
	if routeControllerID == "" {
		segs := strings.Split(routeControllerPath, "/")
		routeControllerID = segs[len(segs)-1]
	}

	id, err := getOrGenerateID(d, m, func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicyRouteControllerBgpNeighborExists(routeControllerID, id, connector, d, m)
	})
	if err != nil {
		return err
	}

	obj := resourceNsxtPolicyRouteControllerBgpNeighborToStruct(d)

	log.Printf("[INFO] Creating Route Controller BGP Neighbor %s for Route Controller %s", id, routeControllerID)

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err = client.Patch(routeControllerID, id, obj)
	if err != nil {
		return handleCreateError("Route Controller BGP Neighbor", id, err)
	}

	d.SetId(routeControllerID + "/" + id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	idStr := d.Id()
	segs := strings.Split(idStr, "/")
	if len(segs) < 2 {
		return fmt.Errorf("invalid resource ID format: %s", idStr)
	}
	routeControllerID := segs[0]
	id := segs[1]

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, id)
	if err != nil {
		return handleReadError(d, "Route Controller BGP Neighbor", id, err)
	}

	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("nsx_id", id)
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

	return nil
}

func resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	idStr := d.Id()
	segs := strings.Split(idStr, "/")
	if len(segs) < 2 {
		return fmt.Errorf("invalid resource ID format: %s", idStr)
	}
	routeControllerID := segs[0]
	id := segs[1]

	obj := resourceNsxtPolicyRouteControllerBgpNeighborToStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating Route Controller BGP Neighbor %s for Route Controller %s", id, routeControllerID)

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	_, err := client.Update(routeControllerID, id, obj)
	if err != nil {
		return handleUpdateError("Route Controller BGP Neighbor", id, err)
	}

	return resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpNeighborDelete(d *schema.ResourceData, m interface{}) error {
	idStr := d.Id()
	segs := strings.Split(idStr, "/")
	if len(segs) < 2 {
		return fmt.Errorf("invalid resource ID format: %s", idStr)
	}
	routeControllerID := segs[0]
	id := segs[1]

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	log.Printf("[INFO] Deleting Route Controller BGP Neighbor %s for Route Controller %s", id, routeControllerID)

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Delete(routeControllerID, id)
	if err != nil {
		return handleDeleteError("Route Controller BGP Neighbor", id, err)
	}

	return nil
}
