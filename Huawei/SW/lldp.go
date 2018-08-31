package HuaweiSW

import (
	"fmt"
	"github.com/ircop/discoverer/base"
	"github.com/ircop/discoverer/util/text"
	"strings"
)

// GetLldp for HuaweiSW profile
func (p *Profile) GetLldp() ([]discoverer.LldpNeighbor, error) {
	p.Debug("starting HuaweiSW.GetLldp()")
	neighbors := make([]discoverer.LldpNeighbor, 0)

	patterns := make(map[string]string)
	patterns["ifname"] = `(?ms)^(?P<ifname>[^\s]+) has \d+ nei.+`
	patterns["neis"] = `(?ms)^Neighbor index[^\n]+\nChassis ty[^\n]+\nChassis ID\s+:(?P<cid>[^\n]+)\nPo[^\n]+\nPort ID\s+:(?P<pid>[^\n]+)\n`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return neighbors, nil
	}

	result, err := p.Cli.Cmd("display lldp neighbor")
	if err != nil {
		return neighbors, fmt.Errorf("Cannot 'display lldp neighbor': %s", err.Error())
	}
	p.Debug(result)

	parts, err := text.SplitByParts(result, `^[^\s]+ has \d+ nei`)
	if err != nil {
		return neighbors, err
	}

	for _, part := range parts {
		out1 := p.ParseSingle(regexps["ifname"], part)
		ifname := strings.Trim(out1["ifname"], " ")
		if ifname == "" {
			continue
		}
		out2 := p.ParseMultiple(regexps["neis"], part)
		for _, nei := range out2 {
			cid := strings.Trim(nei["cid"], " ")
			pid := strings.Trim(nei["pid"], " ")
			if cid == "" || pid == "" {
				p.Debug("Warning: no cid/pid for %s (%s/%s)", ifname, cid, pid)
				continue
			}

			//p.Debug("%s: %s/%s", ifname, cid, pid)
			neighbor := discoverer.LldpNeighbor{
				LocalPort:ifname,
				ChassisID:cid,
				PortID:pid,
			}
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors, nil
}
