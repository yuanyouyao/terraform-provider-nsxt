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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	bgpneighmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockResourceNsxtPolicyRouteControllerBgpNeighbor(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpNeighSDK := bgpneighmocks.NewMockNeighborsClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpNeighborClientContext{
		Client:     mockBgpNeighSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpNeighborClient
	defer func() { cliRouteControllerBgpNeighborClient = originalCli }()
	cliRouteControllerBgpNeighborClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpNeighborClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	rcPath := "/infra/route-controllers/rc-1"
	neighborID := "neigh-1"
	neighborPath := "/infra/route-controllers/rc-1/bgp/neighbors/neigh-1"
	revision := int64(1)
	displayName := "neigh-1-display"
	description := "mock neighbor"
	address := "192.168.1.1"
	remoteAs := "65001"
	allowAsIn := true
	enabled := true
	holdDown := int64(180)
	keepAlive := int64(60)
	maxHop := int64(1)
	password := "secret"
	gwIps := []string{"192.168.1.254"}
	srcIps := []string{"192.168.1.2"}

	t.Run("Create success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().Get(rcID, neighborID).Return(model.RouteControllerBgpNeighborConfig{}, vapiErrors.NotFound{})
		mockBgpNeighSDK.EXPECT().Patch(rcID, neighborID, gomock.Any()).Return(nil)
		mockBgpNeighSDK.EXPECT().Get(rcID, neighborID).Return(model.RouteControllerBgpNeighborConfig{
			Id:                  &neighborID,
			Path:                &neighborPath,
			Revision:            &revision,
			DisplayName:         &displayName,
			Description:         &description,
			NeighborAddress:     &address,
			RemoteAsNum:         &remoteAs,
			AllowAsIn:           &allowAsIn,
			Enabled:             &enabled,
			HoldDownTime:        &holdDown,
			KeepAliveTime:       &keepAlive,
			MaximumHopLimit:     &maxHop,
			Password:            &password,
			GatewayIps:          gwIps,
			SourceAddresses:     srcIps,
			GracefulRestartMode: &displayName, // using displayName as placeholder for helper_only
			Bfd: &model.BgpBfdConfig{
				Enabled:  &enabled,
				Interval: &holdDown,
				Multiple: &maxHop,
			},
			RouteFiltering: []model.BgpRouteFiltering{
				{
					AddressFamily:   &displayName, // using displayName as placeholder for IPV4
					Enabled:         &enabled,
					InRouteFilters:  []string{"/infra/prefix-lists/in"},
					OutRouteFilters: []string{"/infra/prefix-lists/out"},
					MaximumRoutes:   &holdDown,
				},
			},
		}, nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"nsx_id":                neighborID,
			"display_name":          displayName,
			"description":           description,
			"neighbor_address":      address,
			"remote_as_num":         remoteAs,
			"allow_as_in":           allowAsIn,
			"enabled":               enabled,
			"hold_down_time":        int(holdDown),
			"keep_alive_time":       int(keepAlive),
			"maximum_hop_limit":     int(maxHop),
			"password":              password,
			"gateway_ips":           []interface{}{gwIps[0]},
			"source_addresses":      []interface{}{srcIps[0]},
			"bfd_config": []interface{}{
				map[string]interface{}{
					"enabled":  enabled,
					"interval": int(holdDown),
					"multiple": int(maxHop),
				},
			},
			"route_filtering": []interface{}{
				map[string]interface{}{
					"address_family":    displayName,
					"enabled":           enabled,
					"in_route_filters":  []interface{}{"/infra/prefix-lists/in"},
					"out_route_filters": []interface{}{"/infra/prefix-lists/out"},
					"maximum_routes":    int(holdDown),
				},
			},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+neighborID, d.Id())
		assert.Equal(t, neighborPath, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, address, d.Get("neighbor_address"))
		assert.Equal(t, remoteAs, d.Get("remote_as_num"))
		assert.Equal(t, allowAsIn, d.Get("allow_as_in"))
		assert.Equal(t, enabled, d.Get("enabled"))
	})

	t.Run("Update success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().Update(rcID, neighborID, gomock.Any()).Return(model.RouteControllerBgpNeighborConfig{
			Id:       &neighborID,
			Path:     &neighborPath,
			Revision: &revision,
		}, nil)
		mockBgpNeighSDK.EXPECT().Get(rcID, neighborID).Return(model.RouteControllerBgpNeighborConfig{
			Id:       &neighborID,
			Path:     &neighborPath,
			Revision: &revision,
		}, nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"neighbor_address":      address,
			"remote_as_num":         remoteAs,
			"revision":              0,
		})
		d.SetId(rcID + "/" + neighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().Delete(rcID, neighborID).Return(nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"neighbor_address":      address,
			"remote_as_num":         remoteAs,
		})
		d.SetId(rcID + "/" + neighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborDelete(d, m)
		require.NoError(t, err)
	})
}
