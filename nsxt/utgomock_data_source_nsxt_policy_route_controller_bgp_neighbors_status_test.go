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
	bgpneighstatusmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp/neighbors"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgpNeighborsStatus(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpNeighStatusSDK := bgpneighstatusmocks.NewMockStatusClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpNeighborsStatusClientContext{
		Client:     mockBgpNeighStatusSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpNeighborsStatusClient
	defer func() { cliRouteControllerBgpNeighborsStatusClient = originalCli }()
	cliRouteControllerBgpNeighborsStatusClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpNeighborsStatusClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	neighborAddress := "192.168.1.1"
	connectionState := "ESTABLISHED"
	connectionDropCount := int64(2)
	establishedConnectionCount := int64(3)
	gracefulRestartMode := "HELPER_ONLY"
	holdTime := int64(180)
	isDynamic := false
	keepAliveInterval := int64(60)
	localPort := int64(179)
	messagesReceived := int64(100)
	messagesSent := int64(101)
	neighborEdgeNode := "edge-1"
	neighborPath := "/infra/route-controllers/rc-1/bgp/neighbors/neigh-1"
	neighborRouterID := "10.0.0.1"
	remoteAsNumber := "65001"
	remotePort := int64(179)
	remoteSitePath := "/infra/sites/site-1"
	routeControllerPath := "/infra/route-controllers/rc-1"
	sourceAddress := "192.168.1.2"
	timeSinceEstablished := int64(3600)
	totalInPrefixCount := int64(10)
	totalOutPrefixCount := int64(20)
	bgpType := "USER"
	vnaPath := "/infra/sites/site-1/enforcement-points/ep-1/virtual-network-appliances/vna-1"

	announcedCapabilities := []string{"cap-1", "cap-2"}
	negotiatedCapabilities := []string{"cap-1"}

	afType := "IPV4_UNICAST"
	afInPrefixCount := int64(10)
	afOutPrefixCount := int64(20)

	t.Run("List success", func(t *testing.T) {
		mockBgpNeighStatusSDK.EXPECT().List(rcID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerBgpNeighborsStatusListResult{
			Results: []model.RouteControllerBgpNeighborStatus{
				{
					NeighborAddress:             &neighborAddress,
					ConnectionState:             &connectionState,
					ConnectionDropCount:         &connectionDropCount,
					EstablishedConnectionCount:  &establishedConnectionCount,
					GracefulRestartMode:         &gracefulRestartMode,
					HoldTime:                    &holdTime,
					IsDynamic:                   &isDynamic,
					KeepAliveInterval:           &keepAliveInterval,
					LocalPort:                   &localPort,
					MessagesReceived:            &messagesReceived,
					MessagesSent:                &messagesSent,
					NeighborEdgeNode:            &neighborEdgeNode,
					NeighborPath:                &neighborPath,
					NeighborRouterId:            &neighborRouterID,
					RemoteAsNumber:              &remoteAsNumber,
					RemotePort:                  &remotePort,
					RemoteSitePath:              &remoteSitePath,
					RouteControllerPath:         &routeControllerPath,
					SourceAddress:               &sourceAddress,
					TimeSinceEstablished:        &timeSinceEstablished,
					TotalInPrefixCount:          &totalInPrefixCount,
					TotalOutPrefixCount:         &totalOutPrefixCount,
					Type_:                       &bgpType,
					VirtualNetworkAppliancePath: &vnaPath,
					AnnouncedCapabilities:       announcedCapabilities,
					NegotiatedCapability:        negotiatedCapabilities,
					AddressFamilies: []model.BgpAddressFamily{
						{
							Type_:          &afType,
							InPrefixCount:  &afInPrefixCount,
							OutPrefixCount: &afOutPrefixCount,
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighborsStatus()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpNeighborsStatusRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())

		results := d.Get("results").([]interface{})
		assert.Len(t, results, 1)

		item := results[0].(map[string]interface{})
		assert.Equal(t, neighborAddress, item["neighbor_address"])
		assert.Equal(t, connectionState, item["connection_state"])
		assert.Equal(t, int(connectionDropCount), item["connection_drop_count"])
		assert.Equal(t, int(establishedConnectionCount), item["established_connection_count"])
		assert.Equal(t, gracefulRestartMode, item["graceful_restart_mode"])
		assert.Equal(t, int(holdTime), item["hold_time"])
		assert.Equal(t, isDynamic, item["is_dynamic"])
		assert.Equal(t, int(keepAliveInterval), item["keep_alive_interval"])
		assert.Equal(t, int(localPort), item["local_port"])
		assert.Equal(t, int(messagesReceived), item["messages_received"])
		assert.Equal(t, int(messagesSent), item["messages_sent"])
		assert.Equal(t, neighborEdgeNode, item["neighbor_edge_node"])
		assert.Equal(t, neighborPath, item["neighbor_path"])
		assert.Equal(t, neighborRouterID, item["neighbor_router_id"])
		assert.Equal(t, remoteAsNumber, item["remote_as_number"])
		assert.Equal(t, int(remotePort), item["remote_port"])
		assert.Equal(t, remoteSitePath, item["remote_site_path"])
		assert.Equal(t, routeControllerPath, item["route_controller_path"])
		assert.Equal(t, sourceAddress, item["source_address"])
		assert.Equal(t, int(timeSinceEstablished), item["time_since_established"])
		assert.Equal(t, int(totalInPrefixCount), item["total_in_prefix_count"])
		assert.Equal(t, int(totalOutPrefixCount), item["total_out_prefix_count"])
		assert.Equal(t, bgpType, item["type"])
		assert.Equal(t, vnaPath, item["virtual_network_appliance_path"])

		announcedCaps := item["announced_capabilities"].([]interface{})
		assert.Len(t, announcedCaps, 2)
		assert.Equal(t, "cap-1", announcedCaps[0])

		negotiatedCaps := item["negotiated_capabilities"].([]interface{})
		assert.Len(t, negotiatedCaps, 1)
		assert.Equal(t, "cap-1", negotiatedCaps[0])

		afs := item["address_families"].([]interface{})
		assert.Len(t, afs, 1)
		afItem := afs[0].(map[string]interface{})
		assert.Equal(t, afType, afItem["type"])
		assert.Equal(t, int(afInPrefixCount), afItem["in_prefix_count"])
		assert.Equal(t, int(afOutPrefixCount), afItem["out_prefix_count"])
	})
}
