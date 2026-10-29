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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtVpcLBProfilePagination(t *testing.T) {
	type profile struct{ id, name string }
	type page struct {
		profiles []profile
		cursor   *string
		unnamed  bool
		err      error
	}
	specs := []struct {
		name          string
		newDataSource func() *schema.Resource
		typeName      string
		resourceTypes string
		newValue      func(t *testing.T, id, name string) *data.StructValue
	}{
		{
			name: "App", newDataSource: dataSourceNsxtVpcLBAppProfile, typeName: "TCP",
			resourceTypes: "(LBFastTcpProfile OR LBFastUdpProfile OR LBHttpProfile)",
			newValue: func(t *testing.T, id, name string) *data.StructValue {
				return lbAppProfileStructValue(t, id, name, model.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE)
			},
		},
		{
			name: "Monitor", newDataSource: dataSourceNsxtVpcLBMonitor, typeName: "TCP",
			resourceTypes: "(LBHttpMonitorProfile OR LBHttpsMonitorProfile OR LBIcmpMonitorProfile OR LBPassiveMonitorProfile OR LBTcpMonitorProfile OR LBUdpMonitorProfile)",
			newValue: func(t *testing.T, id, name string) *data.StructValue {
				return lbMonitorProfileStructValue(t, id, name, model.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE)
			},
		},
		{
			name: "Persistence", newDataSource: dataSourceNsxtVpcLBPersistenceProfile, typeName: "COOKIE",
			resourceTypes: "(LBCookiePersistenceProfile OR LBGenericPersistenceProfile OR LBSourceIpPersistenceProfile)",
			newValue: func(t *testing.T, id, name string) *data.StructValue {
				return vpcLbDsPersistValue(t, id, name, model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)
			},
		},
		{
			name: "ClientSsl", newDataSource: dataSourceNsxtVpcLBClientSslProfile,
			resourceTypes: "LBClientSslProfile",
			newValue: func(t *testing.T, id, name string) *data.StructValue {
				return lbPersistenceProfileStructValue(t, vpcLbDsClientSslItem(id, name), model.LBClientSslProfileBindingType())
			},
		},
		{
			name: "ServerSsl", newDataSource: dataSourceNsxtVpcLBServerSslProfile,
			resourceTypes: "LBServerSslProfile",
			newValue: func(t *testing.T, id, name string) *data.StructValue {
				return lbPersistenceProfileStructValue(t, vpcLbDsServerSslItem(id, name), model.LBServerSslProfileBindingType())
			},
		},
	}
	next, last, empty := "next-page", "last-page", ""
	cases := []struct {
		name          string
		pages         []page
		inputName     string
		typeOnly      bool
		requestedType string
		wantID        string
		errMsg        string
	}{
		{
			name: "Empty search response", inputName: "target", errMsg: "not found",
			pages: []page{{}},
		},
		{
			name: "ANY type-only match on later page", typeOnly: true, requestedType: "ANY", wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "ANY type-only ambiguity across pages", typeOnly: true, requestedType: "ANY", errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "first"}}, cursor: &next},
				{profiles: []profile{{"two", "second"}}},
			},
		},
		{
			name: "Skip unnamed profiles before a later name match", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"unnamed", ""}}, unnamed: true, cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Type-only lookup matches unnamed profiles", typeOnly: true, wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", ""}}, unnamed: true},
			},
		},
		{
			name: "Type-only lookup matches empty display names once", typeOnly: true, wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", ""}}},
			},
		},
		{
			name: "Literal name with spaces and special characters", inputName: "Production Pool (east)+[blue]/svc:443*?", wantID: "target-id",
			pages: []page{{profiles: []profile{{"target-id", "Production Pool (east)+[blue]/svc:443*?"}}}},
		},
		{
			name: "Literal prefix with spaces and special characters", inputName: "Production Pool (east)+[blue]/svc:443*?", wantID: "target-id",
			pages: []page{{profiles: []profile{{"target-id", "Production Pool (east)+[blue]/svc:443*?-suffix"}}}},
		},
		{
			name: "Exact match on third page", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"other", "other"}}, cursor: &next},
				{cursor: &last},
				{profiles: []profile{{"target-id", "target"}}, cursor: &empty},
			},
		},
		{
			name: "Prefix match on later page", inputName: "target", wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", "target-long"}}},
			},
		},
		{
			name: "Later exact match beats earlier prefix", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"prefix-id", "target-long"}}, cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Duplicate exact matches across pages", inputName: "target", errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "target"}}, cursor: &next},
				{profiles: []profile{{"two", "target"}}},
			},
		},
		{
			name: "Duplicate prefix matches across pages", inputName: "target", errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "target-one"}}, cursor: &next},
				{profiles: []profile{{"two", "target-two"}}},
			},
		},
		{
			name: "Later page error after a match", inputName: "target", errMsg: "Failed to read",
			pages: []page{
				{profiles: []profile{{"one", "target"}}, cursor: &next},
				{err: vapiErrors.InternalServerError{}},
			},
		},
		{
			name: "Not found after all pages", inputName: "target", errMsg: "not found",
			pages: []page{
				{profiles: []profile{{"one", "other"}}, cursor: &next},
				{profiles: []profile{{"two", "another"}}},
			},
		},
		{
			name: "Type match on later page", typeOnly: true, wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Duplicate type matches across pages", typeOnly: true, errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "first"}}, cursor: &next},
				{profiles: []profile{{"two", "second"}}},
			},
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			for _, tc := range cases {
				if tc.typeOnly && spec.typeName == "" {
					continue
				}
				t.Run(tc.name, func(t *testing.T) {
					vpcLbPrfSetNsxVersion(t, "9.0.0")
					var total int64
					for _, pg := range tc.pages {
						total += int64(len(pg.profiles))
						if pg.err != nil {
							total++
						}
					}
					stub := &seqQueryListClient{}
					var cursors []*string
					var cursor *string
					for _, pg := range tc.pages {
						results := make([]*data.StructValue, 0, len(pg.profiles))
						for _, p := range pg.profiles {
							item := spec.newValue(t, p.id, p.name)
							if pg.unnamed {
								item.SetField("display_name", data.NewOptionalValue(nil))
							}
							results = append(results, item)
						}
						stub.responses = append(stub.responses, model.SearchResponse{
							Results: results, ResultCount: &total, Cursor: pg.cursor,
						})
						stub.errs = append(stub.errs, pg.err)
						cursors = append(cursors, cursor)
						cursor = pg.cursor
					}
					t.Cleanup(setupCliQueryClientStub(t, stub))
					raw := map[string]interface{}{"display_name": tc.inputName}
					if tc.typeOnly {
						raw["type"] = spec.typeName
						if tc.requestedType != "" {
							raw["type"] = tc.requestedType
						}
					}
					ds := spec.newDataSource()
					d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbWithContext(raw))
					err := ds.Read(d, newGoMockProviderClient())
					expectedQuery := "resource_type:" + spec.resourceTypes +
						` AND marked_for_delete:false AND path:\/orgs\/default\/projects\/proj1\/vpcs\/vpc1\/*`
					require.Equal(t, len(tc.pages), stub.call)
					assert.Equal(t, cursors, stub.cursors)
					for _, query := range stub.queries {
						assert.Equal(t, expectedQuery, query)
					}
					if tc.errMsg != "" {
						require.Error(t, err)
						assert.Contains(t, err.Error(), tc.errMsg)
						return
					}
					require.NoError(t, err)
					assert.Equal(t, tc.wantID, d.Id())
				})
			}
		})
	}
}
