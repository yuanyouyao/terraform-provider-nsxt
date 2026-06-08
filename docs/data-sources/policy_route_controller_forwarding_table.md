---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_forwarding_table"
description: Policy route controller forwarding table data source.
---

# nsxt_policy_route_controller_forwarding_table

This data source provides the forwarding table and downloaded CSV forwarding table for a Route Controller on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_forwarding_table" "test" {
  route_controller_id = nsxt_policy_route_controller.test.nsx_id
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the Route Controller. Exactly one of `route_controller_id` or `route_controller_path` must be specified.
* `route_controller_path` - (Optional) The policy path of the Route Controller. Exactly one of `route_controller_id` or `route_controller_path` must be specified.
* `network_prefix` - (Optional) Network address filter parameter.
* `route_source` - (Optional) Filter routes based on the source from which route is learned.
* `virtual_network_appliance_path` - (Optional) Contains string path of virtual network appliance.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `forwarding_table` - The forwarding table entries from the route controller.
    * `count` - Number of route entries.
    * `error_message` - Error message if any.
    * `status` - Status of the forwarding table retrieval.
    * `transport_node_path` - Contains string path of transport node.
    * `virtual_network_appliance_path` - Contains string path of virtual network appliance.
    * `route_entries` - List of forwarding entries.
        * `admin_distance` - Admin distance.
        * `black_hole` - Value of this field will be true if given routes are null routes.
        * `lr_component_id` - Logical router component ID.
        * `lr_component_type` - Logical router component type.
        * `network` - Network CIDR.
        * `next_hop` - Next hop address.
        * `next_hop_gateway` - Next hop gateway path.
        * `route_type` - Route type in forwarding table.
* `forwarding_table_csv` - The forwarding table in CSV format downloaded from the route controller.
