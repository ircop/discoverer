package EltexMES

import (
	"fmt"
	"github.com/ircop/dproto"
	"github.com/ircop/discoverer/util/text"
	"net"
	"regexp"
	"strings"
)

// GetIps for EltexMES
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	ipifs := make([]*dproto.Ipif, 0)
	p.Log("Starting EltexMES.GetIps()")

	result, err := p.Cli.Cmd("sh ip interface")
	if err != nil {
		return ipifs, fmt.Errorf("Cannot 'sh ip interface': %s", err.Error())
	}
	p.Debug(result)

	// some models shows both IP and GW tables; GW first. Cut it
	reIfHeader, err := regexp.Compile(`IP Address\s+I\/F\s+`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile ip/gw split regex: %s", reIfHeader)
	}
	parts := reIfHeader.Split(result, -1)
	if len(parts) < 2 {
		p.Log("Cannot split ip/gw tables: split result: %d", len(parts))
	}
	result = parts[1]


	reSvi, err := regexp.Compile(`vlan\s\d+`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile svi regex")
	}

	rows := text.ParseTable(result, "^--", "", false, false)
	for _, row := range rows {
		if len(row) < 4 {
			continue
		}

		ifname := strings.Trim(row[1], " ")
		if !reSvi.Match([]byte(ifname)) {
			continue
		}
		ifname = strings.Replace(ifname, " ", "", -1)

		cidr := strings.Trim(row[0], " ")
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return ipifs, fmt.Errorf("Cannot parse cidr '%s': %s", cidr, err.Error())
		}
		if ip.String() == "0.0.0.0" {
			continue
		}

		maskString := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
		mask := net.ParseIP(maskString)
		if mask == nil {
			continue
		}
		intf := dproto.Ipif{
			Interface:ifname,
			IP:ip.String(),
			Mask:mask.String(),
		}
		ipifs = append(ipifs, &intf)
	}

	return ipifs, nil
}