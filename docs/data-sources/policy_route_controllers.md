---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controllers"
description: Policy route controllers data source.
---

# nsxt_policy_route_controllers

This data source provides information about multiple Route Controllers on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controllers" "all" {}

data "nsxt_policy_route_controllers" "filtered" {
  display_name = ".*-rc"
}
```

## Argument Reference

* `display_name` - (Optional) Display name of Route Controller. Supports regular expressions.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - Mapping of Route Controller instance ID by display name.
