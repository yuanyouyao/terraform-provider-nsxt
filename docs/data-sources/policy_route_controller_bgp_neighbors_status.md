---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp_neighbors_status"
description: Policy route controller BGP Neighbors Status data source.
---

# nsxt_policy_route_controller_bgp_neighbors_status

This data source provides the runtime status of BGP Neighbors configured on a Route Controller in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_neighbors_status" "test" {
  route_controller_id = "rc-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `bgp_neighbor_type` - (Optional) BGP neighbor type: `INTER_SR` or `USER`.
* `enforcement_point_path` - (Optional) Enforcement point path.
* `stats_type` - (Optional) Stats type.
* `transport_node_id` - (Optional) Transport node ID.
* `virtual_network_appliance_path` - (Optional) Virtual network appliance path.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `results` - List of BGP neighbor statuses. Each status contains the following attributes:
  * `neighbor_address` - The IP of the BGP neighbor.
  * `connection_state` - Current state of the BGP session.
  * `connection_drop_count` - Count of connection drops.
  * `established_connection_count` - Count of connections established.
  * `graceful_restart_mode` - Current state of graceful restart of BGP neighbor.
  * `hold_time` - Hold time.
  * `is_dynamic` - Whether the neighbor is dynamic.
  * `keep_alive_interval` - Keep alive interval.
  * `local_port` - Local TCP port.
  * `messages_received` - Count of messages received from the neighbor.
  * `messages_sent` - Count of messages sent to the neighbor.
  * `neighbor_edge_node` - Inter-SR neighbor edge node policy path.
  * `neighbor_path` - Policy intent path of dynamic bgp neighbor.
  * `neighbor_router_id` - Router ID of the BGP neighbor.
  * `remote_as_number` - AS number of the BGP neighbor.
  * `remote_port` - Remote TCP port.
  * `remote_site_path` - Remote site path.
  * `route_controller_path` - Route controller path.
  * `source_address` - Source IP address.
  * `time_since_established` - Time (in seconds) since connection was established.
  * `total_in_prefix_count` - Sum of in prefixes counts.
  * `total_out_prefix_count` - Sum of out prefixes counts.
  * `type` - BGP neighbor type.
  * `virtual_network_appliance_path` - Virtual network appliance path.
  * `announced_capabilities` - BGP capabilities sent to BGP neighbor.
  * `negotiated_capabilities` - BGP capabilities negotiated with BGP neighbor.
  * `address_families` - List of address families of BGP neighbor:
    * `type` - BGP address family type.
    * `in_prefix_count` - Count of in prefixes.
    * `out_prefix_count` - Count of out prefixes.
