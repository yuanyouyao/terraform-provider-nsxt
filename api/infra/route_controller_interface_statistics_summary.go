//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/interfaces/statistics"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerInterfaceStatisticsSummaryClientContext utl.ClientContext

func NewRouteControllerInterfaceStatisticsSummaryClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerInterfaceStatisticsSummaryClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewSummaryClient(connector)
	default:
		return nil
	}
	return &RouteControllerInterfaceStatisticsSummaryClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerInterfaceStatisticsSummaryClientContext) Get(routeControllerIdParam, interfaceIdParam string) (model0.RouteControllerInterfaceStatisticsSummary, error) {
	var obj model0.RouteControllerInterfaceStatisticsSummary
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.SummaryClient)
		obj, err = client.Get(routeControllerIdParam, interfaceIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
