package DLinkDGS3100

import (
	"github.com/ircop/discoverer/base"
	"regexp"
	"fmt"
	"net"
	"strings"
)

// GetIps for DLinkDGS3100
func (p *Profile) GetIps() ([]discoverer.IPInterface, error) {
	addresses := make([]discoverer.IPInterface, 0)
	p.Log("Starting DLinkDGS3100.GetIps()")

	result, err := p.Cli.Cmd("show ipif")
	if err != nil {
		panic(err)
	}
	p.Debug(result)


	re, err := regexp.Compile(`IP Address\s+:\s+(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)([^\n]+)?\nSubnet Mask\s+:\s+(?P<mask>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)([^\n]+)?\nVlan name\s+:\s+(?P<ifname>[^\s]+)`)
	if err != nil {
		return addresses, fmt.Errorf("Cannot compile ipif regex: %s", err.Error())
	}

	out := p.ParseSingle(re, result)
	ip := net.ParseIP(out["ip"])
	mask := net.ParseIP(out["ip"])
	ifname := strings.Trim(out["ifname"], " ")
	if ip == nil || mask == nil {
		return addresses, fmt.Errorf("Cannot parpse ip/mask (%s/%s)", out["ip"], out["mask"])
	}
	if ifname == "" {
		return addresses, fmt.Errorf("Cannot parse ipif vlan name (%s)", ifname)
	}

	addr := discoverer.IPInterface{
		Interface:ifname,
		IP:ip,
		Mask:mask,
	}
	addresses = append(addresses, addr)

	return addresses, nil
}
