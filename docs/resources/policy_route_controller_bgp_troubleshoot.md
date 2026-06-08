---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_route_controller_bgp_troubleshoot"
description: Policy route controller BGP Troubleshoot resource.
---

# nsxt_policy_route_controller_bgp_troubleshoot

This resource provides a way to configure BGP Troubleshoot settings on a Route Controller in NSX.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller_bgp_troubleshoot" "test" {
  route_controller_path       = nsxt_policy_route_controller.test.path
  bfd_control_pkt_diagnostics = true
  bgp_session_diagnostics     = true
  system_diagnostics          = false
}
```

## Argument Reference

* `route_controller_path` - (Required) The policy path of the route controller.
* `bfd_control_pkt_diagnostics` - (Optional) Flag to enable/disable the collection of the timestamps of sent and received BFD control messages per BFD peer session. Defaults to `true`.
* `bgp_session_diagnostics` - (Optional) Flag to enable/disable the collection of the timestamps of sent and received Keep-Alive messages per BGP peer session, and the session states. Defaults to `true`.
* `system_diagnostics` - (Optional) Flag to enable/disable the collection of system diagnostic data such as ARP, Ping, CPU stats, etc. Defaults to `false`.
* `tag` - (Optional) Opaque identifiers meaningful to the API user.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `path` - Path of the resource on NSX Policy Manager.
* `revision` - The _revision property describes the current revision of the resource.

## Import

An existing BGP Troubleshoot configuration can be imported using the Route Controller ID or path:

```bash
terraform import nsxt_policy_route_controller_bgp_troubleshoot.test <route-controller-id>
```
