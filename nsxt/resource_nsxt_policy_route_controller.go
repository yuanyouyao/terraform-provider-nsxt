// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var cliRouteControllersClient = infra.NewRouteControllersClient

var routeControllerSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"ha_mode": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      model.RouteController_HA_MODE_STANDBY,
			Description:  "High-availability mode for route controller.",
			ValidateFunc: validation.StringInSlice([]string{model.RouteController_HA_MODE_STANDBY}, false),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "HaMode",
		},
	},
	"virtual_network_appliance_cluster_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Policy path for virtual network appliance cluster.",
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "VirtualNetworkApplianceClusterPath",
		},
	},
}

const routeControllerPathExample = "/infra/route-controllers/[route-controller-id]"

func resourceNsxtPolicyRouteController() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerCreate,
		Read:   resourceNsxtPolicyRouteControllerRead,
		Update: resourceNsxtPolicyRouteControllerUpdate,
		Delete: resourceNsxtPolicyRouteControllerDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(routeControllerPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(routeControllerSchema),
	}
}

func resourceNsxtPolicyRouteControllerExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error

	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliRouteControllersClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyRouteControllerCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyRouteControllerExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.RouteController{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, routeControllerSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating RouteController with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliRouteControllersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("RouteController", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyRouteControllerRead(d, m)
}

func resourceNsxtPolicyRouteControllerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliRouteControllersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "RouteController", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.StructToSchema(elem, d, routeControllerSchema, "", nil); err != nil {
		return err
	}

	return nil
}

func resourceNsxtPolicyRouteControllerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.RouteController{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, routeControllerSchema, "", nil); err != nil {
		return err
	}

	sessionContext := getSessionContext(d, m)
	client := cliRouteControllersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("RouteController", id, err)
	}

	return resourceNsxtPolicyRouteControllerRead(d, m)
}

func resourceNsxtPolicyRouteControllerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	connector := getPolicyConnector(m)

	sessionContext := getSessionContext(d, m)
	client := cliRouteControllersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("RouteController", id, err)
	}

	return nil
}
