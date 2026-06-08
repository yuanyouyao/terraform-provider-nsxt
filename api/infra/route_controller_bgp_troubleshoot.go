//nolint:revive
package infra

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers/bgp"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerBgpTroubleshootClientContext utl.ClientContext

func NewRouteControllerBgpTroubleshootClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RouteControllerBgpTroubleshootClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewTroubleshootClient(connector)
	default:
		return nil
	}
	return &RouteControllerBgpTroubleshootClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RouteControllerBgpTroubleshootClientContext) Delete(routerControllerIdParam string) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.TroubleshootClient)
		err = client.Delete(routerControllerIdParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpTroubleshootClientContext) Get(routerControllerIdParam string) (model0.BgpTroubleshootConfig, error) {
	var obj model0.BgpTroubleshootConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.TroubleshootClient)
		obj, err = client.Get(routerControllerIdParam)
		if err != nil {
			return obj, err
		}
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerBgpTroubleshootClientContext) Patch(routerControllerIdParam string, bgpTroubleshootConfigParam model0.BgpTroubleshootConfig) error {
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.TroubleshootClient)
		err = client.Patch(routerControllerIdParam, bgpTroubleshootConfigParam)
	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RouteControllerBgpTroubleshootClientContext) Update(routerControllerIdParam string, bgpTroubleshootConfigParam model0.BgpTroubleshootConfig) (model0.BgpTroubleshootConfig, error) {
	var obj model0.BgpTroubleshootConfig
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.TroubleshootClient)
		obj, err = client.Update(routerControllerIdParam, bgpTroubleshootConfigParam)
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
