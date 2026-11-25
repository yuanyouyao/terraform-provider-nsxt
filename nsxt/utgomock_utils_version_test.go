//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/node/VersionClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/node/VersionClient.go VersionClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	nodeAPI "github.com/vmware/terraform-provider-nsxt/api/nsx/node"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nodemocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/node"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockNsxtVersionOverride(t *testing.T) {
	previous := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = previous })
	ctrl := gomock.NewController(t)
	_, restore := setupVersionMock(ctrl)
	t.Cleanup(restore)

	// No Get expectation: configuring an override must never contact the version API.
	util.NsxVersion = "4.2.0"
	require.NoError(t, initNSXVersion(nil, "9.2.0.0.25956201"))
	assert.Equal(t, "9.2.0.0.25956201", util.NsxVersion)
	assert.True(t, util.NsxVersionHigherOrEqual("9.2.0"))
}

func TestMockNsxtVersionOverrideValidation(t *testing.T) {
	validate := Provider().Schema["nsx_version"].ValidateFunc
	for _, value := range []string{"", "9.2.0", "9.2.0.0.25956201"} {
		_, errs := validate(value, "nsx_version")
		assert.Empty(t, errs, value)
	}
	for _, value := range []string{"latest", "9.2", "9.2.0.", " 9.2.0", "9.2.0-beta"} {
		_, errs := validate(value, "nsx_version")
		assert.NotEmpty(t, errs, value)
	}
}

func TestMockNsxtVersionOverrideInitialization(t *testing.T) {
	previous := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = previous })
	ctrl := gomock.NewController(t)
	_, restore := setupVersionMock(ctrl)
	t.Cleanup(restore)
	t.Setenv("NSXT_VERSION", "9.2.0.0.25956201")

	for _, deferred := range []bool{false, true} {
		util.NsxVersion = ""
		data := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
			"host":                 "example.invalid",
			"username":             "admin",
			"password":             "test-only-password",
			"global_manager":       true,
			"on_demand_connection": deferred,
		})
		clients := &nsxtClients{}
		require.NoError(t, configurePolicyConnectorData(data, clients))
		if deferred {
			assert.Empty(t, util.NsxVersion)
			getPolicyConnector(*clients)
		}
		assert.Equal(t, "9.2.0.0.25956201", util.NsxVersion)
	}
}

func setupVersionMock(ctrl *gomock.Controller) (*nodemocks.MockVersionClient, func()) {
	mockSDK := nodemocks.NewMockVersionClient(ctrl)
	mockWrapper := &nodeAPI.NodeVersionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliVersionClient
	cliVersionClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *nodeAPI.NodeVersionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVersionClient = original }
}

func TestMockNsxtGetNSXVersion(t *testing.T) {
	t.Run("prefers ProductVersion when present", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		nodeVersion := "4.1.0.0.0"
		productVersion := "4.1.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion, ProductVersion: &productVersion}, nil)

		version, err := getNSXVersion(nil)
		require.NoError(t, err)
		assert.Equal(t, "4.1.0", version)
	})

	t.Run("falls back to NodeVersion when ProductVersion is absent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		nodeVersion := "4.1.0.0.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion}, nil)

		version, err := getNSXVersion(nil)
		require.NoError(t, err)
		assert.Equal(t, "4.1.0.0.0", version)
	})

	t.Run("Get error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{}, errors.New("connection error"))

		_, err := getNSXVersion(nil)
		require.Error(t, err)
	})
}

func TestMockNsxtInitNSXVersion(t *testing.T) {
	t.Run("success sets util.NsxVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		// NodeVersion must be non-nil: getNSXVersion unconditionally dereferences
		// it for a debug log before checking ProductVersion.
		nodeVersion := "4.1.0.0.0"
		productVersion := "4.1.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion, ProductVersion: &productVersion}, nil)

		err := initNSXVersion(nil, "")
		require.NoError(t, err)
	})

	t.Run("error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{}, errors.New("connection error"))

		err := initNSXVersion(nil, "")
		require.Error(t, err)
	})
}
