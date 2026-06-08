// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerBgp_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_route_controller_bgp.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpDataSourceTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "local_as_num", "65001"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpDataSourceTemplate() string {
	name := getAccTestResourceName()

	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller_bgp" "test" {
  route_controller_path              = nsxt_policy_route_controller.test.path
  ecmp                               = true
  local_as_num                       = "65001"
  multipath_relax                    = true
  peer_route_convergence_timer       = 5
  graceful_restart_mode              = "HELPER_ONLY"
  graceful_restart_timer             = 120
  graceful_restart_stale_route_timer = 600
}

data "nsxt_policy_route_controller_bgp" "test" {
  route_controller_id = nsxt_policy_route_controller_bgp.test.id
}
`, name)
}
