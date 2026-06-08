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

func TestMockResourceNsxtPolicyRouteControllerBgpTroubleshoot(t *testing.T) {
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

	t.Run("Create success", func(t *testing.T) {
		mockBgpTroubleSDK.EXPECT().Patch(rcID, gomock.Any()).Return(nil)

		mockBgpTroubleSDK.EXPECT().Get(rcID).Return(model.BgpTroubleshootConfig{
			Path:                     &path,
			Revision:                 &revision,
			BfdControlPktDiagnostics: &bfdDiag,
			BgpSessionDiagnostics:    &bgpDiag,
			SystemDiagnostics:        &sysDiag,
		}, nil)

		resource := resourceNsxtPolicyRouteControllerBgpTroubleshoot()
		d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
			"route_controller_path":       "/infra/route-controllers/" + rcID,
			"bfd_control_pkt_diagnostics": bfdDiag,
			"bgp_session_diagnostics":     bgpDiag,
			"system_diagnostics":          sysDiag,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpTroubleshootCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, path, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
	})

	t.Run("Read success", func(t *testing.T) {
		mockBgpTroubleSDK.EXPECT().Get(rcID).Return(model.BgpTroubleshootConfig{
			Path:                     &path,
			Revision:                 &revision,
			BfdControlPktDiagnostics: &bfdDiag,
			BgpSessionDiagnostics:    &bgpDiag,
			SystemDiagnostics:        &sysDiag,
		}, nil)

		resource := resourceNsxtPolicyRouteControllerBgpTroubleshoot()
		d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
		d.SetId(rcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpTroubleshootRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, path, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
	})

	t.Run("Update success", func(t *testing.T) {
		newSysDiag := true
		newRevision := int64(2)

		mockBgpTroubleSDK.EXPECT().Update(rcID, gomock.Any()).Return(model.BgpTroubleshootConfig{
			Path:                     &path,
			Revision:                 &newRevision,
			BfdControlPktDiagnostics: &bfdDiag,
			BgpSessionDiagnostics:    &bgpDiag,
			SystemDiagnostics:        &newSysDiag,
		}, nil)

		mockBgpTroubleSDK.EXPECT().Get(rcID).Return(model.BgpTroubleshootConfig{
			Path:                     &path,
			Revision:                 &newRevision,
			BfdControlPktDiagnostics: &bfdDiag,
			BgpSessionDiagnostics:    &bgpDiag,
			SystemDiagnostics:        &newSysDiag,
		}, nil)

		resource := resourceNsxtPolicyRouteControllerBgpTroubleshoot()
		d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
			"route_controller_path":       "/infra/route-controllers/" + rcID,
			"revision":                    int(revision),
			"bfd_control_pkt_diagnostics": bfdDiag,
			"bgp_session_diagnostics":     bgpDiag,
			"system_diagnostics":          newSysDiag,
		})
		d.SetId(rcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpTroubleshootUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, int(newRevision), d.Get("revision"))
		assert.Equal(t, newSysDiag, d.Get("system_diagnostics"))
	})

	t.Run("Delete success", func(t *testing.T) {
		mockBgpTroubleSDK.EXPECT().Delete(rcID).Return(nil)

		resource := resourceNsxtPolicyRouteControllerBgpTroubleshoot()
		d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
		d.SetId(rcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpTroubleshootDelete(d, m)
		require.NoError(t, err)
	})
}
