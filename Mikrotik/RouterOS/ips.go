package RouterOS

import (
	"fmt"
	"github.com/ircop/discoverer/base"
	"net"
	"strings"
)

// GetIps for RouterOS
func (p *Profile) GetIps() ([]discoverer.IPInterface, error) {
	ipifs := make([]discoverer.IPInterface, 0)
	p.Log("Starting RouterOS.GetIps()")

	patterns := make(map[string]string, 0)
	patterns["interfaces"] = `(?m:\s\d+\s+address=(?P<cidr>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d+))[^\n]+actual-interface=(?P<ifname>[^\s]+)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return ipifs, err
	}

	result, err := p.Cli.Cmd("/ip address print detail without-paging")
	if err != nil {
		return ipifs, nil
	}
	p.Debug(result)

	parts := p.ParseMultiple(regexps["interfaces"], result)
	for _, part := range parts {
		cidrString := strings.Trim(part["cidr"], " ")
		iface := strings.Trim(part["ifname"], " ")
		if cidrString == "" || iface == "" {
			p.Log("Warning! Cannot parse ip/interface (%s/%s)", cidrString, iface)
			continue
		}

		ip, network, err := net.ParseCIDR(cidrString)
		if err != nil {
			p.Log("Error! Cannot parse cidr '%s': %s", cidrString, err.Error())
			continue
		}
		maskString := fmt.Sprintf("%d.%d.%d.%d", network.Mask[0], network.Mask[1], network.Mask[2], network.Mask[3])
		mask := net.ParseIP(maskString)

		if ip == nil || mask == nil {
			p.Log("Cannot parse ip/net of '%s'", cidrString)
			continue
		}

		ipif := discoverer.IPInterface{
			Interface:iface,
			IP:ip,
			Mask:mask,
		}
		ipifs = append(ipifs, ipif)
	}

	return ipifs, nil
}
