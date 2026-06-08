---
subcategory: "EVPN"
page_title: "NSXT: policy_route_controller_bgp_neighbor_advertised_routes"
description: Policy route controller BGP Neighbor Advertised Routes data source.
---

# nsxt_policy_route_controller_bgp_neighbor_advertised_routes

This data source provides the routes advertised by a BGP Neighbor configured on a Route Controller in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_neighbor_advertised_routes" "test" {
  route_controller_id = "rc-1"
  neighbor_id         = "neigh-1"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `neighbor_id` - (Optional) The ID of the BGP neighbor. Exactly one of `neighbor_id` or `neighbor_path` must be set.
* `neighbor_path` - (Optional) The policy path of the BGP neighbor. Exactly one of `neighbor_id` or `neighbor_path` must be set.
* `enforcement_point_path` - (Optional) Enforcement point path.
* `neighbor_address` - (Optional) Dynamically discovered BGP neighbor address.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `results` - List of advertised routes. Each item contains the following attributes:
  * `enforcement_point_path` - Enforcement point path.
  * `neighbor_path` - BGP neighbor policy path.
  * `virtual_network_appliance_routes` - List of routes per virtual network appliance (transport node):
    * `source_address` - BGP neighbor source IP address.
    * `transport_node_id` - Transport node ID.
    * `routes` - List of route details:
      * `as_path` - BGP AS path attribute.
      * `esi` - Ethernet Segment Identifier for EVPN routes.
      * `eth_tag` - Ethernet Tag ID for EVPN routes.
      * `evpn_route_type` - EVPN route type for EVPN routes.
      * `local_pref` - BGP Local Preference attribute.
      * `med` - BGP Multi Exit Discriminator attribute.
      * `network` - CIDR network address.
      * `next_hop` - Next hop IP address.
      * `rd` - BGP Route Distinguisher attribute.
      * `rmac` - Router MAC address for EVPN routes.
      * `rmac_len` - Router MAC address length for EVPN routes.
      * `weight` - BGP Weight attribute.
