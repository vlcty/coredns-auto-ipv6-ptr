# coredns-auto-ipv6-ptr

Some services require that RDNS requests resolve to PTR records. With this CoreDNS plugin, you can generate these PTR records on the fly based on the requested IPv6 address. The plugin translates the requested address and appends a suffix.

## Translation process

Let's say the plugin receives a PTR request for `4.4.b.d.b.4.e.f.f.f.0.0.4.5.0.5.2.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa.`. If there is a preset found the preset value will be used in the anwer. If there is none the regular translation process starts:

1) Strip `.ip6.arpa.`: `4.4.b.d.b.4.e.f.f.f.0.0.4.5.0.5.2.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2`
2) Remove all dots: `44bdb4efff004505200b00308bd01002`
3) Reverse the string: `20010db80300b002505400fffe4bdb44`
4) Append the suffix and return the result: `20010db80300b002505400fffe4bdb44.mydomain.tld`

## Corefile example

Possible plugin arguments:

| Argument | Default value | Description |
|-|-|-|
| suffix | | The suffix to append when regular translating happens |
| ttl | 900 | The TTL value the answer should have in seconds |

Let's say your provider allocated `2001:db8:300:b000::/56` to you. You sliced two subnets out of it:

1) 2001:db8:300:b000::/64 => lan.myhost.tld
2) 2001:db8:300:b001::/64 => servers.myhost.tld

You Corefile would look something like this:

```
0.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    log
    autoipv6ptr {
        suffix lan.mydomain.tld
    }
}

1.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    log
    autoipv6ptr {
        suffix servers.mydomain.tld
        ttl 60
    }
}
```

## Working with known hosts

If you have knows hosts you want to return specific PTR records you can do this via the `file` or `secondary` plugin. However there is a catch to this! `file` and `secondary` are so called backends which return NXDOMAIN when no record was found. You can find a patch provided by GitHub user @dorchain in the file `file-fallthrough.patch`. This little patch makes `file` and `secondary` falling through if no record was found. Apply it via git patch from your CoreDNS root directory:

> git apply plugin/autoipv6ptr/file-fallthrough.patch

And build CoreDNS to your needs. Sample Corefile:

```
0.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    log
    file your.reverse.zone
    autoipv6ptr {
        suffix lan.mydomain.tld
    }
}

1.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa {
    log
    secondary {
        transfer from your.master.dns
    }
    autoipv6ptr {
        suffix servers.mydomain.tld
        ttl 60
    }
}
```

## Building a ready-to-use coredns binary using Docker

Using the docker infrastructure it's easy for you to build a working binary with the plugin:

> docker build --pull --no-cache --output type=local,dest=result -f Dockerfile.build .

If everything checks out you'll find an x86_64 binary locally under `result/coredns`.
