# desec-dyndns-client

This is a basic DynDNS Client for [deSEC.io](https://desec.io/) with Dual Stack support. This is very rudimentary, offers no special functions and may break in the future!

## Getting Started

You need to provide Ipv4 and Ipv6 network to the container, because it needs to check both address types. 

The following Podman command shows, how to run this in a container:

```sh
podman run --network=host -e DYNDNS_DOMAIN=host.example.com -e DYNDNS_TOKEN=YOUR_TOKEN git.maltech.io/maltech/desec-dyndns-client/desec-dyndns-client:1.0.0
```
