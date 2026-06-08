// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerBgpRouteTable_basic(t *testing.T) {
	name := getAccTestResourceName()
	testDataSourceName := "data.nsxt_policy_route_controller_bgp_route_table.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpRouteTableDataSourceTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpRouteTableDataSourceTemplate(name string) string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller_bgp_route_table" "test" {
  route_controller_id            = nsxt_policy_route_controller.test.nsx_id
  virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path # placeholder
}
`, name)
}
