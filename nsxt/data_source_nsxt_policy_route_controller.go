// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyRouteController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"path":         getPathSchema(),
			"display_name": getOptionalDisplayNameSchema(true),
			"description":  getDescriptionSchema(),
			"ha_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "High-availability mode for route controller.",
			},
			"virtual_network_appliance_cluster_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy path for virtual network appliance cluster.",
			},
		},
	}
}

func dataSourceNsxtPolicyRouteControllerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	obj, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "RouteController", nil, false)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.RouteControllerBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	routeController := dataValue.(model.RouteController)

	d.SetId(*routeController.Id)
	d.Set("display_name", routeController.DisplayName)
	d.Set("description", routeController.Description)
	d.Set("path", routeController.Path)

	if routeController.HaMode != nil {
		d.Set("ha_mode", routeController.HaMode)
	}
	if routeController.VirtualNetworkApplianceClusterPath != nil {
		d.Set("virtual_network_appliance_cluster_path", routeController.VirtualNetworkApplianceClusterPath)
	}

	return nil
}
