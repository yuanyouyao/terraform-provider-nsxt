// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
}

func TestAccDataSourceNsxtPolicyRouteControllerBgpNeighbor_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_route_controller_bgp_neighbor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpNeighborDataSourceTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "neighbor_address", "192.168.1.1"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpNeighborDataSourceTemplate() string {
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
  display_name          = "neigh-test"
  neighbor_address      = "192.168.1.1"
  remote_as_num         = "65001"
  allow_as_in           = true
  enabled               = true
  hold_down_time        = 180
  keep_alive_time       = 60
  maximum_hop_limit     = 1
}

data "nsxt_policy_route_controller_bgp_neighbor" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
  id                  = nsxt_policy_route_controller_bgp_neighbor.test.nsx_id
}
`, name)
}
