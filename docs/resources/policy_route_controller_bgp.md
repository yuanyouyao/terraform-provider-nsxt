---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_route_controller_bgp"
description: A resource to configure a Route Controller BGP routing config.
---

# nsxt_policy_route_controller_bgp

This resource provides a method for the management of a Route Controller BGP routing configuration.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller_bgp" "test" {
  route_controller_path              = nsxt_policy_route_controller.test.path
  ecmp                               = true
  local_as_num                       = "65001"
  multipath_relax                    = true
  peer_route_convergence_timer       = 5
  graceful_restart_mode              = "HELPER_ONLY"
  graceful_restart_timer             = 120
  graceful_restart_stale_route_timer = 600
}
```

## Argument Reference

The following arguments are supported:

* `route_controller_path` - (Required) The policy path of the route controller.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `ecmp` - (Optional) Flag to enable ECMP. Default is `true`.
* `local_as_num` - (Optional) BGP AS number in ASPLAIN/ASDOT format.
* `multipath_relax` - (Optional) Flag to enable BGP multipath relax option.
* `peer_route_convergence_timer` - (Optional) Extra time in seconds the router must wait before sending the UP notification.
* `graceful_restart_mode` - (Optional) BGP Graceful Restart Configuration Mode. Supported values: `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`. Default is `HELPER_ONLY`.
* `graceful_restart_timer` - (Optional) Maximum time taken (in seconds) for a BGP session to be established after a restart. Default is `120`.
* `graceful_restart_stale_route_timer` - (Optional) Maximum time (in seconds) before stale routes are removed from the RIB when BGP restarts. Default is `600`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource (Route Controller ID).
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_route_controller_bgp.test PATH
```

The above command imports Route Controller BGP routing config named `test` with the policy path `PATH` or Route Controller ID.
