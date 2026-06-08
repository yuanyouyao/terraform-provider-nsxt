---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_interface"
description: Policy route controller Interface data source.
---

# nsxt_policy_route_controller_interface

This data source provides information about Route Controller Interface configuration on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_interface" "test" {
  route_controller_id = "rc-1"
  id                  = "iface-1"
}
```

## Argument Reference

* `id` - (Required) The ID of the route controller interface.
* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - The NSX path of the policy resource.
* `display_name` - Display name of the resource.
* `description` - Description of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `tag` - A list of scope + tag pairs associated with this resource.
* `mtu` - MTU size.
* `urpf_mode` - Unicast Reverse Path Forwarding mode.
* `floating_ip_subnets` - IP address and subnet specification for VIP floating IP address subnets.
  * `ip_addresses` - IP addresses assigned to interface.
  * `prefix_len` - Subnet prefix length.
* `interface_address` - Route Controller Interface Address parameters.
  * `portgroup_id` - DV port group identifier discovered from vCenter.
  * `virtual_network_appliance_path` - Policy path for virtual network appliance.
  * `interface_subnet` - IP address and subnet specification for interface.
    * `ip_addresses` - IP addresses assigned to interface.
    * `prefix_len` - Subnet prefix length.
