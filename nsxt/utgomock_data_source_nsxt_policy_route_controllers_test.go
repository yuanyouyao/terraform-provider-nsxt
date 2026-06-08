//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtPolicyRouteControllersRead(t *testing.T) {
	rt := "RouteController"
	rcName := "rc-fooname"
	rcPath := "/infra/route-controllers/rc-1"
	rcID := "rc-1"

	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id:           &rcID,
		DisplayName:  &rcName,
		Path:         &rcPath,
		ResourceType: &rt,
	})

	t.Run("success", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results:     []*data.StructValue{sv},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyRouteControllersRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		items := d.Get("items").(map[string]interface{})
		require.Contains(t, items, rcID)
		assert.Equal(t, rcName, items[rcID].(string))
		assert.NotEmpty(t, d.Id())
	})

	t.Run("success with regex match", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results:     []*data.StructValue{sv},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ".*foo.*",
		})
		err := dataSourceNsxtPolicyRouteControllersRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		items := d.Get("items").(map[string]interface{})
		require.Contains(t, items, rcID)
		assert.Equal(t, rcName, items[rcID].(string))
	})

	t.Run("success with regex non-match", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results:     []*data.StructValue{sv},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "bar",
		})
		err := dataSourceNsxtPolicyRouteControllersRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		items := d.Get("items").(map[string]interface{})
		assert.Empty(t, items)
	})

	t.Run("list error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search fail")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyRouteControllersRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error in listing the Route Controller items")
	})
}
