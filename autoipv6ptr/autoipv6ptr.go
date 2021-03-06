package autoipv6ptr

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

const AUTOIPV6PTR_PLUGIN_NAME string = "autoipv6ptr"

type AutoIPv6PTR struct {
	Next plugin.Handler

	// Presets are static entries which should not be generated
	Presets map[string]string
	TTL     uint32

	Suffix string
}

// ServeDNS implements the plugin.Handler interface.
func (v6ptr AutoIPv6PTR) ServeDNS(ctx context.Context, writer dns.ResponseWriter, request *dns.Msg) (int, error) {
	if request.Question[0].Qtype != dns.TypePTR {
		return plugin.NextOrFailure(v6ptr.Name(), v6ptr.Next, ctx, writer, request)
	}

	var responsePtrValue string

	if ptrValue, found := v6ptr.Presets[request.Question[0].Name]; found {
		responsePtrValue = ptrValue
	} else {
		responsePtrValue = request.Question[0].Name
		responsePtrValue = RemoveIP6DotArpa(responsePtrValue)
		responsePtrValue = RemoveDots(responsePtrValue)
		responsePtrValue = ReverseString(responsePtrValue)
		responsePtrValue += "." + v6ptr.Suffix + "."
	}

	message := new(dns.Msg)
	message.SetReply(request)
	hdr := dns.RR_Header{Name: request.Question[0].Name, Ttl: v6ptr.TTL, Class: dns.ClassINET, Rrtype: dns.TypePTR}
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

func (a AutoIPv6PTR) Name() string { return AUTOIPV6PTR_PLUGIN_NAME }

func init() {
	plugin.Register(AUTOIPV6PTR_PLUGIN_NAME, setup)
}

func setup(c *caddy.Controller) error {
	v6ptr := AutoIPv6PTR{}
	v6ptr.TTL = 900

	for c.Next() {
		switch c.Val() {
		case "presetsfile":
			if err := parsePresetsFile(c.RemainingArgs()[0], &v6ptr); err != nil {
				return plugin.Error(AUTOIPV6PTR_PLUGIN_NAME, err)
			}

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

func parsePresetsFile(filepath string, v6ptr *AutoIPv6PTR) error {
	v6ptr.Presets = make(map[string]string)

	file, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0

	for scanner.Scan() {
		counter++
		presets := strings.Split(scanner.Text(), ";")

		if len(presets) == 2 {
			ip6ArpaValue, reverseError := dns.ReverseAddr(presets[0])

			if reverseError != nil {
				return reverseError
			} else {
				v6ptr.Presets[ip6ArpaValue] = presets[1] + "."
			}
		} else {
			return errors.New(fmt.Sprintf("Presets error: Two items expected in line %d", counter))
		}
	}

	return nil
}
