// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcLBClientSslProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcLBClientSslProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcLBClientSslProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Client SSL Profile data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.LBClientSslProfile
	if objID != "" {
		// Get by id
		c := cliVpcLbClientSslProfilesClient(sessionContext, connector)
		objGet, err := c.Get(parents[0], parents[1], parents[2], objID)
		if err != nil {
			return handleDataSourceReadError(d, "VPC LBClientSslProfile", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining VPC LBClientSslProfile ID or name during read")
	} else {
		// Match names locally to preserve literal prefixes, including spaces and punctuation.
		resourceTypes := "LBClientSslProfile"
		resultValues, err := listPolicyResources(connector, sessionContext, resourceTypes, nil)
		if err != nil {
			return handleListError("VPC LBClientSslProfile", err)
		}
		var perfectMatch []model.LBClientSslProfile
		var prefixMatch []model.LBClientSslProfile
		converter := bindings.NewTypeConverter()
		// Prefer an exact name match over a prefix match across all search pages.
		for _, resultValue := range resultValues {
			converted, errs := converter.ConvertToGolang(resultValue, model.LBClientSslProfileBindingType())
			if len(errs) > 0 {
				return fmt.Errorf("Error while converting VPC LBClientSslProfile for name %q: %w", objName, errs[0])
			}
			objInList, ok := converted.(model.LBClientSslProfile)
			if !ok {
				return fmt.Errorf("Unexpected VPC LBClientSslProfile search result type %T", converted)
			}
			if objInList.DisplayName == nil {
				continue
			}
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBClientSslProfiles with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBClientSslProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VPC LBClientSslProfile with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
