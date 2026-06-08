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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

var cliRouteControllerBgpClient = infra.NewRouteControllerBgpClient

func resourceNsxtPolicyRouteControllerBgp() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerBgpCreate,
		Read:   resourceNsxtPolicyRouteControllerBgpRead,
		Update: resourceNsxtPolicyRouteControllerBgpUpdate,
		Delete: resourceNsxtPolicyRouteControllerBgpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyRouteControllerBgpImporter,
		},

		Schema: map[string]*schema.Schema{
			"route_controller_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "The policy path of the route controller.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the route controller BGP routing config.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The revision of the resource.",
			},
			"tag": getTagsSchema(),
			"ecmp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to enable ECMP.",
			},
			"local_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateASPlainOrDot,
				Description:  "BGP AS number in ASPLAIN/ASDOT format.",
			},
			"multipath_relax": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Flag to enable BGP multipath relax option.",
			},
			"peer_route_convergence_timer": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Extra time in seconds the router must wait before sending the UP notification.",
			},
			"graceful_restart_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
				ValidateFunc: validation.StringInSlice([]string{model.BgpGracefulRestartConfig_MODE_DISABLE, model.BgpGracefulRestartConfig_MODE_GR_AND_HELPER, model.BgpGracefulRestartConfig_MODE_HELPER_ONLY}, false),
				Description:  "BGP Graceful Restart Configuration Mode.",
			},
			"graceful_restart_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      120,
				ValidateFunc: validation.IntBetween(1, 3600),
				Description:  "Maximum time taken (in seconds) for a BGP session to be established after a restart.",
			},
			"graceful_restart_stale_route_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validation.IntBetween(1, 3600),
				Description:  "Maximum time (in seconds) before stale routes are removed from the RIB when BGP restarts.",
			},
		},
	}
}

func resourceNsxtPolicyRouteControllerBgpImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	var routeControllerID string
	if strings.HasPrefix(importID, "/infra/route-controllers/") {
		routeControllerID = getResourceIDFromResourcePath(importID, "route-controllers")
		if routeControllerID == "" {
			segs := strings.Split(importID, "/")
			if len(segs) >= 4 {
				routeControllerID = segs[3]
			}
		}
	} else {
		routeControllerID = importID
	}

	if routeControllerID == "" {
		return nil, fmt.Errorf("invalid import ID: %s", importID)
	}

	d.Set("route_controller_path", "/infra/route-controllers/"+routeControllerID)
	d.SetId(routeControllerID)
	return []*schema.ResourceData{d}, nil
}

func resourceNsxtPolicyRouteControllerBgpToStruct(d *schema.ResourceData) model.RouteControllerBgpRoutingConfig {
	tags := getPolicyTagsFromSchema(d)

	obj := model.RouteControllerBgpRoutingConfig{
		Tags: tags,
	}

	if val, ok := d.GetOkExists("ecmp"); ok {
		b := val.(bool)
		obj.Ecmp = &b
	}
	if val, ok := d.GetOkExists("local_as_num"); ok {
		s := val.(string)
		obj.LocalAsNum = &s
	}
	if val, ok := d.GetOkExists("multipath_relax"); ok {
		b := val.(bool)
		obj.MultipathRelax = &b
	}
	if val, ok := d.GetOkExists("peer_route_convergence_timer"); ok {
		i := int64(val.(int))
		obj.PeerRouteConvergenceTimer = &i
	}

	mode := d.Get("graceful_restart_mode").(string)
	restartTimer := int64(d.Get("graceful_restart_timer").(int))
	staleTimer := int64(d.Get("graceful_restart_stale_route_timer").(int))

	timerStruct := model.BgpGracefulRestartTimer{
		RestartTimer:    &restartTimer,
		StaleRouteTimer: &staleTimer,
	}
	restartConfigStruct := model.BgpGracefulRestartConfig{
		Mode:  &mode,
		Timer: &timerStruct,
	}
	obj.GracefulRestartConfig = &restartConfigStruct

	return obj
}

func resourceNsxtPolicyRouteControllerBgpCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerPath := d.Get("route_controller_path").(string)
	routeControllerID := getResourceIDFromResourcePath(routeControllerPath, "route-controllers")
	if routeControllerID == "" {
		segs := strings.Split(routeControllerPath, "/")
		routeControllerID = segs[len(segs)-1]
	}

	obj := resourceNsxtPolicyRouteControllerBgpToStruct(d)

	log.Printf("[INFO] Creating Route Controller BGP Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Patch(routeControllerID, obj)
	if err != nil {
		return handleCreateError("Route Controller BGP Config", routeControllerID, err)
	}

	d.SetId(routeControllerID)
	return resourceNsxtPolicyRouteControllerBgpRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP ID")
	}

	client := cliRouteControllerBgpClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID)
	if err != nil {
		return handleReadError(d, "Route Controller BGP Config", routeControllerID, err)
	}

	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setPolicyTagsInSchema(d, obj.Tags)

	if obj.Ecmp != nil {
		d.Set("ecmp", *obj.Ecmp)
	}
	if obj.LocalAsNum != nil {
		d.Set("local_as_num", *obj.LocalAsNum)
	}
	if obj.MultipathRelax != nil {
		d.Set("multipath_relax", *obj.MultipathRelax)
	}
	if obj.PeerRouteConvergenceTimer != nil {
		d.Set("peer_route_convergence_timer", int(*obj.PeerRouteConvergenceTimer))
	}

	if obj.GracefulRestartConfig != nil {
		if obj.GracefulRestartConfig.Mode != nil {
			d.Set("graceful_restart_mode", *obj.GracefulRestartConfig.Mode)
		}
		if obj.GracefulRestartConfig.Timer != nil {
			if obj.GracefulRestartConfig.Timer.RestartTimer != nil {
				d.Set("graceful_restart_timer", int(*obj.GracefulRestartConfig.Timer.RestartTimer))
			}
			if obj.GracefulRestartConfig.Timer.StaleRouteTimer != nil {
				d.Set("graceful_restart_stale_route_timer", int(*obj.GracefulRestartConfig.Timer.StaleRouteTimer))
			}
		}
	}

	return nil
}

func resourceNsxtPolicyRouteControllerBgpUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP ID")
	}

	obj := resourceNsxtPolicyRouteControllerBgpToStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating Route Controller BGP Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	_, err := client.Update(routeControllerID, obj)
	if err != nil {
		return handleUpdateError("Route Controller BGP Config", routeControllerID, err)
	}

	return resourceNsxtPolicyRouteControllerBgpRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpDelete(d *schema.ResourceData, m interface{}) error {
	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	log.Printf("[INFO] Deleting Route Controller BGP Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Delete(routeControllerID)
	if err != nil {
		return handleDeleteError("Route Controller BGP Config", routeControllerID, err)
	}

	return nil
}
