//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/bgp/neighbors"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpNeighborsStatusClientContext utl.ClientContext

func NewRouteControllerBgpNeighborsStatusClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpNeighborsStatusClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewStatusClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpNeighborsStatusClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpNeighborsStatusClientContext) List(routeControllerIdParam string, bgpNeighborTypeParam, cursorParam, enforcementPointPathParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam, statsTypeParam, transportNodeIdParam, virtualNetworkAppliancePathParam *string) (model0.RouteControllerBgpNeighborsStatusListResult, error) {
	var obj model0.RouteControllerBgpNeighborsStatusListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.StatusClient)
		obj, err = client.List(routeControllerIdParam, bgpNeighborTypeParam, cursorParam, enforcementPointPathParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam, statsTypeParam, transportNodeIdParam, virtualNetworkAppliancePathParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
