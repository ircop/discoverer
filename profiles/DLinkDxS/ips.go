package DLinkDxS

import (
	"fmt"
	"github.com/ircop/dproto"
	"net"
	"regexp"
	"strings"
)

// GetIps for DLinkDxS
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	ipifs := make([]*dproto.Ipif, 0)
	p.Log("Starting DLinkDxS.GetIps()")

	patterns := make(map[string]string)
	patterns["vlan"] = `(?msi:^VLAN( name)?\s+:(\s+)(?P<vlan>[^\s\n]+))`
	patterns["ip"] = `(?msi:^ip(v4)? address\s+:\s+(?P<ip>\b(?:\d{1,3}\.){3}\d{1,3}\b)(\s|\n))`
	patterns["mask"] = `(?msi:^subnet mask\s+:\s+(?P<mask>\b(?:\d{1,3}\.){3}\d{1,3}\b)(\s|\n))`
	patterns["cidr"] = `(?msi:^ip(v4)? address\s+:\s+(?P<cidr>\b(?:\d{1,3}\.){3}\d{1,3}\b\/\d+)(\s|\n))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile regexps: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show ipif")
	if err != nil {
		return ipifs, fmt.Errorf("Error during 'show ipif': %s", err.Error())
	}
	p.Debug(result)

	result = strings.ToLower(result)
	parts := strings.Split(result, "ip interface")

	for i := range parts {
		part := parts[i]
		out := p.ParseSingle(regexps["vlan"], part)
		vlan := strings.Trim(out["vlan"], " ")
		if vlan == "" {
			continue
		}

		out = p.ParseSingle(regexps["ip"], part)
		ipstring := strings.Trim(out["ip"], " ")
		out = p.ParseSingle(regexps["mask"], part)
		maskstring := strings.Trim(out["mask"], " ")
		out = p.ParseSingle(regexps["cidr"], part)
		cidr := strings.Trim(out["cidr"], " ")

		if cidr == "" && (maskstring == "" && ipstring == "") {
			continue
		}

		iface := dproto.Ipif{
			Interface:vlan,
		}

		// find out ip&mask
		if cidr != "" {
			ip, net, err := net.ParseCIDR(cidr)
			if err == nil {
				iface.IP = ip.String()
				mask := fmt.Sprintf("%d.%d.%d.%d", net.Mask[0], net.Mask[1], net.Mask[2], net.Mask[3])
				iface.Mask = mask
				ipifs = append(ipifs, &iface)
				continue
			}
		}

		ip := net.ParseIP(ipstring)
		if ip == nil {
			continue
		}
		mask := net.ParseIP(maskstring)
		if mask == nil {
			continue
		}
		iface.IP = ip.String()
		iface.Mask = mask.String()
		ipifs = append(ipifs, &iface)
	}

	return ipifs, nil
}

