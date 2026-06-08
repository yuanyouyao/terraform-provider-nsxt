//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpRouteTableClientContext utl.ClientContext

func NewRouteControllerBgpRouteTableClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpRouteTableClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewBgpRouteTableClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpRouteTableClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpRouteTableClientContext) List(routeControllerIdParam string, virtualNetworkAppliancePathParam string, networkPrefixParam *string) (model0.BgpRIBListResult, error) {
	var obj model0.BgpRIBListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.BgpRouteTableClient)
		obj, err = client.List(routeControllerIdParam, virtualNetworkAppliancePathParam, networkPrefixParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
