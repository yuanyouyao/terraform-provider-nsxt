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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ifacesmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockResourceNsxtPolicyRouteControllerInterface(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfacesSDK := ifacesmocks.NewMockInterfacesClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerInterfacesClientContext{
		Client:     mockIfacesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerInterfacesClientForResource
	defer func() { cliRouteControllerInterfacesClientForResource = originalCli }()
	cliRouteControllerInterfacesClientForResource = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerInterfacesClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	rcPath := "/infra/route-controllers/rc-1"
	ifaceID := "iface-1"
	ifacePath := "/infra/route-controllers/rc-1/interfaces/iface-1"
	revision := int64(1)
	displayName := "iface-test"
	description := "mock interface"
	mtu := int64(1500)
	urpfMode := "NONE"

	portgroupID := "dvportgroup-1"
	vnaPath := "/infra/sites/site-1/enforcement-points/ep-1/virtual-network-appliances/vna-1"
	ipAddress := "192.168.1.1"
	prefixLen := int64(24)

	t.Run("Create success", func(t *testing.T) {
		mockIfacesSDK.EXPECT().Get(rcID, ifaceID).Return(model.RouteControllerInterface{}, vapiErrors.NotFound{})
		mockIfacesSDK.EXPECT().Patch(rcID, ifaceID, gomock.Any()).Return(nil)
		mockIfacesSDK.EXPECT().Get(rcID, ifaceID).Return(model.RouteControllerInterface{
			Id:          &ifaceID,
			Path:        &ifacePath,
			Revision:    &revision,
			DisplayName: &displayName,
			Description: &description,
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
		}, nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"nsx_id":                ifaceID,
			"display_name":          displayName,
			"description":           description,
			"mtu":                   int(mtu),
			"urpf_mode":             urpfMode,
			"floating_ip_subnets": []interface{}{
				map[string]interface{}{
					"prefix_len":   int(prefixLen),
					"ip_addresses": []interface{}{"192.168.1.100"},
				},
			},
			"interface_address": []interface{}{
				map[string]interface{}{
					"portgroup_id":                   portgroupID,
					"virtual_network_appliance_path": vnaPath,
					"interface_subnet": []interface{}{
						map[string]interface{}{
							"prefix_len":   int(prefixLen),
							"ip_addresses": []interface{}{ipAddress},
						},
					},
				},
			},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID+"/"+ifaceID, d.Id())
		assert.Equal(t, ifacePath, d.Get("path"))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, int(mtu), d.Get("mtu"))
		assert.Equal(t, urpfMode, d.Get("urpf_mode"))
	})

	t.Run("Update success", func(t *testing.T) {
		mockIfacesSDK.EXPECT().Update(rcID, ifaceID, gomock.Any()).Return(model.RouteControllerInterface{
			Id:          &ifaceID,
			Path:        &ifacePath,
			Revision:    &revision,
			DisplayName: &displayName,
			Description: &description,
			Mtu:         &mtu,
			UrpfMode:    &urpfMode,
		}, nil)
		mockIfacesSDK.EXPECT().Get(rcID, ifaceID).Return(model.RouteControllerInterface{
			Id:          &ifaceID,
			Path:        &ifacePath,
			Revision:    &revision,
			DisplayName: &displayName,
			Description: &description,
			Mtu:         &mtu,
			UrpfMode:    &urpfMode,
		}, nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"nsx_id":                ifaceID,
			"display_name":          displayName,
			"description":           description,
			"mtu":                   int(mtu),
			"urpf_mode":             urpfMode,
			"interface_address": []interface{}{
				map[string]interface{}{
					"portgroup_id":                   portgroupID,
					"virtual_network_appliance_path": vnaPath,
					"interface_subnet": []interface{}{
						map[string]interface{}{
							"prefix_len":   int(prefixLen),
							"ip_addresses": []interface{}{ipAddress},
						},
					},
				},
			},
		})
		d.SetId(rcID + "/" + ifaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockIfacesSDK.EXPECT().Delete(rcID, ifaceID).Return(nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(rcID + "/" + ifaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceDelete(d, m)
		require.NoError(t, err)
	})
}
