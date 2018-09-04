package JunOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"regexp"
	"strings"
)

// GetLldp for JunOS profile
func (p *Profile) GetLldp() ([]dproto.LldpNeighbor, error) {
	p.Debug("starting JunOS.GetLldp()")
	neighbors := make([]dproto.LldpNeighbor, 0)

	result, err := p.Cli.Cmd("show lldp neighbors")
	if err != nil {
		return neighbors, fmt.Errorf("Cannot 'show lldp neighbors': %s", err.Error())
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?m:Chassis ID\s+:\s+(?P<cid>[^\s]+)\n.+\nPort ID\s+:\s+(?P<pid>[^\s]+)\n)`)
	if err != nil {
		return neighbors, fmt.Errorf("Cannot compile neighbor regex: %s", err.Error())
	}

	rows := text.ParseTable(result, `^Local`, "", true)
	for _, row := range rows {
		ifname := strings.Trim(row[0], " ")
		ifname = strings.Replace(ifname, ".0", "", -1)
		if ifname == "" {
			p.Log("WARNING! Empty interface name (%+v)", row)
			continue
		}

		r, err := p.Cli.Cmd("show lldp neighbors interface "+ifname)
		if err != nil {
			p.Log("Error! Cannot 'show lldp neighbors interface %s': %s", ifname, err.Error())
			continue
		}
		p.Debug(r)

		out := p.ParseSingle(re, r)
		cid := strings.Trim(out["cid"], " ")
		pid := strings.Trim(out["pid"], " ")
		if cid == "" || pid == "" {
			p.Log("Warning! Cannot parse cid/pid (%s/%s)", cid,pid)
			continue
		}

		nei := dproto.LldpNeighbor{
			LocalPort:ifname,
			ChassisID:cid,
			PortID:pid,
		}
		neighbors = append(neighbors, nei)
	}

	return neighbors, nil
}
