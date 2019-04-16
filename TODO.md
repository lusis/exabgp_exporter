# TODO items and notes

## Add flows to exported stats

We capture the data now in the global announcement/withdraw tracking.
It just needs to have an Gauge added for it

## ~~Consider reworking how we track routes being available~~

After reading up on the bgp spec more trying to recall years old experiences, I realized that tracking nlri by next-hop is irrelevant.
The behaviour of routers on handling next hop is dependent on ibgp vs ebgp.
There's a different behaviour for each but regardless, getting an update for the same nlri either changes the nexthop to self or replaces the existing entry.

What this means is honestly we don't even CARE about what the next hop is for state exporting. We just need to know if we've announced a route for a given network.