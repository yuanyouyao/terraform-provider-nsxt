// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyRouteControllerBgpTroubleshoot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpTroubleshootRead,

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
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the route controller BGP troubleshoot config.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The revision of the resource.",
			},
			"tag": getTagsSchema(),
			"bfd_control_pkt_diagnostics": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable/disable the collection of the timestamps of sent and received BFD control messages per BFD peer session.",
			},
			"bgp_session_diagnostics": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable/disable the collection of the timestamps of sent and received Keep-Alive messages per BGP peer session, and the session states.",
			},
			"system_diagnostics": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable/disable the collection of system diagnostic data such as ARP, Ping, CPU stats, etc.",
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpTroubleshootRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerBgpTroubleshootClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID)
	if err != nil {
		return fmt.Errorf("error reading Route Controller BGP Troubleshoot Config for Route Controller %s: %v", routeControllerID, err)
	}

	d.SetId(routeControllerID)
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
