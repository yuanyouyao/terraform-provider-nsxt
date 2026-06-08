---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_route_controller_interface"
description: A resource to configure a Route Controller Interface.
---

# nsxt_policy_route_controller_interface

This resource provides a method for the management of a Route Controller Interface configuration.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `route_controller_path` - (Required) The policy path of the route controller.
* `interface_address` - (Required) Route Controller Interface Address parameters.
  * `portgroup_id` - (Required) DV port group identifier discovered from vCenter.
  * `virtual_network_appliance_path` - (Required) Policy path for virtual network appliance.
  * `interface_subnet` - (Required) IP address and subnet specification for interface.
    * `ip_addresses` - (Required) IP addresses assigned to interface.
    * `prefix_len` - (Required) Subnet prefix length.
* `display_name` - (Optional) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `mtu` - (Optional) MTU size.
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode. Supported values: `NONE`, `STRICT`. Default is `NONE`.
* `floating_ip_subnets` - (Optional) IP address and subnet specification for VIP floating IP address subnets.
  * `ip_addresses` - (Required) IP addresses assigned to interface.
  * `prefix_len` - (Required) Subnet prefix length.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource (Route Controller ID + "/" + Interface ID).
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_route_controller_interface.test ROUTE_CONTROLLER_ID/INTERFACE_ID
```

The above command imports Route Controller Interface named `test` with Route Controller ID and Interface ID.
