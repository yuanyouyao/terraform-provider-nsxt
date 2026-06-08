// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerBgpTroubleshoot_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_route_controller_bgp_troubleshoot.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpTroubleshootDataSourceTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "id"),
					resource.TestCheckResourceAttr(testDataSourceName, "system_diagnostics", "true"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpTroubleshootDataSourceTemplate() string {
	name := getAccTestResourceName()

	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller_bgp_troubleshoot" "test" {
  route_controller_path       = nsxt_policy_route_controller.test.path
  bfd_control_pkt_diagnostics = true
  bgp_session_diagnostics     = true
  system_diagnostics          = true
}

data "nsxt_policy_route_controller_bgp_troubleshoot" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
  depends_on          = [nsxt_policy_route_controller_bgp_troubleshoot.test]
}
`, name)
}
