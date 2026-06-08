---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp_neighbors"
description: Policy route controller BGP Neighbors data source.
---

# nsxt_policy_route_controller_bgp_neighbors

This data source provides information about all Route Controller BGP Neighbors configured on a Route Controller in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_neighbors" "test" {
  route_controller_id = "rc-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `display_name` - (Optional) Display name of BGP Neighbor. Supports regular expressions.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - Mapping of BGP Neighbor instance ID by display name.
