package EltexMES

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"github.com/ircop/discoverer/util/text"
	"strconv"
)

// GetVlans for EltexMES
func (p *Profile) GetVlans() ([]discoverer.Vlan, error) {
	p.Log("Starting DLinkDxS.GetVlans()")
	vlans := make([]discoverer.Vlan, 0)

	result, err := p.Cli.Cmd("show vlan")
	if err != nil {
		return vlans, fmt.Errorf("Cannot execute 'show vlan': %s", err.Error())
	}
	p.Debug(result)

	rows := text.ParseTable(result, "^--", "", false)

	for _, row := range rows {
		vidStr := row[0]
		tag := row[2]
		unt := row[3]
		vid, err := strconv.ParseInt(vidStr, 10, 64)
		if err != nil {
			p.Log("Cannot parse vlan id '%s'", vidStr)
			continue
		}
		untagged := p.ExpandInterfaceRange(unt)
		tagged := p.ExpandInterfaceRange(tag)

		vlan := discoverer.Vlan{
			ID:vid,
			Name:vidStr,
			TrunkPorts:tagged,
			AccessPorts:untagged,
		}
		vlans = append(vlans, vlan)
	}

	return vlans, nil
}
