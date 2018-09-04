package JunOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"net"
	"regexp"
	"strings"
)

// GetIps for JunOS
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	ipifs := make([]*dproto.Ipif, 0)
	p.Log("Starting JunOS.GetIps()")

	result, err := p.Cli.Cmd("show interfaces terse")
	if err != nil {
		return ipifs, fmt.Errorf("Cannot 'show int terse': %s", err.Error())
	}
	p.Debug(result)

	reDestination, err := regexp.Compile(`(?m:Destination:\s+(?P<net>[^/]+)/(?P<shortmask>\d+),\sLocal:\s+(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b),)`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile destination regexp: %s", err.Error())
	}
	reSkip, err := regexp.Compile(`^(bme|jsrv)`)
	if err != nil {
		return ipifs, fmt.Errorf("Cannot compile if skip regex: %s", err.Error())
	}

	rows := text.ParseTable(result, `Interface\s+`, "", true)
	for _, row := range rows {
		if len(row) < 4 {
			p.Log("Warning! Interfaces row len is %d", len(row))
			continue
		}

		ifname := strings.Trim(row[0], " ")
		if ifname == "" {
			continue
		}
		if strings.HasSuffix(ifname, ".0") {
			continue
		}
		if reSkip.Match([]byte(ifname)) {
			continue
		}

		r, e := p.Cli.Cmd(fmt.Sprintf("show interfaces %s", ifname))
		if e != nil {
			p.Log("Error! Cannot 'show interfaces %s': %s", ifname, e.Error())
			continue
		}
		p.Debug(r)

		out := p.ParseMultiple(reDestination, r)
		for _, part := range out {
			subnet := strings.Trim(part["net"], " ")
			shortmask := strings.Trim(part["shortmask"], " ")
			ipString := strings.Trim(part["ip"], " ")
			ip := net.ParseIP(ipString)
			if subnet == "" || shortmask == "" || ip == nil {
				p.Log("WARNING! Cannot parse subnet/shortmask/ip (%s/%s/%s)", subnet, shortmask, ipString)
				continue
			}

			ip, network, err := net.ParseCIDR(fmt.Sprintf("%s/%s", ip.String(), shortmask))
			if err != nil {
				p.Log("Error! Cannot parse cidr from ip/network (%s/%s): %s", ip.String(), shortmask, err.Error())
				continue
			}
			maskString := fmt.Sprintf("%d.%d.%d.%d", network.Mask[0], network.Mask[1], network.Mask[2], network.Mask[3])
			mask := net.ParseIP(maskString)
			if mask == nil {
				continue
			}

			ipif := dproto.Ipif{
				Interface:ifname,
				IP:ip.String(),
				Mask:mask.String(),
			}
			ipifs = append(ipifs, &ipif)
		}
	}

	return ipifs, nil
}
