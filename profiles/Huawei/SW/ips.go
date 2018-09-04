package HuaweiSW

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"net"
	"strings"
)

// GetIps for HuaweiSW
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	addresses := make([]*dproto.Ipif, 0)
	p.Log("Starting HuaweiSW.GetIps()")

	patterns := make(map[string]string)
	patterns["ifname"] = `(?mi:^(?P<ifname>[^\s]+)\scurrent state)`
	patterns["cidr"] = `(?mi:^Internet Address is\s(?P<cidr>[^\s\n]+)(\sSub)?)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return addresses, err
	}

	result, err := p.Cli.Cmd("display ip interface")
	if err != nil {
		return addresses, err
	}
	p.Debug(result)

	parts, err := text.SplitByParts(result, `^[^\s]+\s+current state`)
	if err != nil {
		return addresses, fmt.Errorf("Cannot split string by ip interfaces: %s", err.Error())
	}

	for _, part := range parts {
		out := p.ParseSingle(regexps["ifname"], part)
		ifname := strings.Trim(out["ifname"], " ")
		if ifname == "" {
			continue
		}
		out2 := p.ParseMultiple(regexps["cidr"], part)
		for _, o := range out2 {
			cidr := strings.Trim(o["cidr"], " ")
			ip, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				p.Log("Error parsing cidr '%s' (%s): %s", cidr, ifname, err.Error())
				continue
			}

			maskString := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
			mask := net.ParseIP(maskString)
			if mask == nil {
				p.Log("Error parsing mask for '%s'", ifname)
				continue
			}

			iface := dproto.Ipif{
				Interface:ifname,
				IP:ip.String(),
				Mask:mask.String(),
			}
			addresses = append(addresses, &iface)
		}
	}


	return addresses, nil
}
