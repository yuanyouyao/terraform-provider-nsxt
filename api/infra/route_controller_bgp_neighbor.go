//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/bgp"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpNeighborClientContext utl.ClientContext

func NewRouteControllerBgpNeighborClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpNeighborClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewNeighborsClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpNeighborClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpNeighborClientContext) Delete(routerControllerIdParam string, neighborIdParam string) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.NeighborsClient)
		err = client.Delete(routerControllerIdParam, neighborIdParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpNeighborClientContext) Get(routerControllerIdParam string, neighborIdParam string) (model0.RouteControllerBgpNeighborConfig, error) {
	var obj model0.RouteControllerBgpNeighborConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.NeighborsClient)
		obj, err = client.Get(routerControllerIdParam, neighborIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerBgpNeighborClientContext) List(routerControllerIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.RouteControllerBgpNeighborConfigListResult, error) {
	var obj model0.RouteControllerBgpNeighborConfigListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.NeighborsClient)
		obj, err = client.List(routerControllerIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerBgpNeighborClientContext) Patch(routerControllerIdParam string, neighborIdParam string, routeControllerBgpNeighborConfigParam model0.RouteControllerBgpNeighborConfig) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.NeighborsClient)
		err = client.Patch(routerControllerIdParam, neighborIdParam, routeControllerBgpNeighborConfigParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpNeighborClientContext) Update(routerControllerIdParam string, neighborIdParam string, routeControllerBgpNeighborConfigParam model0.RouteControllerBgpNeighborConfig) (model0.RouteControllerBgpNeighborConfig, error) {
	var obj model0.RouteControllerBgpNeighborConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.NeighborsClient)
		obj, err = client.Update(routerControllerIdParam, neighborIdParam, routeControllerBgpNeighborConfigParam)
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
