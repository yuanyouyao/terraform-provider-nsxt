// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceNsxtPolicyRouteControllerInterface_basic(t *testing.T) {
	name := getAccTestResourceName()
	resourceName := "nsxt_policy_route_controller_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerInterfaceTemplate(name, "iface-test", "terraform created", 1500),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "iface-test"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(resourceName, "urpf_mode", "NONE"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerInterfaceTemplate(name, "iface-test-updated", "terraform updated", 1400),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "iface-test-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform updated"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1400"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerInterfaceTemplate(name, displayName, description string, mtu int) string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "terraform created"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller_interface" "test" {
  route_controller_path = nsxt_policy_route_controller.test.path
  display_name          = "%s"
  description           = "%s"
  mtu                   = %d
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
`, name, displayName, description, mtu)
}
