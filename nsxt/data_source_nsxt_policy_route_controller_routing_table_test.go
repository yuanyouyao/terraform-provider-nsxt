// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerRoutingTable_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_route_controller_routing_table.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.1")
			testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyRouteControllerRoutingTableTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "route_controller_id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "routing_table_csv"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "routing_table.#"),
				),
			},
		},
	})
}

func testAccNSXPolicyRouteControllerRoutingTableTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + `
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "tf-acc-rc-rt-01"
  nsx_id                                 = "tf-acc-rc-rt-01"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller_routing_table" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
  depends_on          = [nsxt_policy_route_controller.test]
}
`
}
