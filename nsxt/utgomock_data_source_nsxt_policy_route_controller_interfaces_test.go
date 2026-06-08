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
	ifacesmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerInterfaces(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfacesSDK := ifacesmocks.NewMockInterfacesClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerInterfacesClientContext{
		Client:     mockIfacesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerInterfacesClient
	defer func() { cliRouteControllerInterfacesClient = originalCli }()
	cliRouteControllerInterfacesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerInterfacesClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	ifaceID := "iface-1"
	displayName := "iface-test"
	description := "terraform created"
	path := "/infra/route-controllers/rc-1/interfaces/iface-1"
	revision := int64(1)
	mtu := int64(1500)
	urpfMode := "NONE"

	portgroupID := "dvportgroup-1"
	vnaPath := "/infra/sites/site-1/enforcement-points/ep-1/virtual-network-appliances/vna-1"
	ipAddress := "192.168.1.1"
	prefixLen := int64(24)

	t.Run("List success", func(t *testing.T) {
		mockIfacesSDK.EXPECT().List(rcID, nil, nil, nil, nil, nil, nil).Return(model.RouteControllerInterfaceListResult{
			Results: []model.RouteControllerInterface{
				{
					Id:          &ifaceID,
					DisplayName: &displayName,
					Description: &description,
					Path:        &path,
					Revision:    &revision,
					Mtu:         &mtu,
					UrpfMode:    &urpfMode,
					FloatingIpSubnets: []model.InterfaceSubnet{
						{
							IpAddresses: []string{"192.168.1.100"},
							PrefixLen:   &prefixLen,
						},
					},
					InterfaceAddress: []model.RouteControllerInterfaceAddress{
						{
							PortgroupId:                 &portgroupID,
							VirtualNetworkAppliancePath: &vnaPath,
							InterfaceSubnet: []model.InterfaceSubnet{
								{
									IpAddresses: []string{ipAddress},
									PrefixLen:   &prefixLen,
								},
							},
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerInterfaces()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerInterfacesRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())

		results := d.Get("results").([]interface{})
		assert.Len(t, results, 1)

		item := results[0].(map[string]interface{})
		assert.Equal(t, ifaceID, item["id"])
		assert.Equal(t, displayName, item["display_name"])
		assert.Equal(t, description, item["description"])
		assert.Equal(t, path, item["path"])
		assert.Equal(t, int(revision), item["revision"])
		assert.Equal(t, int(mtu), item["mtu"])
		assert.Equal(t, urpfMode, item["urpf_mode"])

		floatingSubnets := item["floating_ip_subnets"].([]interface{})
		assert.Len(t, floatingSubnets, 1)
		subItem := floatingSubnets[0].(map[string]interface{})
		assert.Equal(t, int(prefixLen), subItem["prefix_len"])
		assert.Equal(t, []interface{}{"192.168.1.100"}, subItem["ip_addresses"])

		interfaceAddress := item["interface_address"].([]interface{})
		assert.Len(t, interfaceAddress, 1)
		addrItem := interfaceAddress[0].(map[string]interface{})
		assert.Equal(t, portgroupID, addrItem["portgroup_id"])
		assert.Equal(t, vnaPath, addrItem["virtual_network_appliance_path"])

		subnets := addrItem["interface_subnet"].([]interface{})
		assert.Len(t, subnets, 1)
		subItem2 := subnets[0].(map[string]interface{})
		assert.Equal(t, int(prefixLen), subItem2["prefix_len"])
		assert.Equal(t, []interface{}{ipAddress}, subItem2["ip_addresses"])
	})
}
