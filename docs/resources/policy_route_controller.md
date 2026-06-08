---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_route_controller"
description: A resource to configure a Route Controller.
---

# nsxt_policy_route_controller

This resource provides a method for the management of a Route Controller.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "test-rc"
  description                            = "Terraform provisioned Route Controller"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ha_mode` - (Optional) High-availability mode for route controller. Currently only `ACTIVE_STANDBY` is supported. Default is `ACTIVE_STANDBY`.
* `virtual_network_appliance_cluster_path` - (Required) Policy path for virtual network appliance cluster.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_route_controller.test PATH
```

The above command imports Route Controller named `test` with the policy path `PATH`.
