//nolint:revive
package infra

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/route_controllers"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RouteControllerForwardingTableClientContext struct {
	Client     interface{}
	ClientType utl.ClientType
	ProjectID  string
	VPCID      string
	Host       string
	HTTPClient *http.Client
	Username   string
	Password   string
	Token      string
	Cookie     string
	XSRF       string
}

func NewRouteControllerForwardingTableClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector, host string, httpClient *http.Client, username, password, token, cookie, xsrf string) *RouteControllerForwardingTableClientContext {
	var client interface{}

	switch sessionContext.ClientType {
	case utl.Local:
		client = client0.NewForwardingTableClient(connector)
	default:
		return nil
	}
	return &RouteControllerForwardingTableClientContext{
		Client:     client,
		ClientType: sessionContext.ClientType,
		ProjectID:  sessionContext.ProjectID,
		VPCID:      sessionContext.VPCID,
		Host:       host,
		HTTPClient: httpClient,
		Username:   username,
		Password:   password,
		Token:      token,
		Cookie:     cookie,
		XSRF:       xsrf,
	}
}

func (c RouteControllerForwardingTableClientContext) Get(routeControllerIdParam string, networkPrefixParam *string, routeSourceParam *string, virtualNetworkAppliancePathParam *string) (model0.RouteControllerRoutingTableListResult, error) {
	var obj model0.RouteControllerRoutingTableListResult
	var err error

	switch c.ClientType {
	case utl.Local:
		client := c.Client.(client0.ForwardingTableClient)
		obj, err = client.Get(routeControllerIdParam, networkPrefixParam, routeSourceParam, virtualNetworkAppliancePathParam)
	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RouteControllerForwardingTableClientContext) Download(routeControllerIdParam string, networkPrefixParam *string, routeSourceParam *string, virtualNetworkAppliancePathParam *string) (string, error) {
	urlStr := fmt.Sprintf("https://%s/policy/api/v1/infra/route-controllers/%s/forwarding-table/download", c.Host, routeControllerIdParam)

	params := url.Values{}
	if networkPrefixParam != nil {
		params.Add("network_prefix", *networkPrefixParam)
	}
	if routeSourceParam != nil {
		params.Add("route_source", *routeSourceParam)
	}
	if virtualNetworkAppliancePathParam != nil {
		params.Add("virtual_network_appliance_path", *virtualNetworkAppliancePathParam)
	}
	if len(params) > 0 {
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}
	if c.XSRF != "" {
		req.Header.Set("X-XSRF-TOKEN", c.XSRF)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to download forwarding table, status: %s, body: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
