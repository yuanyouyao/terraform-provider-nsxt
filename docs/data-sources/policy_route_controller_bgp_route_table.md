# nsxt_policy_route_controller_bgp_route_table

This data source provides the BGP routing table entries from the specified virtual network appliance node for the given route controller.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_route_table" "test" {
  route_controller_id            = "rc-id"
  virtual_network_appliance_path = "/infra/virtual-network-appliances/vna-1"
  network_prefix                 = "10.0.0.0/24"
}
```

## Argument Reference

* `route_controller_id` - (Optional) The ID of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `route_controller_path` - (Optional) The policy path of the route controller. Exactly one of `route_controller_id` or `route_controller_path` must be set.
* `virtual_network_appliance_path` - (Required) Policy path of virtual network appliance to retrieve BGP routes from.
* `network_prefix` - (Optional) Network address filter parameter.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `results` - The list of BGP route table entries.
  * `enforcement_point_path` - Enforcement point path.
  * `gateway_path` - Policy path of the gateway.
  * `last_update_timestamp` - Timestamp of the last update.
  * `transport_node_id` - ID of the transport node.
  * `transport_node_path` - Policy path of the transport node.
  * `route_details` - List of route details.
    * `as_path` - BGP AS path.
    * `bestpath` - Whether the route is the best path.
    * `community` - BGP community string.
    * `esi` - Ethernet Segment Identifier.
    * `eth_tag` - Ethernet Tag.
    * `evpn_route_type` - EVPN route type.
    * `extended_community` - BGP extended community string.
    * `large_community` - BGP large community string.
    * `local_pref` - Local preference value.
    * `med` - Multi-Exit Discriminator value.
    * `multipath` - Whether the route is a multipath route.
    * `network` - Network address CIDR.
    * `path_from` - Source of the path.
    * `peer_id` - BGP peer ID.
    * `rd` - Route Distinguisher.
    * `rmac` - Router MAC address.
    * `route_origin` - BGP route origin.
    * `stale` - Whether the route is stale.
    * `valid` - Whether the route is valid.
    * `vni` - VXLAN Network Identifier.
    * `weight` - Route weight.
    * `nexthops` - List of next hops.
      * `ip_address` - IP address of the next hop.
      * `scope` - Scope of the next hop.
