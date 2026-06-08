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
	rtmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgpRouteTable(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpRouteTableSDK := rtmocks.NewMockBgpRouteTableClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpRouteTableClientContext{
		Client:     mockBgpRouteTableSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpRouteTableClient
	defer func() { cliRouteControllerBgpRouteTableClient = originalCli }()
	cliRouteControllerBgpRouteTableClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpRouteTableClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	vnaPath := "/infra/virtual-network-appliances/vna-1"
	prefix := "10.0.0.0/24"
	epPath := "/infra/enforcement-points/ep-1"
	gwPath := "/infra/tier-0s/t0-1"
	timestamp := int64(1234567890)
	tnID := "tn-1"
	tnPath := "/infra/transport-nodes/tn-1"

	asPath := "65001 65002"
	best := true
	comm := "65001:100"
	esi := "esi-1"
	ethTag := int64(10)
	evpnType := int64(2)
	extComm := "ext-1"
	largeComm := "large-1"
	localPref := int64(100)
	med := int64(50)
	multi := false
	network := "10.0.0.0/24"
	pathFrom := "path-1"
	peerID := "peer-1"
	rd := "rd-1"
	rmac := "rmac-1"
	origin := "IGP"
	stale := false
	valid := true
	vni := int64(1000)
	weight := int64(32768)

	afi := "afi-1"
	ip := "192.168.1.1"
	scope := "scope-1"
	used := true

	t.Run("List success", func(t *testing.T) {
		mockBgpRouteTableSDK.EXPECT().List(rcID, vnaPath, &prefix).Return(model.BgpRIBListResult{
			Results: []model.BgpRoutes{
				{
					EnforcementPointPath: &epPath,
					GatewayPath:          &gwPath,
					LastUpdateTimestamp:  &timestamp,
					TransportNodeId:      &tnID,
					TransportNodePath:    &tnPath,
					RouteDetails: []model.BgpRouteDetails{
						{
							AsPath:            &asPath,
							Bestpath:          &best,
							Community:         &comm,
							Esi:               &esi,
							EthTag:            &ethTag,
							EvpnRouteType:     &evpnType,
							ExtendedCommunity: &extComm,
							LargeCommunity:    &largeComm,
							LocalPref:         &localPref,
							Med:               &med,
							Multipath:         &multi,
							Network:           &network,
							PathFrom:          &pathFrom,
							PeerId:            &peerID,
							Rd:                &rd,
							Rmac:              &rmac,
							RouteOrigin:       &origin,
							Stale:             &stale,
							Valid:             &valid,
							Vni:               &vni,
							Weight:            &weight,
							Nexthops: []model.BgpNexthop{
								{
									Afi:   &afi,
									Ip:    &ip,
									Scope: &scope,
									Used:  &used,
								},
							},
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpRouteTable()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id":            rcID,
			"virtual_network_appliance_path": vnaPath,
			"network_prefix":                 prefix,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpRouteTableRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/bgp/route-table", d.Id())

		resList := d.Get("results").([]interface{})
		require.Len(t, resList, 1)
		item := resList[0].(map[string]interface{})
		assert.Equal(t, epPath, item["enforcement_point_path"])
		assert.Equal(t, gwPath, item["gateway_path"])
		assert.Equal(t, int(timestamp), item["last_update_timestamp"])
		assert.Equal(t, tnID, item["transport_node_id"])
		assert.Equal(t, tnPath, item["transport_node_path"])

		rdList := item["route_details"].([]interface{})
		require.Len(t, rdList, 1)
		rdItem := rdList[0].(map[string]interface{})
		assert.Equal(t, asPath, rdItem["as_path"])
		assert.Equal(t, best, rdItem["bestpath"])
		assert.Equal(t, comm, rdItem["community"])
		assert.Equal(t, esi, rdItem["esi"])
		assert.Equal(t, int(ethTag), rdItem["eth_tag"])
		assert.Equal(t, int(evpnType), rdItem["evpn_route_type"])
		assert.Equal(t, extComm, rdItem["extended_community"])
		assert.Equal(t, largeComm, rdItem["large_community"])
		assert.Equal(t, int(localPref), rdItem["local_pref"])
		assert.Equal(t, int(med), rdItem["med"])
		assert.Equal(t, multi, rdItem["multipath"])
		assert.Equal(t, network, rdItem["network"])
		assert.Equal(t, pathFrom, rdItem["path_from"])
		assert.Equal(t, peerID, rdItem["peer_id"])
		assert.Equal(t, rd, rdItem["rd"])
		assert.Equal(t, rmac, rdItem["rmac"])
		assert.Equal(t, origin, rdItem["route_origin"])
		assert.Equal(t, stale, rdItem["stale"])
		assert.Equal(t, valid, rdItem["valid"])
		assert.Equal(t, int(vni), rdItem["vni"])
		assert.Equal(t, int(weight), rdItem["weight"])

		nhList := rdItem["nexthops"].([]interface{})
		require.Len(t, nhList, 1)
		nhItem := nhList[0].(map[string]interface{})
		assert.Equal(t, afi, nhItem["afi"])
		assert.Equal(t, ip, nhItem["ip"])
		assert.Equal(t, scope, nhItem["scope"])
		assert.Equal(t, used, nhItem["used"])
	})
}
