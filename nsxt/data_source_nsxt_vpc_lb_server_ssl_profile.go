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

func dataSourceNsxtVpcLBServerSslProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcLBServerSslProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcLBServerSslProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Server SSL Profile data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.LBServerSslProfile
	if objID != "" {
		// Get by id
		c := cliVpcLbServerSslProfilesClient(sessionContext, connector)
		objGet, err := c.Get(parents[0], parents[1], parents[2], objID)
		if err != nil {
			return handleDataSourceReadError(d, "VPC LBServerSslProfile", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining VPC LBServerSslProfile ID or name during read")
	} else {
		// Match names locally to preserve literal prefixes, including spaces and punctuation.
		resourceTypes := "LBServerSslProfile"
		resultValues, err := listPolicyResources(connector, sessionContext, resourceTypes, nil)
		if err != nil {
			return handleListError("VPC LBServerSslProfile", err)
		}
		var perfectMatch []model.LBServerSslProfile
		var prefixMatch []model.LBServerSslProfile
		converter := bindings.NewTypeConverter()
		// Prefer an exact name match over a prefix match across all search pages.
		for _, resultValue := range resultValues {
			converted, errs := converter.ConvertToGolang(resultValue, model.LBServerSslProfileBindingType())
			if len(errs) > 0 {
				return fmt.Errorf("Error while converting VPC LBServerSslProfile for name %q: %w", objName, errs[0])
			}
			objInList, ok := converted.(model.LBServerSslProfile)
			if !ok {
				return fmt.Errorf("Unexpected VPC LBServerSslProfile search result type %T", converted)
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
				return fmt.Errorf("Found multiple VPC LBServerSslProfiles with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBServerSslProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VPC LBServerSslProfile with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
