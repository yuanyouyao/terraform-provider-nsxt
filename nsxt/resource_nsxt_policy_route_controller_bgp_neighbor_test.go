// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyRouteControllerBgpNeighbor_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_bgp_neighbor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerBgpNeighborCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpNeighborResourceTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerBgpNeighborExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", "neigh-test"),
					resource.TestCheckResourceAttr(testResourceName, "neighbor_address", "192.168.1.1"),
					resource.TestCheckResourceAttr(testResourceName, "remote_as_num", "65001"),
					resource.TestCheckResourceAttr(testResourceName, "allow_as_in", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "hold_down_time", "180"),
					resource.TestCheckResourceAttr(testResourceName, "keep_alive_time", "60"),
					resource.TestCheckResourceAttr(testResourceName, "maximum_hop_limit", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerBgpNeighborResourceTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerBgpNeighborExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", "neigh-test-updated"),
					resource.TestCheckResourceAttr(testResourceName, "neighbor_address", "192.168.1.1"),
					resource.TestCheckResourceAttr(testResourceName, "remote_as_num", "65001"),
					resource.TestCheckResourceAttr(testResourceName, "allow_as_in", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "hold_down_time", "150"),
					resource.TestCheckResourceAttr(testResourceName, "keep_alive_time", "50"),
					resource.TestCheckResourceAttr(testResourceName, "maximum_hop_limit", "2"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpNeighborExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		sessionContext := testAccGetSessionContext()

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor resource %s not found in resources", resourceName)
		}

		idStr := rs.Primary.ID
		segs := strings.Split(idStr, "/")
		if len(segs) < 2 {
			return fmt.Errorf("invalid resource ID format: %s", idStr)
		}
		routeControllerID := segs[0]
		id := segs[1]

		client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err := client.Get(routeControllerID, id)
		if err != nil {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor with ID %s does not exist: %v", idStr, err)
		}

		return nil
	}
}

func testAccNsxtPolicyRouteControllerBgpNeighborCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	sessionContext := testAccGetSessionContext()

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller_bgp_neighbor" {
			continue
		}

		idStr := rs.Primary.ID
		segs := strings.Split(idStr, "/")
		if len(segs) < 2 {
			return fmt.Errorf("invalid resource ID format: %s", idStr)
		}
		routeControllerID := segs[0]
		id := segs[1]

		client := cliRouteControllerBgpNeighborClient(sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err := client.Get(routeControllerID, id)
		if err == nil {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor %s still exists", idStr)
		}
	}

	return nil
}

func testAccNsxtPolicyRouteControllerBgpNeighborResourceTemplate(isCreate bool) string {
	displayName := "neigh-test"
	allowAsIn := "true"
	enabled := "true"
	holdDown := "180"
	keepAlive := "60"
	maxHop := "1"

	if !isCreate {
		displayName = "neigh-test-updated"
		allowAsIn = "false"
		enabled = "false"
		holdDown = "150"
		keepAlive = "50"
		maxHop = "2"
	}

	name := getAccTestResourceName()

	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  route_controller_path = nsxt_policy_route_controller.test.path
  display_name          = "%s"
  neighbor_address      = "192.168.1.1"
  remote_as_num         = "65001"
  allow_as_in           = %s
  enabled               = %s
  hold_down_time        = %s
  keep_alive_time       = %s
  maximum_hop_limit     = %s
}
`, name, displayName, allowAsIn, enabled, holdDown, keepAlive, maxHop)
}
