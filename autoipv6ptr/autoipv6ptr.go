package autoipv6ptr

import (
	"context"
    "strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type AutoIPv6PTR struct {
	Next plugin.Handler

	// Presets are static entries which should not be generated
	Presets map[string]string

	Suffix string
}

// ServeDNS implements the plugin.Handler interface.
func (plugin AutoIPv6PTR) ServeDNS(ctx context.Context, writer dns.ResponseWriter, request *dns.Msg) (int, error) {
	if request.Question[0].Qtype != dns.TypePTR {
		return plugin.NextOrFailure(plugin.Name(), plugin.Next, ctx, writer, request)
	}

	var responsePtrValue string

	if ptrValue, found := plugin.Presets[request.Question[0].Name]; found {
		responsePtrValue = ptrValue
	} else {
		responsePtrValue = request.Question[0].Name
	    responsePtrValue = RemoveIP6DotArpa(responsePtrValue)
	    responsePtrValue = RemoveDots(responsePtrValue)
	    responsePtrValue = ReverseString(responsePtrValue)
	    responsePtrValue += "." + plugin.Suffix + "."
	}

	message := new(dns.Msg)
	message.SetReply(request)
	hdr := dns.RR_Header{Name: request.Question[0].Name, Ttl: 900, Class: dns.ClassINET, Rrtype: dns.TypePTR}
	message.Answer = []dns.RR{&dns.PTR{Hdr: hdr, Ptr: responsePtrValue}}

	writer.WriteMsg(message)
	return 0, nil
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

// Name implements the Handler interface.
func (a AutoIPv6PTR) Name() string { return "plugin" }
