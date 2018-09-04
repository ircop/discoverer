package CiscoIOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"net"
	"strings"
)

// GetIps for CiscoIOS
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	addresses := make([]*dproto.Ipif, 0)
	p.Log("Starting CiscoIOS.GetIps()")

	result, err := p.Cli.Cmd("show ip int")
	if err != nil {
		return addresses, err
	}
	p.Debug(result)

	patterns := make(map[string]string)
	patterns["split"] = `Local Proxy ARP`
	patterns["ifname"] = `(?m:^(\s+)?(?P<ifname>.+?)\s+is\s+([^\s]+),\s+line protocol[^\n]+\n)`
	patterns["cidrs"] = `(Internet address is|Secondary address) (?P<cidr>\S+)\n`
	patterns["convert"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return addresses, err
	}

	parts := regexps["split"].Split(result, -1)
	for _, part := range parts {
		out := p.ParseSingle(regexps["ifname"], part)
		ifname := strings.Trim(out["ifname"], " ")
		if ifname == "" {
			continue
		}
		ifname, _ = p.ConvertIfname(ifname, regexps["convert"])

		cidrs := make([]string, 0)
		outx := p.ParseMultiple(regexps["cidrs"], part)
		for _, x := range outx {
			cidr := strings.Trim(x["cidr"], " ")
			cidrs = append(cidrs, cidr)

			ip, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				p.Log("Cannot parse cidr '%s': %s", cidr, err.Error())
				continue
			}
			maskString := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
			mask := net.ParseIP(maskString)
			if mask == nil {
				continue
			}
			intf := dproto.Ipif {
				Interface:ifname,
				IP:ip.String(),
				Mask:mask.String(),
			}

			addresses = append(addresses, &intf)
		}
	}

	return addresses, nil
}
