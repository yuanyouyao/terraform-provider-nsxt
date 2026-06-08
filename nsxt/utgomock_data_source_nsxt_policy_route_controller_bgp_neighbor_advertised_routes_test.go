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
	bgpneighroutesmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp/neighbors"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutes(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpNeighRoutesSDK := bgpneighroutesmocks.NewMockAdvertisedRoutesClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpNeighborAdvertisedRoutesClientContext{
		Client:     mockBgpNeighRoutesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpNeighborAdvertisedRoutesClient
	defer func() { cliRouteControllerBgpNeighborAdvertisedRoutesClient = originalCli }()
	cliRouteControllerBgpNeighborAdvertisedRoutesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpNeighborAdvertisedRoutesClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	neighborID := "neigh-1"
	epPath := "/infra/enforcement-points/default"
	neighborPath := "/infra/route-controllers/rc-1/bgp/neighbors/neigh-1"
	sourceAddress := "192.168.1.1"
	transportNodeID := "tn-1"

	asPath := "65001 65002"
	esi := "esi-1"
	ethTag := int64(10)
	evpnRouteType := int64(2)
	localPref := int64(100)
	med := int64(20)
	network := "10.0.0.0/24"
	nextHop := "192.168.1.2"
	rd := "1:1"
	rmac := "00:50:56:af:12:34"
	rmacLen := int64(48)
	weight := int64(0)

	t.Run("List success", func(t *testing.T) {
		mockBgpNeighRoutesSDK.EXPECT().List(rcID, neighborID, nil, nil, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerBgpNeighborRoutesListResult{
			Results: []model.RouteControllerBgpNeighborRoutes{
				{
					EnforcementPointPath: &epPath,
					NeighborPath:         &neighborPath,
					VirtualNetworkApplianceRoutes: []model.RoutesPerTransportNode{
						{
							SourceAddress:   &sourceAddress,
							TransportNodeId: &transportNodeID,
							Routes: []model.RouteDetails{
								{
									AsPath:        &asPath,
									Esi:           &esi,
									EthTag:        &ethTag,
									EvpnRouteType: &evpnRouteType,
									LocalPref:     &localPref,
									Med:           &med,
									Network:       &network,
									NextHop:       &nextHop,
									Rd:            &rd,
									Rmac:          &rmac,
									RmacLen:       &rmacLen,
									Weight:        &weight,
								},
							},
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutes()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"neighbor_id":         neighborID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpNeighborAdvertisedRoutesRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+neighborID, d.Id())

		results := d.Get("results").([]interface{})
		assert.Len(t, results, 1)

		item := results[0].(map[string]interface{})
		assert.Equal(t, epPath, item["enforcement_point_path"])
		assert.Equal(t, neighborPath, item["neighbor_path"])

		vnaRoutes := item["virtual_network_appliance_routes"].([]interface{})
		assert.Len(t, vnaRoutes, 1)

		vnaItem := vnaRoutes[0].(map[string]interface{})
		assert.Equal(t, sourceAddress, vnaItem["source_address"])
		assert.Equal(t, transportNodeID, vnaItem["transport_node_id"])

		routes := vnaItem["routes"].([]interface{})
		assert.Len(t, routes, 1)

		routeItem := routes[0].(map[string]interface{})
		assert.Equal(t, asPath, routeItem["as_path"])
		assert.Equal(t, esi, routeItem["esi"])
		assert.Equal(t, int(ethTag), routeItem["eth_tag"])
		assert.Equal(t, int(evpnRouteType), routeItem["evpn_route_type"])
		assert.Equal(t, int(localPref), routeItem["local_pref"])
		assert.Equal(t, int(med), routeItem["med"])
		assert.Equal(t, network, routeItem["network"])
		assert.Equal(t, nextHop, routeItem["next_hop"])
		assert.Equal(t, rd, routeItem["rd"])
		assert.Equal(t, rmac, routeItem["rmac"])
		assert.Equal(t, int(rmacLen), routeItem["rmac_len"])
		assert.Equal(t, int(weight), routeItem["weight"])
	})
}
