// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerState_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_route_controller_state.test"

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
				Config: testAccNSXPolicyRouteControllerStateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "route_controller_id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "logical_gateway_id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "per_node_status.#"),
				),
			},
		},
	})
}

func testAccNSXPolicyRouteControllerStateTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + `
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "tf-acc-rc-state-01"
  nsx_id                                 = "tf-acc-rc-state-01"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller_state" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
  depends_on          = [nsxt_policy_route_controller.test]
}
`
}
