# containerlab

Couple of exercise to get started with containerlab.dev

## Setup

```
docker run --rm -it --privileged --network host --pid host -d \
    --name containerlab \
    -v /var/lib/docker/containers:/var/lib/docker/containers \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/run/netns:/var/run/netns \
    -v "$PWD":"$PWD" \
    -w "$PWD" \
    ghcr.io/srl-labs/clab \
    containerlab deploy -t l2l3.clab.yml

docker exec -it clab-l2l3-sw sr_cli      # the switch
docker exec -it clab-l2l3-h1 sh          # host 1
docker exec -it clab-l2l3-h2 sh          # host 2
```

## Exercise 1: L2, the switch as a “smart cable”

On the switch, put both ports into one bridge (SR Linux calls a bridge a mac-vrf):

```
show interface all

enter candidate
set / interface ethernet-1/1 admin-state enable
set / interface ethernet-1/1 subinterface 0 type bridged
set / interface ethernet-1/2 admin-state enable
set / interface ethernet-1/2 subinterface 0 type bridged
set / network-instance lan type mac-vrf
set / network-instance lan interface ethernet-1/1.0
set / network-instance lan interface ethernet-1/2.0
commit now
```

Give the hosts addresses in the same subnet:

```
# h1
ip addr add 10.0.0.1/24 dev eth1
# h2
ip addr add 10.0.0.2/24 dev eth1
# h1
ping -c3 10.0.0.2
ip neigh
```

What to notice:

- The ping works, even though the switch has no IP address at all.
- `ip neigh` on h1 shows h2’s own MAC address. h1 talks to h2 directly, and the switch is invisible.
- On the switch, `show network-instance lan bridge-table mac-table all` lists both hosts’ MACs. That table is all an L2 switch knows: “this MAC is behind this port.”

# Exercise 2: L2 can’t cross subnets

Move the hosts into different subnets:

```
# h1
ip addr flush dev eth1; ip addr add 10.0.1.1/24 dev eth1
# h2
ip addr flush dev eth1; ip addr add 10.0.2.2/24 dev eth1
# h1
ping -c3 10.0.2.2       # → Network unreachable
```

The cable still works, but h1 has no idea how to reach 10.0.2.0/24. Moving between IP networks needs a router. That’s what L3 is for.

## Exercise 3: L3, the switch as a router

Turn the ports into routed ports, each with its own IP address:

```
enter candidate
delete / network-instance lan
delete / interface ethernet-1/1 subinterface 0
delete / interface ethernet-1/2 subinterface 0
set / interface ethernet-1/1 subinterface 0 ipv4 admin-state enable
set / interface ethernet-1/1 subinterface 0 ipv4 address 10.0.1.254/24
set / interface ethernet-1/2 subinterface 0 ipv4 admin-state enable
set / interface ethernet-1/2 subinterface 0 ipv4 address 10.0.2.254/24
set / network-instance default interface ethernet-1/1.0
set / network-instance default interface ethernet-1/2.0
commit now
```

Configure routes for the other host’s network via the switch:

```
# h1
ip route add 10.0.2.0/24 via 10.0.1.254
# h2
ip route add 10.0.1.0/24 via 10.0.2.254
# h1
ping -c3 10.0.2.2
traceroute 10.0.2.2
```

What to notice:

- The ping works across two different networks.
- The reply shows `ttl=63` instead of 64: one router hop decremented it. traceroute shows that hop as 10.0.1.254.
- `ip neigh` on h1 now shows only the switch’s MAC, not h2’s. h1 hands every packet to its gateway and never sees h2 at L2.

On the switch, `show network-instance default ipv4 route` lists the two directly connected networks (route type local). An L3 device forwards each packet by looking up its destination IP in this table.

## Cleanup

```
docker rm -f containerlab clab-l2l3-sw clab-l2l3-h1 clab-l2l3-h2
docker network rm clab
rm -rf clab-l2l3
```