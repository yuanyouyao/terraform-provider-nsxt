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

var cliRouteControllerInterfaceStatisticsSummaryClient = infra.NewRouteControllerInterfaceStatisticsSummaryClient

func dataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummaryRead,

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
			"interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the route controller interface.",
			},
			"interface_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy path of the interface.",
			},
			"last_update_timestamp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of the last update.",
			},
			"rx": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "RX counters.",
				Elem: &schema.Resource{
					Schema: getLogicalRouterPortCountersSchema(),
				},
			},
			"tx": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TX counters.",
				Elem: &schema.Resource{
					Schema: getLogicalRouterPortCountersSchema(),
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummaryRead(d *schema.ResourceData, m interface{}) error {
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

	interfaceID := d.Get("interface_id").(string)

	client := cliRouteControllerInterfaceStatisticsSummaryClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(routeControllerID, interfaceID)
	if err != nil {
		return fmt.Errorf("error getting route controller interface statistics summary for Route Controller %s and Interface %s: %v", routeControllerID, interfaceID, err)
	}

	d.SetId(routeControllerID + "/" + interfaceID + "/statistics/summary")

	if obj.InterfacePath != nil {
		d.Set("interface_path", *obj.InterfacePath)
	}
	if obj.LastUpdateTimestamp != nil {
		d.Set("last_update_timestamp", int(*obj.LastUpdateTimestamp))
	}
	if obj.Rx != nil {
		d.Set("rx", flattenLogicalRouterPortCounters(obj.Rx))
	}
	if obj.Tx != nil {
		d.Set("tx", flattenLogicalRouterPortCounters(obj.Tx))
	}

	return nil
}
