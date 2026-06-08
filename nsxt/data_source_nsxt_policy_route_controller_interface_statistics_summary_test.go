// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerInterfaceStatisticsSummary_basic(t *testing.T) {
	name := getAccTestResourceName()
	testDataSourceName := "data.nsxt_policy_route_controller_interface_statistics_summary.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerInterfaceStatisticsSummaryDataSourceTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "interface_path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerInterfaceStatisticsSummaryDataSourceTemplate(name string) string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller_interface" "test" {
  route_controller_path = nsxt_policy_route_controller.test.path
  display_name          = "iface-test"
  description           = "terraform created"
  mtu                   = 1500
  urpf_mode             = "NONE"

  floating_ip_subnets {
    prefix_len   = 24
    ip_addresses = ["192.168.1.100"]
  }

  interface_address {
    portgroup_id                   = "dvportgroup-1"
    virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path # placeholder
    interface_subnet {
      prefix_len   = 24
      ip_addresses = ["192.168.1.1"]
    }
  }
}

data "nsxt_policy_route_controller_interface_statistics_summary" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
  interface_id        = nsxt_policy_route_controller_interface.test.nsx_id
}
`, name)
}
