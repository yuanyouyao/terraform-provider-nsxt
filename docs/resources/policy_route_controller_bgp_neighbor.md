---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_route_controller_bgp_neighbor"
description: A resource to configure a Route Controller BGP Neighbor.
---

# nsxt_policy_route_controller_bgp_neighbor

This resource provides a method for the management of a Route Controller BGP Neighbor configuration.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  route_controller_path = nsxt_policy_route_controller.test.path
  display_name          = "neigh-test"
  neighbor_address      = "192.168.1.1"
  remote_as_num         = "65001"
  allow_as_in           = true
  enabled               = true
  hold_down_time        = 180
  keep_alive_time       = 60
  maximum_hop_limit     = 1

  bfd_config {
    enabled  = true
    interval = 500
    multiple = 3
  }

  route_filtering {
    address_family    = "IPV4"
    enabled           = true
    in_route_filters  = ["/infra/prefix-lists/in-filter"]
    out_route_filters = ["/infra/prefix-lists/out-filter"]
    maximum_routes    = 1000
  }
}
```

## Argument Reference

The following arguments are supported:

* `route_controller_path` - (Required) The policy path of the route controller.
* `neighbor_address` - (Required) Neighbor IP Address.
* `remote_as_num` - (Required) ASN of the neighbor in ASPLAIN format.
* `display_name` - (Optional) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `allow_as_in` - (Optional) Flag to enable allow_as_in option for BGP neighbor. Default is `false`.
* `enabled` - (Optional) Flag to enable/disable BGP peering. Default is `true`.
* `gateway_ips` - (Optional) Next hop gateway IP addresses to reach non-directly connected BGP peers.
* `graceful_restart_mode` - (Optional) BGP Graceful Restart Configuration Mode. Supported values: `DISABLE`, `HELPER_ONLY`. Default is `HELPER_ONLY`.
* `hold_down_time` - (Optional) Wait time in seconds before declaring peer dead. Default is `180`.
* `keep_alive_time` - (Optional) Interval in seconds between keep alive messages sent to peer. Default is `60`.
* `maximum_hop_limit` - (Optional) Maximum number of hops allowed to reach BGP neighbor. Default is `1`.
* `password` - (Optional) Password for BGP neighbor authentication.
* `source_addresses` - (Optional) Source IP addresses for BGP peering.
* `bfd_config` - (Optional) BFD configuration for failure detection.
  * `enabled` - (Optional) Flag to enable/disable BFD configuration. Default is `false`.
  * `interval` - (Optional) Time interval between heartbeat packets in milliseconds. Default is `500`.
  * `multiple` - (Optional) Number of times heartbeat packet is missed before BFD declares the neighbor is down. Default is `3`.
* `route_filtering` - (Optional) Enable address families and route filtering in each direction.
  * `address_family` - (Required) Address family type. Supported values: `IPV4`, `IPV6`, `L2VPN_EVPN`.
  * `enabled` - (Optional) Flag to enable/disable address family. Default is `true`.
  * `in_route_filters` - (Optional) Prefix-list or route map paths for IN direction.
  * `out_route_filters` - (Optional) Prefix-list or route map paths for OUT direction.
  * `maximum_routes` - (Optional) Maximum number of routes for the address family.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource (Route Controller ID + "/" + Neighbor ID).
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_route_controller_bgp_neighbor.test ROUTE_CONTROLLER_ID/NEIGHBOR_ID
```

The above command imports Route Controller BGP Neighbor named `test` with Route Controller ID and Neighbor ID.
