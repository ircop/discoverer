package DLinkDGS3100

import (
	"fmt"
	"github.com/ircop/dproto"
	"github.com/ircop/discoverer/util/text"
	"regexp"
	"strconv"
	"strings"
)

	// GetInterfaces for DLinkDGS3100 profile
	func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
		p.Debug("Starting DLinkDxS.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	result, err := p.Cli.Cmd("sh ports description")
	if err != nil {
		p.Debug(result)
		return interfaces, fmt.Errorf("Error getting interfaces: %s", err.Error())
	}
	p.Debug(result)

	rows := text.ParseTable(result, `^----`, "", false, false)
	for _, row := range rows {
		if len(row) < 1 {
			continue
		}
		ifname := strings.Trim(row[0], " ")
		desc := ""
		if len(row) > 1 {
			desc = strings.Trim(row[1], " ")
		}

		newInt := dproto.Interface{
			Name:        ifname,
			Shortname:   ifname,
			Description: desc,
			Type:        dproto.InterfaceType_PHISYCAL,
		}
		interfaces[ifname] = &newInt
	}

	// todo: portchannels
	rePos, err := regexp.Compile(`(?msi:Group ID\s+:\s+(?P<id>\d+)\nMember Port\s+:\s+(?P<ports>[^\s]+))`)
	if err != nil {
		return interfaces, fmt.Errorf("Cannot compile port-channel regex: %s", err.Error())
	}

	result, err = p.Cli.Cmd("show link_aggregation")
	if err != nil {
		return interfaces, fmt.Errorf("Cannot 'show link_aggregation': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseMultiple(rePos, result)
	for _, part := range out {
		gid := strings.Trim(part["id"], " ")
		ports := strings.Trim(part["ports"], " ")
		if _, err = strconv.ParseInt(gid, 10, 64); err != nil {
			p.Log("Cannot parse port-channel id '%s'", gid)
			continue
		}
		ifaces := p.ExpandInterfaceRange(ports)
		newInt := dproto.Interface{
			Name:      "ch"+gid,
			Shortname: "ch"+gid,
			LldpID:    "ch"+gid,
			PoMembers: ifaces,
			Type:      dproto.InterfaceType_AGGREGATED,
		}
		interfaces["ch"+gid] = &newInt
	}

	ipifs, err := p.getIpifs()
	if err != nil {
		return interfaces, err
	}
	for name, _ := range ipifs {
		newInt := dproto.Interface{
			Name:      name,
			Shortname: name,
			Type:      dproto.InterfaceType_SVI,
			LldpID:    name,
		}
		interfaces[name] = &newInt
	}

	return interfaces, nil
}

func (p *Profile) getIpifs() (map[string]string, error) {
	ipifs := make(map[string]string)

	result, err := p.Cli.Cmd("show ipif")
	if err != nil {
		p.Debug(result)
		return ipifs, fmt.Errorf("Error getting ipifs: %s", err.Error())
	}
	p.Debug(result)

	reIpif, err := regexp.Compile(`(?msi:^Vlan name\s+:\s+(?P<ipif>[^\n]+)\n)`)
	if err != nil {
		return ipifs, fmt.Errorf("Failed to compile ipif regex: %s", err.Error())
	}

	out := p.ParseMultiple(reIpif, result)
	for _, part := range out {
		ipif := strings.Trim(part["ipif"], " ")
		if ipif != "" {
			ipifs[ipif] = ipif
		}
	}

	return ipifs, nil
}