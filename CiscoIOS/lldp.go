package CiscoIOS

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"strings"
)

// GetLldp for CiscoIOS profile
func (p *Profile) GetLldp() ([]discoverer.LldpNeighborship, error) {
	p.Debug("starting CiscoIOS.GetLldp()")
	neighbors := make([]discoverer.LldpNeighborship, 0)

	out, err := p.Cli.Cmd("show lldp neighbors")
	if err != nil {
		return neighbors, fmt.Errorf("Cannot get lldp neighbors: %s", err.Error())
	}
	p.Debug(strings.Replace(out, "%", "%%", -1))

	patterns := make(map[string]string)
	//patterns["locals"] = `(?i:\n[^\n]+\s+(?P<local_if>(?:Fa|Gi|Te)\d+[\d/\.]*).+)`
	patterns["locals"] = `(?msi:^[^\n\s]+\s+(?P<local_if>(?:Fa|Gi|Te)\d+[\d/\.]*)\s+([^\n\s]+\s+)?[^\n]+\n)`
	patterns["chassis"] = `(?i:Chassis id:\s*(?P<chassis_id>\S+)(\s+)?\n)`
	patterns["port"] = `(?i:Port id:\s*(?P<port_id>[^\n\s]+)(\s+)?\n)`
	patterns["ifname"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return neighbors, err
	}

	locals := p.ParseMultiple(regexps["locals"], out)
	for _, l := range locals {
		// Loop trouhg local interfaces
		local := strings.Trim(l["local_if"], " ")
		if local == "" {
			p.Log("Something wrong: empty local LLDP interface")
			continue
		}

		// Get detailed output
		r, err := p.Cli.Cmd(fmt.Sprintf("show lldp neighbors %s detail", local))
		if err != nil {
			p.Log("Cannot get detail information about neighbor on interface '%s': %s", local, err.Error())
			continue
		}
		p.Debug(strings.Replace(r, "%", "%%", -1))

		out1 := p.ParseSingle(regexps["chassis"], r)
		out2 := p.ParseSingle(regexps["port"], r)
		cid := strings.Trim(out1["chassis_id"], " ")
		pid := strings.Trim(out2["port_id"], " ")
		if cid == "" || pid == "" {
			p.Log("Cannot parse chassis/port id on lldp interface '%s' ('%s'/'%s')", local, cid, pid)
			continue
		}

		local, err = p.ConvertIfname(local, regexps["ifname"])
		if err != nil {
			p.Log("Cannot convert local lldp ifname to short one (%s)", local)
			continue
		}

		item := discoverer.LldpNeighborship{
			PortName:local,
		}
		members := make([]discoverer.LldpNeighbor,0)
		members = append(members, discoverer.LldpNeighbor{ChassisID:cid,PortID:pid})
		item.Members = members

		neighbors = append(neighbors, item)
	}

	return neighbors, nil
}