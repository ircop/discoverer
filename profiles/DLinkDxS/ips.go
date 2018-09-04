package DLinkDxS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"net"
	"regexp"
	"strings"
)

// GetIps for DLinkDxS
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	ipifs := make([]*dproto.Ipif, 0)
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
	re2, err := regexp.Compile(`(?:Interface Name|IP Interface)\s+:\s+(?P<ifname>\S+)\s*\nVLAN Name\s+:\s+(?P<vlan_name>\S+)\s*\n([^\n]+\n)?([^\n]+\n)?IP(v4)? Address\s+:\s+(?P<cidr>\S+)\s+\(\S+\)\s*\n`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile ipifs regex")
	}

	out := p.ParseMultiple(re, result)
	if len(out) > 0 {
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

			Interface := dproto.Ipif{
				Interface: vlan,
				IP:        ip.String(),
				Mask:      mask.String(),
			}
			ipifs = append(ipifs, &Interface)
		}
	}

	//3200 c1
	out = p.ParseMultiple(re2, result)
	if len(out) > 0 {
		for _, iface := range out {
			name := strings.Trim(iface["ifname"], " ")
			cidrString := strings.Trim(iface["cidr"], " ")
			vlan := strings.Trim(iface["vlan_name"], " ")
			if name == "" || cidrString == "" || vlan == "" {
				p.Log("Something wrong: no name/cidr/vlan for IPIF (%s/%s/%s/%s)", name, cidrString, vlan)
				continue
			}

			ip, ipnet, err := net.ParseCIDR(cidrString)
			if err != nil {
				p.Log("Cannot parse cidr '%s': %s", cidrString, err.Error())
				continue
			}

			maskString := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
			mask := net.ParseIP(maskString)
			if mask == nil {
				continue
			}
			intf := dproto.Ipif{
				Interface:vlan,
				IP:ip.String(),
				Mask:mask.String(),
			}

			ipifs = append(ipifs, &intf)
		}
	}

	return ipifs, nil
}
