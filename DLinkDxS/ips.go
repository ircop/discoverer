package DLinkDxS

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"regexp"
	"strings"
	"net"
)

// GetIps for DLinkDxS
func (p *Profile) GetIps() ([]discoverer.IPInterface, error) {
	ipifs := make([]discoverer.IPInterface, 0)
	p.Log("Starting DLinkDxS.GetIps()")

	result, err := p.Cli.Cmd("show ipif")
	if err != nil {
		return ipifs, fmt.Errorf("Error during 'show ipif': %s", err.Error())
	}
	p.Debug(result)


	re, err := regexp.Compile(`(?:Interface Name|IP Interface)\s+:\s+(?P<ifname>\S+)\s*\nIP Address\s+:\s+(?P<ip>\S+)\s+\(\S+\)\s*\nSubnet Mask\s+:\s+(?P<mask>\S+)\s*\n(Interface Admin[^\n]+\n)?VLAN Name\s+:\s+(?P<vlan_name>\S+)\s*\n`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile ipifs regex")
	}

	out := p.ParseMultiple(re, result)
	for _, iface := range out {
		name := strings.Trim(iface["ifname"], " ")
		ipstring := strings.Trim(iface["ip"], " ")
		maskstring := strings.Trim(iface["mask"], " ")
		vlan := strings.Trim(iface["vlan_name"], " ")
		if name == "" || ipstring == "" || maskstring == "" || vlan == "" {
			p.Log("Something wrong: no name/ip/mask/vlan for IPIF (%s/%s/%s/%s)", name, ipstring, maskstring, vlan)
			continue
		}

		ip := net.ParseIP(ipstring)
		if ip == nil {
			p.Log("Error: wrong ip address '%s'", ipstring)
			continue
		}
		mask := net.ParseIP(maskstring)
		if mask == nil {
			p.Log("Error: wrong mask '%s'", maskstring)
		}


		Interface := discoverer.IPInterface{
			Interface:vlan,
			IP: ip,
			Mask:mask,
		}
		ipifs = append(ipifs, Interface)
	}

	return ipifs, nil
}
