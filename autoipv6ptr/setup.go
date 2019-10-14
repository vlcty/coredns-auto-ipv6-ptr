package autoipv6ptr

import (
	"os"
	"bufio"
	"strings"
	"strconv"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin("autoipv6ptr", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	v6ptr := AutoIPv6PTR{}
	v6ptr.TTL = 900

	presetsFilePath := ""

	for c.Next() {
		switch c.Val() {
		case "presetsfile":
			presetsFilePath = c.RemainingArgs()[0]

		case "suffix":
			v6ptr.Suffix = c.RemainingArgs()[0]

		case "ttl":
			possibleTTL := c.RemainingArgs()[0]
			ttl, err := strconv.ParseUint(possibleTTL, 10, 32)

			if err != nil {
				return plugin.Error("autoipv6ptr", err)
			} else {
				v6ptr.TTL = uint32(ttl)
			}

		default:
			// Maybe log something? :-)
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		v6ptr.Next = next
		return v6ptr
	})

	v6ptr.Presets = make(map[string]string)

	if len(presetsFilePath) != 0 {
		file, err := os.Open(presetsFilePath)
	    if err != nil {
	        return err
	    }
	    defer file.Close()

	    scanner := bufio.NewScanner(file)

	    for scanner.Scan() {
			presets := strings.Split(scanner.Text(), ";")

			v6ptr.Presets[presets[0]] = presets[1] + "."
	    }
	}

	return nil
}
