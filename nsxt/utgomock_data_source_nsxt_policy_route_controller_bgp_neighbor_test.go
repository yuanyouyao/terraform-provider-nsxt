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

func TestMockDataSourceNsxtPolicyRouteControllerBgpNeighbor(t *testing.T) {
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
	neighborPath := "/infra/route-controllers/rc-1/bgp/neighbors/neigh-1"
	revision := int64(1)
	displayName := "neigh-1-display"
	address := "192.168.1.1"
	remoteAs := "65001"

	t.Run("Read success", func(t *testing.T) {
		mockBgpNeighSDK.EXPECT().Get(rcID, neighborID).Return(model.RouteControllerBgpNeighborConfig{
			Id:              &neighborID,
			Path:            &neighborPath,
			Revision:        &revision,
			DisplayName:     &displayName,
			NeighborAddress: &address,
			RemoteAsNum:     &remoteAs,
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"id":                  neighborID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+neighborID, d.Id())
		assert.Equal(t, neighborPath, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, address, d.Get("neighbor_address"))
		assert.Equal(t, remoteAs, d.Get("remote_as_num"))
	})
}
