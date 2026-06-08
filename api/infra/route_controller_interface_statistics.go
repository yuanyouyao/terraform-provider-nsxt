//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/interfaces"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerInterfaceStatisticsClientContext utl.ClientContext

func NewRouteControllerInterfaceStatisticsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerInterfaceStatisticsClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewStatisticsClient(connector)
	default:
		return nil
	}
	return &RouteControllerInterfaceStatisticsClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerInterfaceStatisticsClientContext) Get(routeControllerIdParam, interfaceIdParam string, cursorParam *string, edgePathParam *string, enforcementPointPathParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, sourceParam *string, statsTypeParam *string, transportNodeIdParam *string) (model0.RouteControllerInterfaceStatistics, error) {
	var obj model0.RouteControllerInterfaceStatistics
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.StatisticsClient)
		obj, err = client.Get(routeControllerIdParam, interfaceIdParam, cursorParam, edgePathParam, enforcementPointPathParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam, sourceParam, statsTypeParam, transportNodeIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
