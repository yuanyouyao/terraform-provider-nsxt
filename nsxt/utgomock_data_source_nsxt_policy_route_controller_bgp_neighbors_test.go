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
	bgpneighmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgpNeighbors(t *testing.T) {
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
	neighborID := "neigh-1"
	displayName := "neigh-1-display"

	t.Run("List success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().List(rcID, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerBgpNeighborConfigListResult{
			Results: []model.RouteControllerBgpNeighborConfig{
				{
					Id:          &neighborID,
					DisplayName: &displayName,
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbors()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpNeighborsRead(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())

		items := d.Get("items").(map[string]interface{})
		assert.Len(t, items, 1)
		assert.Equal(t, displayName, items[neighborID])
	})

	t.Run("List with regex filter success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().List(rcID, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerBgpNeighborConfigListResult{
			Results: []model.RouteControllerBgpNeighborConfig{
				{
					Id:          &neighborID,
					DisplayName: &displayName,
				},
				{
					Id:          &rcID, // different ID
					DisplayName: &rcID, // different name
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbors()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"display_name":        "neigh-.*",
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpNeighborsRead(d, m)
		require.NoError(t, err)

		items := d.Get("items").(map[string]interface{})
		assert.Len(t, items, 1)
		assert.Equal(t, displayName, items[neighborID])
	})
}
