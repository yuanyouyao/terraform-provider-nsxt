//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	rcstatemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerStateRead(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStateSDK := rcstatemocks.NewMockStateClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerStateClientContext{
		Client:     mockStateSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerStateClient
	defer func() { cliRouteControllerStateClient = originalCli }()
	cliRouteControllerStateClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerStateClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	logicalGWID := "lgw-1"
	vnaClusterPath := "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-1"
	lastUpdate := int64(1234567890)
	haStatus := "ACTIVE"
	nodeType := "APPLIANCE"
	sgwID := "sgw-1"
	vnaPath := "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-1/virtual-network-appliances/vna-node-1"

	t.Run("success", func(t *testing.T) {
		mockStateSDK.EXPECT().Get(rcID, nil).Return(model.RouteControllerState{
			LastUpdateTimestamp:                &lastUpdate,
			LogicalGatewayId:                   &logicalGWID,
			VirtualNetworkApplianceClusterPath: &vnaClusterPath,
			PerNodeStatus: []model.RouteControllerStateNodeStatus{
				{
					HighAvailabilityStatus:      &haStatus,
					NodeType:                    &nodeType,
					ServiceGatewayId:            &sgwID,
					VirtualNetworkAppliancePath: &vnaPath,
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerState()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		err := dataSourceNsxtPolicyRouteControllerStateRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, int(lastUpdate), d.Get("last_update_timestamp"))
		assert.Equal(t, logicalGWID, d.Get("logical_gateway_id"))
		assert.Equal(t, vnaClusterPath, d.Get("virtual_network_appliance_cluster_path"))

		nodeStatusList := d.Get("per_node_status").([]interface{})
		require.Len(t, nodeStatusList, 1)
		nodeStatus := nodeStatusList[0].(map[string]interface{})
		assert.Equal(t, haStatus, nodeStatus["high_availability_status"])
		assert.Equal(t, nodeType, nodeStatus["node_type"])
		assert.Equal(t, sgwID, nodeStatus["service_gateway_id"])
		assert.Equal(t, vnaPath, nodeStatus["virtual_network_appliance_path"])
	})

	t.Run("success with path input", func(t *testing.T) {
		mockStateSDK.EXPECT().Get(rcID, nil).Return(model.RouteControllerState{
			LogicalGatewayId: &logicalGWID,
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerState()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_path": "/infra/route-controllers/rc-1",
		})

		err := dataSourceNsxtPolicyRouteControllerStateRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, logicalGWID, d.Get("logical_gateway_id"))
	})

	t.Run("get error", func(t *testing.T) {
		mockStateSDK.EXPECT().Get(rcID, nil).Return(model.RouteControllerState{}, errors.New("state fail"))

		ds := dataSourceNsxtPolicyRouteControllerState()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		err := dataSourceNsxtPolicyRouteControllerStateRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error getting route controller state")
	})
}
