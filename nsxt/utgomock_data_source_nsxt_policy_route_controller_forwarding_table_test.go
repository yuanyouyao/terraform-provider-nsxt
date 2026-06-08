//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ftmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerForwardingTableRead(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockForwardingTableSDK := ftmocks.NewMockForwardingTableClient(ctrl)

	rcID := "rc-1"
	networkPrefix := "10.0.0.0/24"
	routeSource := "BGP"
	vnaPath := "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-1/virtual-network-appliances/vna-node-1"

	csvContent := "network,next_hop,route_type\n10.0.0.0/24,192.168.1.1,BGP"

	httpClient := &http.Client{
		Transport: mockRoundTripper(func(req *http.Request) (*http.Response, error) {
			assert.Contains(t, req.URL.Path, "/infra/route-controllers/rc-1/forwarding-table/download")
			assert.Equal(t, "network_prefix=10.0.0.0%2F24&route_source=BGP&virtual_network_appliance_path=%2Finfra%2Fsites%2Fdefault%2Fenforcement-points%2Fdefault%2Fvirtual-network-appliance-clusters%2Fvna-1%2Fvirtual-network-appliances%2Fvna-node-1", req.URL.RawQuery)
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(csvContent)),
			}, nil
		}),
	}

	mockWrapper := &cliinfra.RouteControllerForwardingTableClientContext{
		Client:     mockForwardingTableSDK,
		ClientType: utl.Local,
		Host:       "localhost",
		HTTPClient: httpClient,
	}

	originalCli := cliRouteControllerForwardingTableClient
	defer func() { cliRouteControllerForwardingTableClient = originalCli }()
	cliRouteControllerForwardingTableClient = func(sessionContext utl.SessionContext, connector client.Connector, host string, hc *http.Client, username, password, token, cookie, xsrf string) *cliinfra.RouteControllerForwardingTableClientContext {
		return mockWrapper
	}

	t.Run("success", func(t *testing.T) {
		count := int64(1)
		status := "SUCCESS"
		tnPath := "/infra/sites/default/enforcement-points/default/transport-nodes/tn-1"
		adminDist := int64(20)
		blackHole := false
		lrCompID := "lr-1"
		lrCompType := "TIER0"
		network := "10.0.0.0/24"
		nextHop := "192.168.1.1"
		nextHopGW := "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-1"
		routeType := "BGP"

		mockForwardingTableSDK.EXPECT().Get(rcID, &networkPrefix, &routeSource, &vnaPath).Return(model.RouteControllerRoutingTableListResult{
			Results: []model.RouteControllerRoutingTable{
				{
					Count:                       &count,
					Status:                      &status,
					TransportNodePath:           &tnPath,
					VirtualNetworkAppliancePath: &vnaPath,
					RouteEntries: []model.RoutingEntry{
						{
							AdminDistance:   &adminDist,
							BlackHole:       &blackHole,
							LrComponentId:   &lrCompID,
							LrComponentType: &lrCompType,
							Network:         &network,
							NextHop:         &nextHop,
							NextHopGateway:  &nextHopGW,
							RouteType:       &routeType,
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerForwardingTable()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id":            rcID,
			"network_prefix":                 networkPrefix,
			"route_source":                   routeSource,
			"virtual_network_appliance_path": vnaPath,
		})

		mockClients := newGoMockProviderClient()
		mockClients.PolicyHTTPClient = httpClient
		mockClients.Host = "localhost"

		err := dataSourceNsxtPolicyRouteControllerForwardingTableRead(d, mockClients)
		require.NoError(t, err)

		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, csvContent, d.Get("forwarding_table_csv"))

		forwardingTableList := d.Get("forwarding_table").([]interface{})
		require.Len(t, forwardingTableList, 1)
		rt := forwardingTableList[0].(map[string]interface{})
		assert.Equal(t, int(count), rt["count"])
		assert.Equal(t, status, rt["status"])
		assert.Equal(t, tnPath, rt["transport_node_path"])
		assert.Equal(t, vnaPath, rt["virtual_network_appliance_path"])

		entriesList := rt["route_entries"].([]interface{})
		require.Len(t, entriesList, 1)
		entry := entriesList[0].(map[string]interface{})
		assert.Equal(t, int(adminDist), entry["admin_distance"])
		assert.Equal(t, blackHole, entry["black_hole"])
		assert.Equal(t, lrCompID, entry["lr_component_id"])
		assert.Equal(t, lrCompType, entry["lr_component_type"])
		assert.Equal(t, network, entry["network"])
		assert.Equal(t, nextHop, entry["next_hop"])
		assert.Equal(t, nextHopGW, entry["next_hop_gateway"])
		assert.Equal(t, routeType, entry["route_type"])
	})
}
