// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyRouteControllerBgpNeighbors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerBgpNeighborsRead,

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
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of BGP Neighbor. Supports regular expressions.",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of BGP Neighbor instance ID by display name.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerBgpNeighborsRead(d *schema.ResourceData, m interface{}) error {
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

	client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	result, err := client.List(routeControllerID, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error listing BGP neighbors for Route Controller %s: %v", routeControllerID, err)
	}

	resultMap := make(map[string]string)
	for _, item := range result.Results {
		id := ""
		if item.Id != nil {
			id = *item.Id
		}
		name := id
		if item.DisplayName != nil && *item.DisplayName != "" {
			name = *item.DisplayName
		}
		if id != "" {
			resultMap[id] = name
		}
	}

	d.SetId(newUUID())

	var re *regexp.Regexp
	if displayNameRegex, ok := d.GetOk("display_name"); ok {
		re, err = regexp.Compile(displayNameRegex.(string))
		if err != nil {
			return err
		}
		filteredMap := make(map[string]string)
		for id, name := range resultMap {
			if re.MatchString(name) {
				filteredMap[id] = name
			}
		}
		d.Set("items", filteredMap)
	} else {
		d.Set("items", resultMap)
	}

	return nil
}
