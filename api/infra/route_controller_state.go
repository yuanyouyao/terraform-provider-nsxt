//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerStateClientContext utl.ClientContext

func NewRouteControllerStateClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerStateClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewStateClient(connector)

	default:
		return nil
	}
	return &RouteControllerStateClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerStateClientContext) Get(routeControllerIdParam string, sourceParam *string) (model0.RouteControllerState, error) {
	var obj model0.RouteControllerState
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.StateClient)
		obj, err = client.Get(routeControllerIdParam, sourceParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
