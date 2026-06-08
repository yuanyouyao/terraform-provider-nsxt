//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func routeControllerToStructValue(t *testing.T, rc model.RouteController) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(rc, model.RouteControllerBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockDataSourceNsxtPolicyRouteControllerRead(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	rcID := "rc-1"
	rcName := "rc-test"
	rcPath := "/infra/route-controllers/rc-1"
	rcDesc := "test desc"
	haMode := "ACTIVE_STANDBY"
	vnaClusterPath := "/infra/virtual-network-appliance-clusters/vna-1"
	rt := "RouteController"

	rc := model.RouteController{
		Id:                                 &rcID,
		DisplayName:                        &rcName,
		Path:                               &rcPath,
		Description:                        &rcDesc,
		HaMode:                             &haMode,
		VirtualNetworkApplianceClusterPath: &vnaClusterPath,
		ResourceType:                       &rt,
	}

	sv := routeControllerToStructValue(t, rc)

	t.Run("success by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results:     []*data.StructValue{sv},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": rcID,
		})

		err := dataSourceNsxtPolicyRouteControllerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, rcName, d.Get("display_name"))
		assert.Equal(t, rcDesc, d.Get("description"))
		assert.Equal(t, rcPath, d.Get("path"))
		assert.Equal(t, haMode, d.Get("ha_mode"))
		assert.Equal(t, vnaClusterPath, d.Get("virtual_network_appliance_cluster_path"))
	})

	t.Run("success by display_name", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results:     []*data.StructValue{sv},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": rcName,
		})

		err := dataSourceNsxtPolicyRouteControllerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
	})
}
