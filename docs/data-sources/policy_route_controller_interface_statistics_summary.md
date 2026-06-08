---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_interface_statistics_summary"
description: Policy route controller Interface statistics summary data source.
---

# nsxt_policy_route_controller_interface_statistics_summary

This data source provides the statistics summary of a specific interface configured on a Route Controller in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_interface_statistics_summary" "test" {
  route_controller_id = "rc-1"
  interface_id        = "iface-1"
}
```

## Argument Reference

* `interface_id` - (Required) The ID of the route controller interface.
* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `interface_path` - Absolute path of the interface.
* `last_update_timestamp` - Timestamp of the last update.
* `rx` - List of RX counters:
  * `total_packets` - Total packets received.
  * `total_bytes` - Total bytes received.
  * `dropped_packets` - Dropped packets.
  * `blocked_packets` - Blocked packets.
  * `firewall_dropped_packets` - Firewall dropped packets.
  * `ipv6_dropped_packets` - IPv6 dropped packets.
  * `no_arp_dropped_packets` - No ARP dropped packets.
  * `no_route_dropped_packets` - No route dropped packets.
  * `rpf_check_dropped_packets` - RPF check dropped packets.
  * `ttl_exceeded_dropped_packets` - TTL exceeded dropped packets.
* `tx` - List of TX counters:
  * `total_packets` - Total packets transmitted.
  * `total_bytes` - Total bytes transmitted.
  * `dropped_packets` - Dropped packets.
  * `blocked_packets` - Blocked packets.
  * `firewall_dropped_packets` - Firewall dropped packets.
  * `ipv6_dropped_packets` - IPv6 dropped packets.
  * `no_arp_dropped_packets` - No ARP dropped packets.
  * `no_route_dropped_packets` - No route dropped packets.
  * `rpf_check_dropped_packets` - RPF check dropped packets.
  * `ttl_exceeded_dropped_packets` - TTL exceeded dropped packets.
