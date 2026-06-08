---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_state"
description: Policy route controller state data source.
---

# nsxt_policy_route_controller_state

This data source provides the runtime state information for a Route Controller on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_state" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the Route Controller. Exactly one of `route_controller_id` or `route_controller_path` must be specified.
* `route_controller_path` - (Optional) The policy path of the Route Controller. Exactly one of `route_controller_id` or `route_controller_path` must be specified.
* `source` - (Optional) The source for which state is retrieved.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `last_update_timestamp` - Timestamp when the status was last updated.
* `logical_gateway_id` - The ID of the route controller logical gateway.
* `virtual_network_appliance_cluster_path` - Policy path of virtual network appliance cluster.
* `per_node_status` - Per node status of the route controller.
    * `high_availability_status` - High availability status on virtual network appliance node.
    * `node_type` - Type of node.
    * `service_gateway_id` - The ID of the service gateway where the status is retrieved.
    * `virtual_network_appliance_path` - Policy path of virtual network appliance where the node status is retrieved.
