// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccNsxtPolicyRouteControllerBgpPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
}

func TestAccResourceNsxtPolicyRouteControllerBgp_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_bgp.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerBgpPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerBgpCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerBgpResourceTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerBgpExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "local_as_num", "65001"),
					resource.TestCheckResourceAttr(testResourceName, "multipath_relax", "true"),
					resource.TestCheckResourceAttr(testResourceName, "peer_route_convergence_timer", "5"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_timer", "120"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_stale_route_timer", "600"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerBgpResourceTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerBgpExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "ecmp", "false"),
					resource.TestCheckResourceAttr(testResourceName, "local_as_num", "65002"),
					resource.TestCheckResourceAttr(testResourceName, "multipath_relax", "false"),
					resource.TestCheckResourceAttr(testResourceName, "peer_route_convergence_timer", "10"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_mode", "DISABLE"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_timer", "150"),
					resource.TestCheckResourceAttr(testResourceName, "graceful_restart_stale_route_timer", "700"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerBgpExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		sessionContext := testAccGetSessionContext()

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteControllerBgp resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy RouteControllerBgp resource ID not set in resources")
		}

		client := cliRouteControllerBgpClient(sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err := client.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Policy RouteControllerBgp with ID %s does not exist: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyRouteControllerBgpCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	sessionContext := testAccGetSessionContext()

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller_bgp" {
			continue
		}

		resourceID := rs.Primary.ID
		client := cliRouteControllerBgpClient(sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err := client.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy RouteControllerBgp %s still exists", resourceID)
		}
	}

	return nil
}

func testAccNsxtPolicyRouteControllerBgpResourceTemplate(isCreate bool) string {
	ecmp := "true"
	localAs := "65001"
	multipathRelax := "true"
	peerTimer := "5"
	mode := "HELPER_ONLY"
	restartTimer := "120"
	staleTimer := "600"

	if !isCreate {
		ecmp = "false"
		localAs = "65002"
		multipathRelax = "false"
		peerTimer = "10"
		mode = "DISABLE"
		restartTimer = "150"
		staleTimer = "700"
	}

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
  ecmp                               = %s
  local_as_num                       = "%s"
  multipath_relax                    = %s
  peer_route_convergence_timer       = %s
  graceful_restart_mode              = "%s"
  graceful_restart_timer             = %s
  graceful_restart_stale_route_timer = %s
}
`, name, ecmp, localAs, multipathRelax, peerTimer, mode, restartTimer, staleTimer)
}
