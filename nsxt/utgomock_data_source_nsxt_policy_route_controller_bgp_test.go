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
	bgpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgp(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpSDK := bgpmocks.NewMockBgpClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpClientContext{
		Client:     mockBgpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpClient
	defer func() { cliRouteControllerBgpClient = originalCli }()
	cliRouteControllerBgpClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	bgpPath := "/infra/route-controllers/rc-1/bgp"
	revision := int64(1)
	ecmp := true
	localAsNum := "65001"

	t.Run("Read success", func(t *testing.T) {
		mockBgpSDK.EXPECT().Get(rcID).Return(model.RouteControllerBgpRoutingConfig{
			Path:       &bgpPath,
			Revision:   &revision,
			Ecmp:       &ecmp,
			LocalAsNum: &localAsNum,
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgp()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, bgpPath, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, ecmp, d.Get("ecmp"))
		assert.Equal(t, localAsNum, d.Get("local_as_num"))
	})
}
