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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	rcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	routeControllerID          = "rc-1"
	routeControllerDisplayName = "rc-fooname"
	routeControllerDescription = "route controller mock"
	routeControllerPath        = "/infra/route-controllers/rc-1"
	routeControllerRevision    = int64(1)
	routeControllerHaMode      = "ACTIVE_STANDBY"
	routeControllerVnaPath     = "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-1"
)

func TestMockResourceNsxtPolicyRouteControllerCreate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK := rcmocks.NewMockRouteControllersClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerClientContext{
		Client:     mockRcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllersClient
	defer func() { cliRouteControllersClient = originalCli }()
	cliRouteControllersClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockRcSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockRcSDK.EXPECT().Get(gomock.Any()).Return(model.RouteController{
			Id:                                 &routeControllerID,
			DisplayName:                        &routeControllerDisplayName,
			Description:                        &routeControllerDescription,
			Path:                               &routeControllerPath,
			Revision:                           &routeControllerRevision,
			HaMode:                             &routeControllerHaMode,
			VirtualNetworkApplianceClusterPath: &routeControllerVnaPath,
		}, nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockRcSDK.EXPECT().Get("existing-id").Return(model.RouteController{Id: &routeControllerID}, nil)

		res := resourceNsxtPolicyRouteController()
		data := minimalRouteControllerData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyRouteControllerRead(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK := rcmocks.NewMockRouteControllersClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerClientContext{
		Client:     mockRcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllersClient
	defer func() { cliRouteControllersClient = originalCli }()
	cliRouteControllersClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockRcSDK.EXPECT().Get(routeControllerID).Return(model.RouteController{
			Id:                                 &routeControllerID,
			DisplayName:                        &routeControllerDisplayName,
			Description:                        &routeControllerDescription,
			Path:                               &routeControllerPath,
			Revision:                           &routeControllerRevision,
			HaMode:                             &routeControllerHaMode,
			VirtualNetworkApplianceClusterPath: &routeControllerVnaPath,
		}, nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, routeControllerDisplayName, d.Get("display_name"))
		assert.Equal(t, routeControllerDescription, d.Get("description"))
		assert.Equal(t, routeControllerPath, d.Get("path"))
		assert.Equal(t, int(routeControllerRevision), d.Get("revision"))
		assert.Equal(t, routeControllerHaMode, d.Get("ha_mode"))
		assert.Equal(t, routeControllerVnaPath, d.Get("virtual_network_appliance_cluster_path"))
	})
}

func TestMockResourceNsxtPolicyRouteControllerUpdate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK := rcmocks.NewMockRouteControllersClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerClientContext{
		Client:     mockRcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllersClient
	defer func() { cliRouteControllersClient = originalCli }()
	cliRouteControllersClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockRcSDK.EXPECT().Update(gomock.Any(), gomock.Any()).Return(model.RouteController{
			Id:                                 &routeControllerID,
			DisplayName:                        &routeControllerDisplayName,
			Description:                        &routeControllerDescription,
			Path:                               &routeControllerPath,
			Revision:                           &routeControllerRevision,
			HaMode:                             &routeControllerHaMode,
			VirtualNetworkApplianceClusterPath: &routeControllerVnaPath,
		}, nil)
		mockRcSDK.EXPECT().Get(gomock.Any()).Return(model.RouteController{
			Id:                                 &routeControllerID,
			DisplayName:                        &routeControllerDisplayName,
			Description:                        &routeControllerDescription,
			Path:                               &routeControllerPath,
			Revision:                           &routeControllerRevision,
			HaMode:                             &routeControllerHaMode,
			VirtualNetworkApplianceClusterPath: &routeControllerVnaPath,
		}, nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerUpdate(d, m)
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyRouteControllerDelete(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK := rcmocks.NewMockRouteControllersClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerClientContext{
		Client:     mockRcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllersClient
	defer func() { cliRouteControllersClient = originalCli }()
	cliRouteControllersClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockRcSDK.EXPECT().Delete(routeControllerID).Return(nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerDelete(d, m)
		require.NoError(t, err)
	})
}

func minimalRouteControllerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":                           routeControllerDisplayName,
		"description":                            routeControllerDescription,
		"ha_mode":                                routeControllerHaMode,
		"virtual_network_appliance_cluster_path": routeControllerVnaPath,
	}
}
