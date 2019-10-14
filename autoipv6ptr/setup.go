package autoipv6ptr

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"errors"
	"fmt"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/caddyserver/caddy"
)

const AUTOIPV6PTR_PLUGIN_NAME string = "autoipv6ptr"

func init() {
	caddy.RegisterPlugin(AUTOIPV6PTR_PLUGIN_NAME, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
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
			v6ptr.Presets[presets[0]] = presets[1] + "."
		} else {
			return errors.New(fmt.Sprintf("Presets error: Two items expected in line %d", counter))
		}
	}

	return nil
}
