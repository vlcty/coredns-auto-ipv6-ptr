package autoipv6ptr

import (
	"os"
	"bufio"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin("plugin", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	plugin := AutoIPv6PTR{}

	presetsFilePath := ""

	for c.Next() {
		switch c.Val() {
		case "presetsfile":
			presetsFilePath = c.RemainingArgs()[0]

		case "suffix":
			plugin.Suffix = c.RemainingArgs()[0]

		default:
			// Maybe log something? :-)
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		plugin.Next = next
		return plugin
	})

	plugin.Presets = make(map[string]string)

	if len(presetsFilePath) != 0 {
		file, err := os.Open(presetsFilePath)
	    if err != nil {
	        return err
	    }
	    defer file.Close()

	    scanner := bufio.NewScanner(file)

	    for scanner.Scan() {
			presets := strings.Split(scanner.Text(), ";")

			plugin.Presets[presets[0]] = presets[1] + "."
	    }
	}

	return nil
}
