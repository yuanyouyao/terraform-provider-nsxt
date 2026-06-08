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

func TestMockDataSourceNsxtPolicyRouteControllerInterface(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInterfacesSDK := rtmocks.NewMockInterfacesClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerInterfacesClientContext{
		Client:     mockInterfacesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerInterfacesClientForDataSource
	defer func() { cliRouteControllerInterfacesClientForDataSource = originalCli }()
	cliRouteControllerInterfacesClientForDataSource = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerInterfacesClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	ifaceID := "iface-1"
	ifacePath := "/infra/route-controllers/rc-1/interfaces/iface-1"
	displayName := "iface-test"
	description := "test desc"
	revision := int64(1)
	mtu := int64(1500)
	urpfMode := "NONE"

	prefixLen := int64(24)
	ipAddrs := []string{"192.168.1.100"}
	floatingSubnet := model.InterfaceSubnet{
		PrefixLen:   &prefixLen,
		IpAddresses: ipAddrs,
	}

	portgroupID := "dvportgroup-1"
	vnaPath := "/infra/virtual-network-appliance-clusters/vna-1"
	ifaceSubnetPrefixLen := int64(24)
	ifaceSubnetIps := []string{"192.168.1.1"}
	ifaceSubnet := model.InterfaceSubnet{
		PrefixLen:   &ifaceSubnetPrefixLen,
		IpAddresses: ifaceSubnetIps,
	}
	ifaceAddr := model.RouteControllerInterfaceAddress{
		PortgroupId:                 &portgroupID,
		VirtualNetworkAppliancePath: &vnaPath,
		InterfaceSubnet:             []model.InterfaceSubnet{ifaceSubnet},
	}

	t.Run("Get success", func(t *testing.T) {
		mockInterfacesSDK.EXPECT().Get(rcID, ifaceID).Return(model.RouteControllerInterface{
			Path:              &ifacePath,
			DisplayName:       &displayName,
			Description:       &description,
			Revision:          &revision,
			Mtu:               &mtu,
			UrpfMode:          &urpfMode,
			FloatingIpSubnets: []model.InterfaceSubnet{floatingSubnet},
			InterfaceAddress:  []model.RouteControllerInterfaceAddress{ifaceAddr},
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
			"id":                  ifaceID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerInterfaceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+ifaceID, d.Id())
		assert.Equal(t, ifacePath, d.Get("path"))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, int(mtu), d.Get("mtu"))
		assert.Equal(t, urpfMode, d.Get("urpf_mode"))

		floatingList := d.Get("floating_ip_subnets").([]interface{})
		require.Len(t, floatingList, 1)
		floatingItem := floatingList[0].(map[string]interface{})
		assert.Equal(t, int(prefixLen), floatingItem["prefix_len"])
		assert.Equal(t, ipAddrs[0], floatingItem["ip_addresses"].([]interface{})[0])

		addrList := d.Get("interface_address").([]interface{})
		require.Len(t, addrList, 1)
		addrItem := addrList[0].(map[string]interface{})
		assert.Equal(t, portgroupID, addrItem["portgroup_id"])
		assert.Equal(t, vnaPath, addrItem["virtual_network_appliance_path"])

		subnetList := addrItem["interface_subnet"].([]interface{})
		require.Len(t, subnetList, 1)
		subnetItem := subnetList[0].(map[string]interface{})
		assert.Equal(t, int(ifaceSubnetPrefixLen), subnetItem["prefix_len"])
		assert.Equal(t, ifaceSubnetIps[0], subnetItem["ip_addresses"].([]interface{})[0])
	})
}
