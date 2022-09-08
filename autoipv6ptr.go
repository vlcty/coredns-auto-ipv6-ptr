package autoipv6ptr

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
)

const AUTOIPV6PTR_PLUGIN_NAME string = "autoipv6ptr"

type AutoIPv6PTR struct {
	Next plugin.Handler

	TTL uint32

	Suffix string
}

// ServeDNS implements the plugin.Handler interface.
func (v6ptr AutoIPv6PTR) ServeDNS(ctx context.Context, writer dns.ResponseWriter, request *dns.Msg) (int, error) {
	if request.Question[0].Qtype != dns.TypePTR {
		return plugin.NextOrFailure(v6ptr.Name(), v6ptr.Next, ctx, writer, request)
	}

	nw := nonwriter.New(writer)

	plugin.NextOrFailure(v6ptr.Name(), v6ptr.Next, ctx, nw, request)

	if len(nw.Msg.Answer) > 0 {
		writer.WriteMsg(nw.Msg)
	} else {
		responsePtrValue := request.Question[0].Name
		responsePtrValue = RemoveIP6DotArpa(responsePtrValue)
		responsePtrValue = RemoveDots(responsePtrValue)
		responsePtrValue = ReverseString(responsePtrValue)
		responsePtrValue += "." + v6ptr.Suffix + "."

		message := new(dns.Msg)
		message.SetReply(request)
		message.Authoritative = true
		message.Rcode = dns.RcodeSuccess
		hdr := dns.RR_Header{Name: request.Question[0].Name, Ttl: v6ptr.TTL, Class: dns.ClassINET, Rrtype: dns.TypePTR}
		message.Answer = []dns.RR{&dns.PTR{Hdr: hdr, Ptr: responsePtrValue}}

		writer.WriteMsg(message)
	}

	return dns.RcodeSuccess, nil
}

func RemoveIP6DotArpa(input string) string {
	return strings.ReplaceAll(input, ".ip6.arpa.", "")
}

func RemoveDots(input string) string {
	return strings.ReplaceAll(input, ".", "")
}

func ReverseString(input string) string {
	// Copied from https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (a AutoIPv6PTR) Name() string { return AUTOIPV6PTR_PLUGIN_NAME }

func init() {
	plugin.Register(AUTOIPV6PTR_PLUGIN_NAME, setup)
}

func setup(c *caddy.Controller) error {
	v6ptr := AutoIPv6PTR{}
	v6ptr.TTL = 900

	for c.Next() {
		switch c.Val() {
		case "suffix":
			suffix := c.RemainingArgs()[0]

			if len(suffix) == 0 {
				return plugin.Error(AUTOIPV6PTR_PLUGIN_NAME, errors.New("Suffix can't be empty"))
			} else {
				v6ptr.Suffix = suffix
			}

		case "ttl":
			possibleTTL := c.RemainingArgs()[0]
			ttl, err := strconv.ParseUint(possibleTTL, 10, 32)

			if err != nil {
				return plugin.Error(AUTOIPV6PTR_PLUGIN_NAME, err)
			} else {
				v6ptr.TTL = uint32(ttl)
			}

		default:
			continue
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		v6ptr.Next = next
		return v6ptr
	})

	return nil
}
