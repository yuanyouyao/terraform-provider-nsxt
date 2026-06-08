// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyRouteControllers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllersRead,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of Route Controller. Supports regular expressions",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of Route Controller instance ID by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"context": getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
		},
	}
}

func dataSourceNsxtPolicyRouteControllersRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	resultMap := make(map[string]string)
	err := policyDataSourceCreateMap(connector, getSessionContext(d, m), "RouteController", resultMap, nil)
	if err != nil {
		return fmt.Errorf("error in listing the Route Controller items : %v", err)
	}
	d.SetId(newUUID())

	// read the display_name, may or may not be regex expression
	var re *regexp.Regexp
	if displayNameRegex, ok := d.GetOk("display_name"); ok {
		re, err = regexp.Compile(displayNameRegex.(string))
		if err != nil {
			return err
		}
		// Filter the resultMap by matching displayname with the regex
		filteredMap := make(map[string]string)
		for id, displayName := range resultMap {
			if re.MatchString(displayName) {
				filteredMap[id] = displayName
			}
		}
		d.Set("items", filteredMap)
	} else {
		// If no display_name is provided, set the resultMap as is
		d.Set("items", resultMap)
	}
	return nil
}
