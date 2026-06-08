// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyRouteControllerBgp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpRead,

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
				Description: "The policy path of the route controller BGP routing config.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The revision of the resource.",
			},
			"tag": getTagsSchema(),
			"ecmp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable ECMP.",
			},
			"local_as_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BGP AS number in ASPLAIN/ASDOT format.",
			},
			"multipath_relax": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to enable BGP multipath relax option.",
			},
			"peer_route_convergence_timer": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Extra time in seconds the router must wait before sending the UP notification.",
			},
			"graceful_restart_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BGP Graceful Restart Configuration Mode.",
			},
			"graceful_restart_timer": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum time taken (in seconds) for a BGP session to be established after a restart.",
			},
			"graceful_restart_stale_route_timer": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum time (in seconds) before stale routes are removed from the RIB when BGP restarts.",
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerBgpClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID)
	if err != nil {
		return fmt.Errorf("error getting route controller BGP config for ID %s: %v", routeControllerID, err)
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

	d.SetId(routeControllerID)
	return nil
}
