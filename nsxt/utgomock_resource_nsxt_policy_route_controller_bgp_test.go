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

func TestMockResourceNsxtPolicyRouteControllerBgp(t *testing.T) {
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
	rcPath := "/infra/route-controllers/rc-1"
	bgpPath := "/infra/route-controllers/rc-1/bgp"
	revision := int64(1)
	ecmp := true
	localAsNum := "65001"
	multipathRelax := true
	peerTimer := int64(5)
	mode := "HELPER_ONLY"
	restartTimer := int64(120)
	staleTimer := int64(600)

	t.Run("Create success", func(t *testing.T) {
		mockBgpSDK.EXPECT().Patch(rcID, gomock.Any()).Return(nil)
		mockBgpSDK.EXPECT().Get(rcID).Return(model.RouteControllerBgpRoutingConfig{
			Path:                      &bgpPath,
			Revision:                  &revision,
			Ecmp:                      &ecmp,
			LocalAsNum:                &localAsNum,
			MultipathRelax:            &multipathRelax,
			PeerRouteConvergenceTimer: &peerTimer,
			GracefulRestartConfig: &model.BgpGracefulRestartConfig{
				Mode: &mode,
				Timer: &model.BgpGracefulRestartTimer{
					RestartTimer:    &restartTimer,
					StaleRouteTimer: &staleTimer,
				},
			},
		}, nil)

		res := resourceNsxtPolicyRouteControllerBgp()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path":              rcPath,
			"ecmp":                               ecmp,
			"local_as_num":                       localAsNum,
			"multipath_relax":                    multipathRelax,
			"peer_route_convergence_timer":       int(peerTimer),
			"graceful_restart_mode":              mode,
			"graceful_restart_timer":             int(restartTimer),
			"graceful_restart_stale_route_timer": int(staleTimer),
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcID, d.Id())
		assert.Equal(t, bgpPath, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, ecmp, d.Get("ecmp"))
		assert.Equal(t, localAsNum, d.Get("local_as_num"))
		assert.Equal(t, multipathRelax, d.Get("multipath_relax"))
		assert.Equal(t, int(peerTimer), d.Get("peer_route_convergence_timer"))
		assert.Equal(t, mode, d.Get("graceful_restart_mode"))
		assert.Equal(t, int(restartTimer), d.Get("graceful_restart_timer"))
		assert.Equal(t, int(staleTimer), d.Get("graceful_restart_stale_route_timer"))
	})

	t.Run("Update success", func(t *testing.T) {
		mockBgpSDK.EXPECT().Update(rcID, gomock.Any()).Return(model.RouteControllerBgpRoutingConfig{
			Path:     &bgpPath,
			Revision: &revision,
		}, nil)
		mockBgpSDK.EXPECT().Get(rcID).Return(model.RouteControllerBgpRoutingConfig{
			Path:     &bgpPath,
			Revision: &revision,
		}, nil)

		res := resourceNsxtPolicyRouteControllerBgp()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
			"revision":              0,
		})
		d.SetId(rcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockBgpSDK.EXPECT().Delete(rcID).Return(nil)

		res := resourceNsxtPolicyRouteControllerBgp()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"route_controller_path": rcPath,
		})
		d.SetId(rcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpDelete(d, m)
		require.NoError(t, err)
	})
}
