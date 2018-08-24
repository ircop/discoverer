package DLinkDGS3100

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"github.com/ircop/discoverer/util/text"
	"strings"
	"regexp"
	"strconv"
)

	// GetInterfaces for DLinkDGS3100 profile
	func (p *Profile) GetInterfaces() (map[string]discoverer.Interface, error) {
		p.Debug("Starting DLinkDxS.GetInterfaces()")
	interfaces := make(map[string]discoverer.Interface)

	result, err := p.Cli.Cmd("sh ports description")
	if err != nil {
		p.Debug(result)
		return interfaces, fmt.Errorf("Error getting interfaces: %s", err.Error())
	}
	p.Debug(result)

	rows := text.ParseTable(result, `^----`, "")
	for _, row := range rows {
		if len(row) < 1 {
			continue
		}
		ifname := strings.Trim(row[0], " ")
		desc := ""
		if len(row) > 1 {
			desc = strings.Trim(row[1], " ")
		}

		newInt := discoverer.Interface{
			Name: ifname,
			Shortname: ifname,
			Description:desc,
			Type: discoverer.IntTypePhisycal,
		}
		interfaces[ifname] = newInt
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
		newInt := discoverer.Interface{
			Name: "ch"+gid,
			Shortname: "ch"+gid,
			LldpID: "ch"+gid,
			PoMembers: ifaces,
		}
		interfaces["ch"+gid] = newInt
	}

	return interfaces, nil
}
