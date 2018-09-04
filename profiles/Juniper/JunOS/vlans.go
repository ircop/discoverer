package JunOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"strings"
	"strconv"
)

// GetVlans for JunOS
func (p *Profile) GetVlans() ([]dproto.Vlan, error) {
	p.Log("Starting JunOS.GetVlans()")
	vlans := make([]dproto.Vlan, 0)

	patterns := make(map[string]string)
	patterns["split"] = `(?ms:^VLAN:)`
	patterns["vlan"] = `(?ms:^VLAN:\s+(?P<name>[^,]+),.+\n802.1Q Tag:\s(?P<tag>\d+)[^\n]+\n[^\n]+\n[^\n]+.\n(?P<ifstring>.*))`
	patterns["ifs"] = `(?ms:^\s+(?P<ifname>[^,]+).\s+(?P<mode>[^,]+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return vlans, err
	}

	result, err := p.Cli.Cmd("show vlans extensive")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'show vlans extensive': %s", err.Error())
	}
	p.Debug(result)

	parts := regexps["split"].Split(result, -1)
	for _, part := range parts {
		part = "VLAN:" + part
		out := p.ParseSingle(regexps["vlan"], part)

		name := strings.Trim(out["name"], " ")
		vidStr := strings.Trim(out["tag"], " ")
		if name == "" || vidStr == "" {
			p.Log("WARNING! Cannot parse vlan name/id (%s/%s)", name, vidStr)
			p.Log(part)
			continue
		}
		vid, err := strconv.ParseInt(vidStr, 10, 64)
		if err != nil {
			p.Log("WARNING! Cannot parse vlan id (%s)", vidStr)
			continue
		}

		vlan := dproto.Vlan{
			Name:name,
			ID:vid,
		}

		ifstring := strings.Trim(out["ifstring"], "\n")
		out2 := p.ParseMultiple(regexps["ifs"], ifstring)
		for _, part2 := range out2 {
			ifname := strings.Trim(part2["ifname"], " ")
			ifname = strings.Replace(ifname, "*", "", -1)
			ifname = strings.Replace(ifname, ".0", "", -1)
			if ifname == "" {
				p.Log("Warning! Empty ifname (vlan '%s') (%+v)", name, part2)
				continue
			}

			mode := strings.Trim(part2["mode"], " ")
			switch mode {
			case "tagged":
				vlan.TrunkPorts = append(vlan.TrunkPorts, ifname)
				break;
			case "untagged":
				vlan.AccessPorts = append(vlan.AccessPorts, ifname)
				break;
			default:
				p.Log("WARNING! Unknown vlan mode '%s' (port '%s', vlan '%s')", mode, ifname, name)
				break
			}
		}

		vlans = append(vlans, vlan)
	}

	p.Debug("COUNT: %d\n", len(vlans))

	return vlans, nil
}
