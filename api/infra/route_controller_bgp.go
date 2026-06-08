//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpClientContext utl.ClientContext

func NewRouteControllerBgpClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewBgpClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpClientContext) Delete(routerControllerIdParam string) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.BgpClient)
		err = client.Delete(routerControllerIdParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpClientContext) Get(routerControllerIdParam string) (model0.RouteControllerBgpRoutingConfig, error) {
	var obj model0.RouteControllerBgpRoutingConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.BgpClient)
		obj, err = client.Get(routerControllerIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerBgpClientContext) Patch(routerControllerIdParam string, routeControllerBgpRoutingConfigParam model0.RouteControllerBgpRoutingConfig) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.BgpClient)
		err = client.Patch(routerControllerIdParam, routeControllerBgpRoutingConfigParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpClientContext) Update(routerControllerIdParam string, routeControllerBgpRoutingConfigParam model0.RouteControllerBgpRoutingConfig) (model0.RouteControllerBgpRoutingConfig, error) {
	var obj model0.RouteControllerBgpRoutingConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.BgpClient)
		obj, err = client.Update(routerControllerIdParam, routeControllerBgpRoutingConfigParam)
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
