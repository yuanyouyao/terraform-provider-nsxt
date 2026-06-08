//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerInterfacesClientContext utl.ClientContext

func NewRouteControllerInterfacesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerInterfacesClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewInterfacesClient(connector)
	default:
		return nil
	}
	return &RouteControllerInterfacesClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerInterfacesClientContext) List(routerControllerIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.RouteControllerInterfaceListResult, error) {
	var obj model0.RouteControllerInterfaceListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.InterfacesClient)
		obj, err = client.List(routerControllerIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerInterfacesClientContext) Get(routerControllerIdParam, interfaceIdParam string) (model0.RouteControllerInterface, error) {
	var obj model0.RouteControllerInterface
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.InterfacesClient)
		obj, err = client.Get(routerControllerIdParam, interfaceIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerInterfacesClientContext) Delete(routerControllerIdParam, interfaceIdParam string) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.InterfacesClient)
		err = client.Delete(routerControllerIdParam, interfaceIdParam)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerInterfacesClientContext) Patch(routerControllerIdParam, interfaceIdParam string, routeControllerInterfaceParam model0.RouteControllerInterface) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.InterfacesClient)
		err = client.Patch(routerControllerIdParam, interfaceIdParam, routeControllerInterfaceParam)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerInterfacesClientContext) Update(routerControllerIdParam, interfaceIdParam string, routeControllerInterfaceParam model0.RouteControllerInterface) (model0.RouteControllerInterface, error) {
	var obj model0.RouteControllerInterface
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.InterfacesClient)
		obj, err = client.Update(routerControllerIdParam, interfaceIdParam, routeControllerInterfaceParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
