---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp"
description: Policy route controller BGP routing config data source.
---

# nsxt_policy_route_controller_bgp

This data source provides information about Route Controller BGP routing configuration on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp" "test" {
  route_controller_id = "rc-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - The NSX path of the policy resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `tag` - A list of scope + tag pairs associated with this resource.
* `ecmp` - Flag to enable ECMP.
* `local_as_num` - BGP AS number in ASPLAIN/ASDOT format.
* `multipath_relax` - Flag to enable BGP multipath relax option.
* `peer_route_convergence_timer` - Extra time in seconds the router must wait before sending the UP notification.
* `graceful_restart_mode` - BGP Graceful Restart Configuration Mode.
* `graceful_restart_timer` - Maximum time taken (in seconds) for a BGP session to be established after a restart.
* `graceful_restart_stale_route_timer` - Maximum time (in seconds) before stale routes are removed from the RIB when BGP restarts.
