// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

var cliRouteControllerBgpTroubleshootClient = infra.NewRouteControllerBgpTroubleshootClient

func resourceNsxtPolicyRouteControllerBgpTroubleshoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerBgpTroubleshootCreate,
		Read:   resourceNsxtPolicyRouteControllerBgpTroubleshootRead,
		Update: resourceNsxtPolicyRouteControllerBgpTroubleshootUpdate,
		Delete: resourceNsxtPolicyRouteControllerBgpTroubleshootDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyRouteControllerBgpTroubleshootImporter,
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
				Description: "The policy path of the route controller BGP troubleshoot config.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The revision of the resource.",
			},
			"tag": getTagsSchema(),
			"bfd_control_pkt_diagnostics": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to enable/disable the collection of the timestamps of sent and received BFD control messages per BFD peer session.",
			},
			"bgp_session_diagnostics": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to enable/disable the collection of the timestamps of sent and received Keep-Alive messages per BGP peer session, and the session states.",
			},
			"system_diagnostics": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to enable/disable the collection of system diagnostic data such as ARP, Ping, CPU stats, etc.",
			},
		},
	}
}

func resourceNsxtPolicyRouteControllerBgpTroubleshootImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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

func resourceNsxtPolicyRouteControllerBgpTroubleshootToStruct(d *schema.ResourceData) model.BgpTroubleshootConfig {
	tags := getPolicyTagsFromSchema(d)

	obj := model.BgpTroubleshootConfig{
		Tags: tags,
	}

	if val, ok := d.GetOkExists("bfd_control_pkt_diagnostics"); ok {
		b := val.(bool)
		obj.BfdControlPktDiagnostics = &b
	}
	if val, ok := d.GetOkExists("bgp_session_diagnostics"); ok {
		b := val.(bool)
		obj.BgpSessionDiagnostics = &b
	}
	if val, ok := d.GetOkExists("system_diagnostics"); ok {
		b := val.(bool)
		obj.SystemDiagnostics = &b
	}

	return obj
}

func resourceNsxtPolicyRouteControllerBgpTroubleshootCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerPath := d.Get("route_controller_path").(string)
	routeControllerID := getResourceIDFromResourcePath(routeControllerPath, "route-controllers")
	if routeControllerID == "" {
		segs := strings.Split(routeControllerPath, "/")
		routeControllerID = segs[len(segs)-1]
	}

	obj := resourceNsxtPolicyRouteControllerBgpTroubleshootToStruct(d)

	log.Printf("[INFO] Creating Route Controller BGP Troubleshoot Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpTroubleshootClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Patch(routeControllerID, obj)
	if err != nil {
		return handleCreateError("Route Controller BGP Troubleshoot Config", routeControllerID, err)
	}

	d.SetId(routeControllerID)
	return resourceNsxtPolicyRouteControllerBgpTroubleshootRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpTroubleshootRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP Troubleshoot ID")
	}

	client := cliRouteControllerBgpTroubleshootClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID)
	if err != nil {
		return handleReadError(d, "Route Controller BGP Troubleshoot Config", routeControllerID, err)
	}

	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setPolicyTagsInSchema(d, obj.Tags)

	if obj.BfdControlPktDiagnostics != nil {
		d.Set("bfd_control_pkt_diagnostics", *obj.BfdControlPktDiagnostics)
	}
	if obj.BgpSessionDiagnostics != nil {
		d.Set("bgp_session_diagnostics", *obj.BgpSessionDiagnostics)
	}
	if obj.SystemDiagnostics != nil {
		d.Set("system_diagnostics", *obj.SystemDiagnostics)
	}

	return nil
}

func resourceNsxtPolicyRouteControllerBgpTroubleshootUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP Troubleshoot ID")
	}

	obj := resourceNsxtPolicyRouteControllerBgpTroubleshootToStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating Route Controller BGP Troubleshoot Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpTroubleshootClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	_, err := client.Update(routeControllerID, obj)
	if err != nil {
		return handleUpdateError("Route Controller BGP Troubleshoot Config", routeControllerID, err)
	}

	return resourceNsxtPolicyRouteControllerBgpTroubleshootRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpTroubleshootDelete(d *schema.ResourceData, m interface{}) error {
	routeControllerID := d.Id()
	if routeControllerID == "" {
		return fmt.Errorf("error obtaining Route Controller BGP Troubleshoot ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	log.Printf("[INFO] Deleting Route Controller BGP Troubleshoot Config for Route Controller %s", routeControllerID)

	client := cliRouteControllerBgpTroubleshootClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err := client.Delete(routeControllerID)
	if err != nil {
		return handleDeleteError("Route Controller BGP Troubleshoot Config", routeControllerID, err)
	}

	return nil
}
