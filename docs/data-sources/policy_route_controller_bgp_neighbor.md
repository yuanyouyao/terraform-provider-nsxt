---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp_neighbor"
description: Policy route controller BGP Neighbor data source.
---

# nsxt_policy_route_controller_bgp_neighbor

This data source provides information about Route Controller BGP Neighbor configuration on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_neighbor" "test" {
  route_controller_id = "rc-1"
  id                  = "neigh-1"
}
```

## Argument Reference

* `id` - (Required) The ID of the BGP neighbor config.
* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - The NSX path of the policy resource.
* `display_name` - Display name of the resource.
* `description` - Description of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `tag` - A list of scope + tag pairs associated with this resource.
* `allow_as_in` - Flag to enable allow_as_in option for BGP neighbor.
* `enabled` - Flag to enable/disable BGP peering.
* `gateway_ips` - Next hop gateway IP addresses to reach non-directly connected BGP peers.
* `graceful_restart_mode` - BGP Graceful Restart Configuration Mode.
* `hold_down_time` - Wait time in seconds before declaring peer dead.
* `keep_alive_time` - Interval in seconds between keep alive messages sent to peer.
* `maximum_hop_limit` - Maximum number of hops allowed to reach BGP neighbor.
* `neighbor_address` - Neighbor IP Address.
* `password` - Password for BGP neighbor authentication.
* `remote_as_num` - ASN of the neighbor in ASPLAIN format.
* `source_addresses` - Source IP addresses for BGP peering.
* `bfd_config` - BFD configuration for failure detection.
  * `enabled` - Flag to enable/disable BFD configuration.
  * `interval` - Time interval between heartbeat packets in milliseconds.
  * `multiple` - Number of times heartbeat packet is missed before BFD declares the neighbor is down.
* `route_filtering` - Enable address families and route filtering in each direction.
  * `address_family` - Address family type.
  * `enabled` - Flag to enable/disable address family.
  * `in_route_filters` - Prefix-list or route map paths for IN direction.
  * `out_route_filters` - Prefix-list or route map paths for OUT direction.
  * `maximum_routes` - Maximum number of routes for the address family.
