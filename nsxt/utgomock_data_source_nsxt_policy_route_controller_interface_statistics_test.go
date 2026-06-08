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
	statsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/interfaces"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerInterfaceStatistics(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsSDK := statsmocks.NewMockStatisticsClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerInterfaceStatisticsClientContext{
		Client:     mockStatsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerInterfaceStatisticsClient
	defer func() { cliRouteControllerInterfaceStatisticsClient = originalCli }()
	cliRouteControllerInterfaceStatisticsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerInterfaceStatisticsClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	ifaceID := "iface-1"
	ifacePath := "/infra/route-controllers/rc-1/interfaces/iface-1"
	tnID := "tn-1"
	subClusterID := "sc-1"
	lrpID := "lrp-1"
	vnaPath := "/infra/sites/site-1/enforcement-points/ep-1/virtual-network-appliances/vna-1"
	timestamp := int64(1234567890)

	totalPackets := int64(1000)
	totalBytes := int64(50000)
	droppedPackets := int64(5)
	blockedPackets := int64(0)

	t.Run("Get success", func(t *testing.T) {
		mockStatsSDK.EXPECT().Get(rcID, ifaceID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerInterfaceStatistics{
			InterfacePath: &ifacePath,
			PerNodeStatistics: []model.RouteControllerInterfaceStatisticsPerNode{
				{
					TransportNodeId:             &tnID,
					SubClusterId:                &subClusterID,
					LogicalRouterPortId:         &lrpID,
					VirtualNetworkAppliancePath: &vnaPath,
					LastUpdateTimestamp:         &timestamp,
					Rx: &model.LogicalRouterPortCounters{
						TotalPackets:   &totalPackets,
						TotalBytes:     &totalBytes,
						DroppedPackets: &droppedPackets,
						BlockedPackets: &blockedPackets,
					},
					Tx: &model.LogicalRouterPortCounters{
						TotalPackets:   &totalPackets,
						TotalBytes:     &totalBytes,
						DroppedPackets: &droppedPackets,
						BlockedPackets: &blockedPackets,
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerInterfaceStatistics()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"interface_id":        ifaceID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerInterfaceStatisticsRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+ifaceID+"/statistics", d.Id())
		assert.Equal(t, ifacePath, d.Get("interface_path"))

		perNode := d.Get("per_node_statistics").([]interface{})
		require.Len(t, perNode, 1)
		node := perNode[0].(map[string]interface{})
		assert.Equal(t, tnID, node["transport_node_id"])
		assert.Equal(t, subClusterID, node["sub_cluster_id"])
		assert.Equal(t, lrpID, node["logical_router_port_id"])
		assert.Equal(t, vnaPath, node["virtual_network_appliance_path"])
		assert.Equal(t, int(timestamp), node["last_update_timestamp"])

		rxList := node["rx"].([]interface{})
		require.Len(t, rxList, 1)
		rx := rxList[0].(map[string]interface{})
		assert.Equal(t, int(totalPackets), rx["total_packets"])
		assert.Equal(t, int(totalBytes), rx["total_bytes"])
		assert.Equal(t, int(droppedPackets), rx["dropped_packets"])
		assert.Equal(t, int(blockedPackets), rx["blocked_packets"])
	})
}
