---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp_troubleshoot"
description: Policy route controller BGP Troubleshoot data source.
---

# nsxt_policy_route_controller_bgp_troubleshoot

This data source provides information about the Route Controller BGP Troubleshoot configuration in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_troubleshoot" "test" {
  route_controller_id = "rc-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - Path of the resource on NSX Policy Manager.
* `revision` - The _revision property describes the current revision of the resource.
* `bfd_control_pkt_diagnostics` - Flag to enable/disable the collection of the timestamps of sent and received BFD control messages per BFD peer session.
* `bgp_session_diagnostics` - Flag to enable/disable the collection of the timestamps of sent and received Keep-Alive messages per BGP peer session, and the session states.
* `system_diagnostics` - Flag to enable/disable the collection of system diagnostic data such as ARP, Ping, CPU stats, etc.
