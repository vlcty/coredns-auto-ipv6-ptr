# coredns-auto-ipv6-ptr

Some services require that a RDNS request resolves to a PTR record. With this CoreDNS plugin, you can generate these PTR records on the fly based on the requested IPv6 address. The plugin translates the requested address and appends a suffix. Additionally, you can create a so-called presets file to answer with a "real" record if a specific request is received.

Examples:

1) `2001:db8:300:b002:5054:ff:fe4b:db44` could be translated to `20010db80300b002505400fffe4bdb44.mydomain.tld`
2) `2001:db8:300:b002:5054:ff:fe4b:db45` could be translated to `myhost.mydomain.tld` via the presets file

## Presets

Each line in a presets file contains a manual override seperated by a comma. The first part contains the request and the second part contains the value which should be returned. For example:

```
5.4.b.d.b.4.e.f.f.f.0.0.4.5.0.5.2.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa.;myhost.mydomain.tld
6.4.b.d.b.4.e.f.f.f.0.0.4.5.0.5.2.0.0.b.0.0.3.0.8.b.d.0.1.0.0.2.ip6.arpa.;firewall.mydomain.tld
```

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
| presetsfile | | The absolute path to the presets file |
| ttl | 900 | The TTL value the answer should have |

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
        presetsfile /var/lib/coredns/presets.servers.mydomain.tld
        ttl 60
    }
}
```

The suffix is a mandatory argument. The presetsfile is optional. The presets file is read on plugin startup.
