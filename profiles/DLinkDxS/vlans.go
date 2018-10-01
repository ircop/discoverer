package DLinkDxS

import (
	"fmt"
	"github.com/ircop/dproto"
	"regexp"
	"strconv"
	"strings"
)

// GetVlans for DLinkDxS
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting DLinkDxS.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)

	result, err := p.Cli.Cmd("show vlan")
	if err != nil {
		return vlans, fmt.Errorf("Cannot get vlans: %s", err.Error())
	}
	p.Debug(result)

	// REGEX FOR D(G|E)Ss
	re, err := regexp.Compile(`(?i:^VID\s+:\s+(?P<vlan_id>\d+)\s+VLAN Name\s+:(?P<vlan_name>.*?)\nVLAN Type\s+:\s+(?P<vlan_type>\S+)\s*?(Adv[^\n]+)?\n((VLAN )?Advertisement\s+:\s+\S+\s*\n)?(Member[^\n]+\nStatic[^\n]+\n)?(Current )?Tagged Ports\s*:(\s+)?(?P<tagged_ports>[^\n]+)?\n(Current )?Untagged Ports\s*:(?P<untagged_ports>[^\n]+)?\n)`)
	// REGEX FOR DXSs
	re2, err2 := regexp.Compile(`^VID\s+:\s+(?P<vlan_id>\d+)\s+VLAN Name\s+:(?P<vlan_name>.*?)\nVLAN TYPE\s+:[^\n]+\s+UserDefinedPid[^\n]+\s+Encap[^\n]+\s+Member ports\s+:\s+(?P<member_ports>[^\n]+)?\s+Static ports[^\n]+\s+Untagged ports\s+:(\s+)?(?P<untagged_ports>[^\n]+)?\n`)
	//re2, err :=
	if err != nil {
		return vlans, fmt.Errorf("Cannot compile vlan regex: %s", err.Error())
	}
	if err2 != nil {
		return vlans, fmt.Errorf("Cannot compile dxs vlan regex: %s", err2.Error())
	}

	parts := strings.Split(result, "\n\n")
	for _, part := range parts {
		part = strings.Trim(part, "\n")
		//fmt.Printf(part)
		out := p.ParseSingle(re, part)
		if len(out) > 1 {
			vidStr := strings.Trim(out["vlan_id"], " ")
			name := strings.Trim(out["vlan_name"], " ")
			taggedStr := strings.Trim(out["tagged_ports"], " ")
			untaggedStr := strings.Trim(out["untagged_ports"], " ")
			if vidStr == "" || name == "" {
				continue
			}

			vid, err := strconv.ParseInt(vidStr, 10, 64)
			if err != nil {
				p.Log("Failed to parse vlan id '%s'", vidStr)
				continue
			}

			vlan := dproto.Vlan{
				Name:        name,
				ID:          vid,
				AccessPorts: p.ExpandInterfaceRange(untaggedStr),
				TrunkPorts:  p.ExpandInterfaceRange(taggedStr),
			}

			vlans = append(vlans, &vlan)
			continue
		}

		out = p.ParseSingle(re2, part)
		if len(out) > 1 {
			vidStr := strings.Trim(out["vlan_id"], " ")
			name := strings.Trim(out["vlan_name"], " ")
			allStr := strings.Trim(out["member_ports"], " ")
			untaggedStr := strings.Trim(out["untagged_ports"], " ")
			if vidStr == "" || name == "" {
				continue
			}

			vid, err := strconv.ParseInt(vidStr, 10, 64)
			if err != nil {
				p.Log("Failed to parse DXS vlan id '%s'", vidStr)
				continue
			}

			untag := p.ExpandInterfaceRange(untaggedStr)
			all := p.ExpandInterfaceRange(allStr)
			tag := make([]string,0)
			// tag = all minus untag
			for _, port := range all {
				isUntag := false
				for _, u := range untag {
					if u == port {
						// this is untag
						isUntag = true
						break
					}
				}
				if !isUntag {
					tag = append(tag, port)
				}
			}

			vlan := dproto.Vlan{
				Name: name,
				ID: vid,
				AccessPorts: untag,
				TrunkPorts: tag,
			}
			vlans = append(vlans, &vlan)
			continue
		}
	}

	return vlans, nil
}
