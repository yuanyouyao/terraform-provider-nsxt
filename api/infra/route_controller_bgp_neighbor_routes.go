//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/bgp/neighbors"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpNeighborRoutesClientContext utl.ClientContext

func NewRouteControllerBgpNeighborRoutesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpNeighborRoutesClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewRoutesClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpNeighborRoutesClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpNeighborRoutesClientContext) List(routeControllerIdParam string, neighborIdParam string, countParam *int64, cursorParam *string, enforcementPointPathParam *string, includedFieldsParam *string, neighborAddressParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.RouteControllerBgpNeighborRoutesListResult, error) {
	var obj model0.RouteControllerBgpNeighborRoutesListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.RoutesClient)
		obj, err = client.List(routeControllerIdParam, neighborIdParam, countParam, cursorParam, enforcementPointPathParam, includedFieldsParam, neighborAddressParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
