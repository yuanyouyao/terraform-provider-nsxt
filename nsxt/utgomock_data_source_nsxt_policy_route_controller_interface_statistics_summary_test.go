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
	statsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/interfaces/statistics"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummary(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsSDK := statsmocks.NewMockSummaryClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerInterfaceStatisticsSummaryClientContext{
		Client:     mockStatsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerInterfaceStatisticsSummaryClient
	defer func() { cliRouteControllerInterfaceStatisticsSummaryClient = originalCli }()
	cliRouteControllerInterfaceStatisticsSummaryClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerInterfaceStatisticsSummaryClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	ifaceID := "iface-1"
	ifacePath := "/infra/route-controllers/rc-1/interfaces/iface-1"
	timestamp := int64(1234567890)

	totalPackets := int64(1000)
	totalBytes := int64(50000)
	droppedPackets := int64(5)
	blockedPackets := int64(0)

	t.Run("Get success", func(t *testing.T) {
		mockStatsSDK.EXPECT().Get(rcID, ifaceID).Return(model.RouteControllerInterfaceStatisticsSummary{
			InterfacePath:       &ifacePath,
			LastUpdateTimestamp: &timestamp,
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
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummary()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"interface_id":        ifaceID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummaryRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+ifaceID+"/statistics/summary", d.Id())
		assert.Equal(t, ifacePath, d.Get("interface_path"))
		assert.Equal(t, int(timestamp), d.Get("last_update_timestamp"))

		rxList := d.Get("rx").([]interface{})
		require.Len(t, rxList, 1)
		rx := rxList[0].(map[string]interface{})
		assert.Equal(t, int(totalPackets), rx["total_packets"])
		assert.Equal(t, int(totalBytes), rx["total_bytes"])
		assert.Equal(t, int(droppedPackets), rx["dropped_packets"])
		assert.Equal(t, int(blockedPackets), rx["blocked_packets"])
	})
}
