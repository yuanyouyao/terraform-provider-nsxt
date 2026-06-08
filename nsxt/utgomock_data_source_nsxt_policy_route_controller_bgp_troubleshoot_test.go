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
	bgptroublemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyRouteControllerBgpTroubleshoot(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpTroubleSDK := bgptroublemocks.NewMockTroubleshootClient(ctrl)
	mockWrapper := &cliinfra.RouteControllerBgpTroubleshootClientContext{
		Client:     mockBgpTroubleSDK,
		ClientType: utl.Local,
	}

	originalCli := cliRouteControllerBgpTroubleshootClient
	defer func() { cliRouteControllerBgpTroubleshootClient = originalCli }()
	cliRouteControllerBgpTroubleshootClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.RouteControllerBgpTroubleshootClientContext {
		return mockWrapper
	}

	rcID := "rc-1"
	path := "/infra/route-controllers/rc-1/bgp/troubleshoot"
	revision := int64(1)
	bfdDiag := true
	bgpDiag := true
	sysDiag := false

	t.Run("Read success", func(t *testing.T) {
		mockBgpTroubleSDK.EXPECT().Get(rcID).Return(model.BgpTroubleshootConfig{
			Path:                     &path,
			Revision:                 &revision,
			BfdControlPktDiagnostics: &bfdDiag,
			BgpSessionDiagnostics:    &bgpDiag,
			SystemDiagnostics:        &sysDiag,
		}, nil)

		ds := dataSourceNsxtPolicyRouteControllerBgpTroubleshoot()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"route_controller_id": rcID,
		})

		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyRouteControllerBgpTroubleshootRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, path, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, bfdDiag, d.Get("bfd_control_pkt_diagnostics"))
		assert.Equal(t, bgpDiag, d.Get("bgp_session_diagnostics"))
		assert.Equal(t, sysDiag, d.Get("system_diagnostics"))
	})
}
