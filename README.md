# coredns-auto-ipv6-ptr

Goal: Generate IPv6 PTR records on the fly.

Additional benefit: Works with known hosts.

## Examples

### Generate PTR records if not found in a zonefile

```
0.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    autoipv6ptr {
        suffix lan.mydomain.tld
    }
    file your.reverse.zone
    log
}
```

### Same as above but with a transferred zone

```
1.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    autoipv6ptr {
        suffix servers.mydomain.tld
        ttl 60
    }
    secondary {
        transfer from your.master.dns
    }
    log
}
```

## Order is everything!

It's necessary that `file` or `seconary` comes right after `autoipv6ptr`! This plugin always calls the next plugin and checks its return. It will only generate a PTR if a negative result comes back.

## Building a ready-to-use coredns binary using Docker

Using the docker infrastructure it's easy for you to build a working binary with the plugin:

> docker build --pull --no-cache --output type=local,dest=result -f Dockerfile.build .

If everything checks out you'll find an x86_64 binary locally under `result/coredns`.

## Testing

Run:

> ./result/coredns -conf tests/Corefile -p 1337

Test with dig. Known record:

> dig @::1 -p 1337 +short -x 2001:db8:300::10   
>success.example.com.

Unknown record:

> dig @::1 -p 1337 +short -x 2001:db8:300::11   
> 20010db8030000000000000000000011.example.com.