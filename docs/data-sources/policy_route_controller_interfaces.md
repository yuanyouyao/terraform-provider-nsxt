---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_interfaces"
description: Policy route controller Interfaces data source.
---

# nsxt_policy_route_controller_interfaces

This data source provides the list of interfaces configured on a Route Controller in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_interfaces" "test" {
  route_controller_id = "rc-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `results` - List of route controller interfaces. Each item contains the following attributes:
  * `id` - Unique identifier of this resource.
  * `display_name` - Display name of this resource.
  * `description` - Description of this resource.
  * `path` - Absolute path of this object.
  * `revision` - Current revision of the resource.
  * `mtu` - MTU size.
  * `urpf_mode` - Unicast Reverse Path Forwarding mode.
  * `floating_ip_subnets` - List of IP address and subnet specifications for VIP floating IP address subnets:
    * `ip_addresses` - IP addresses assigned to interface.
    * `prefix_len` - Subnet prefix length.
  * `interface_address` - List of Route Controller Interface Address object parameters:
    * `portgroup_id` - DV port group identifier discovered from vCenter.
    * `virtual_network_appliance_path` - Policy path for virtual network appliance.
    * `interface_subnet` - List of IP address and subnet specifications for interface:
      * `ip_addresses` - IP addresses assigned to interface.
      * `prefix_len` - Subnet prefix length.
