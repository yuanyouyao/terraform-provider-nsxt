// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestPolicyRouteControllerCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"ha_mode":      "ACTIVE_STANDBY",
}

var accTestPolicyRouteControllerUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"ha_mode":      "ACTIVE_STANDBY",
}

func testAccNsxtPolicyRouteControllerPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
}

func TestAccResourceNsxtPolicyRouteController_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"
	testDataSourceName := "data.nsxt_policy_route_controller.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerCheckDestroy(state, accTestPolicyRouteControllerUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRouteControllerCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", accTestPolicyRouteControllerCreateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "ha_mode", accTestPolicyRouteControllerCreateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testDataSourceName, "virtual_network_appliance_cluster_path"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRouteControllerUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", accTestPolicyRouteControllerUpdateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "ha_mode", accTestPolicyRouteControllerUpdateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testDataSourceName, "virtual_network_appliance_cluster_path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyRouteController_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerCheckDestroy(state, accTestPolicyRouteControllerUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerTemplate(false),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyRouteControllerExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteController resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy RouteController resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyRouteControllerExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy RouteController with ID %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyRouteControllerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller" {
			continue
		}

		resourceID := rs.Primary.ID
		exists, err := resourceNsxtPolicyRouteControllerExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy RouteController %s still exists", resourceID)
		}
	}

	return nil
}

func testAccNsxtPolicyRouteControllerTemplate(isCreate bool) string {
	attrs := accTestPolicyRouteControllerCreateAttributes
	if !isCreate {
		attrs = accTestPolicyRouteControllerUpdateAttributes
	}

	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "%s"
  ha_mode                                = "%s"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller" "test" {
  id = nsxt_policy_route_controller.test.nsx_id
}
`, attrs["display_name"], attrs["description"], attrs["ha_mode"])
}

func testAccNsxtPolicyRouteControllerVnaTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_virtual_network_appliance_cluster" "vna" {
  display_name = "%s"
}
`, os.Getenv("NSXT_TEST_RC_VNA_CLUSTER_NAME"))
}
