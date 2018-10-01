package DLinkDGS3100

import (
	"fmt"
	"github.com/ircop/dproto"
	"regexp"
	"strconv"
	"strings"
)

// GetVlans for DLinkDGS3100
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting DLinkDGS3100.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)

	result, err := p.Cli.Cmd("show vlan")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'show vlan': %s", err.Error())
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?msi:VID\s+:\s+(?P<vid>\d+)\s+VLAN Name\s+:\s+(?P<name>[^\n]+)\nVLAN TYPE[^\n]+\nMember ports\s+:(?P<all_ports>[^\n]+)?\nSt[^\n]+\nUntagged ports\s+:(?P<untagged>[^\n]+)?)`)
	if err != nil {
		return vlans, fmt.Errorf("Cannot compile regex: %s", err.Error())
	}

	out := p.ParseMultiple(re, result)
	for _, part := range out {
		vidStr := strings.Trim(part["vid"], " ")
		name := strings.Trim(part["name"], " ")
		all := p.ExpandInterfaceRange(strings.Trim(part["all_ports"], " "))
		untag := p.ExpandInterfaceRange(strings.Trim(part["untagged"], " "))
		if vidStr == "" || name == "" {
			p.Log("Cannot parse vlan (vid, name = '%s', '%s')", vidStr, name)
			continue
		}

		vid, err := strconv.ParseInt(vidStr, 10, 64)
		if err != nil {
			p.Log("Cannot parse vlan id to integer (%s): %s", vidStr, err.Error())
			continue
		}

		// tag?
		tag := make([]string, 0)
		allLoop:
		for _, port := range all {
			for _, uport := range untag {
				if port == uport {
					continue allLoop
				}
			}
			// not untagged
			tag = append(tag, port)
		}

		vlan := dproto.Vlan{
			ID:vid,
			Name:name,
			AccessPorts:untag,
			TrunkPorts:tag,
		}
		vlans = append(vlans, &vlan)
	}

	return vlans, nil
}
