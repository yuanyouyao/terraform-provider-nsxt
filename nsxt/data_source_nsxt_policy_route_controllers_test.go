// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllers_basic(t *testing.T) {
	checkResourceName := "nsxt_policy_route_controllers_result"

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
				Config: testAccNSXPolicyRouteControllersReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "tf-acc-rc-list-01"),
				),
			},
			{
				Config: testAccNSXPolicyRouteControllersReadTemplateWithRegex(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "regex-tf-acc-rc-list-01"),
				),
			},
		},
	})
}

func testAccNSXPolicyRouteControllersReadTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + `
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "tf-acc-rc-list-01"
  nsx_id                                 = "tf-acc-rc-list-01"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controllers" "all" {
  depends_on = [nsxt_policy_route_controller.test]
}

output "nsxt_policy_route_controllers_result" {
  value      = data.nsxt_policy_route_controllers.all.items["tf-acc-rc-list-01"]
  depends_on = [data.nsxt_policy_route_controllers.all]
}
`
}

func testAccNSXPolicyRouteControllersReadTemplateWithRegex() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + `
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "regex-tf-acc-rc-list-01"
  nsx_id                                 = "regex-tf-acc-rc-list-01"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controllers" "all" {
  display_name = ".*"
  depends_on   = [nsxt_policy_route_controller.test]
}

output "nsxt_policy_route_controllers_result" {
  value      = data.nsxt_policy_route_controllers.all.items["regex-tf-acc-rc-list-01"]
  depends_on = [data.nsxt_policy_route_controllers.all]
}
`
}
