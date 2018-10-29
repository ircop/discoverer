package JunOS

import (
	"fmt"
	"github.com/ircop/dproto"
	"regexp"
	"strings"
	"strconv"
)

// GetVlans for JunOS
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting JunOS.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)

	patterns := make(map[string]string)
	patterns["split"] = `(?ms:^VLAN:)`
	patterns["platform"] = `(?ms:Model:\s+(?P<platform>\S+)\n(Junos:\s+(?P<version>[^\n]+)|JUNOS Base OS boot \[(?P<version2>[^\]]+)\]))`
	patterns["vlan"] = `(?ms:^VLAN:\s+(?P<name>[^,]+),.+\n802.1Q Tag:\s(?P<tag>\d+)[^\n]+\n[^\n]+\n[^\n]+.\n(?P<ifstring>.*))`
	patterns["ifs"] = `(?ms:^(\s+)?(?P<ifname>[^,]+).\s+(?P<mode>[^,]+),)`
	patterns["expand_ifs"] = `(?ms:Number of interfaces[^\n]+\n(?P<ifs>.+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return vlans, err
	}

	result, err := p.Cli.Cmd("show version")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'show version': %s", err.Error())
	}
	out := p.ParseSingle(regexps["platform"], result)
	if strings.Contains(out["platform"], "mx") {
		return p.GetVlansMX()
	}


	result, err = p.Cli.Cmd("show vlans extensive")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'show vlans extensive': %s", err.Error())
	}
	p.Debug(result)

	result = strings.Replace(result, "show vlans extensive", "", -1)

	parts := regexps["split"].Split(result, -1)
	for _, part := range parts {
		if len(part) < 10 {
			continue
		}
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

		//ifstring := strings.Trim(out["ifstring"], "\n")
		//fmt.Printf("%s\n", name)
		//fmt.Printf(ifstring)
		out = p.ParseSingle(regexps["expand_ifs"], part)
		ifstring := strings.Trim(out["ifs"], " ")

		//fmt.Printf(ifstring)
		//return vlans, nil
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
				p.Log(part)
				break
			}
		}

		vlans = append(vlans, &vlan)
	}

	//p.Debug("COUNT: %d\n", len(vlans))

	return vlans, nil
}

// GetVlans for MX series
func (p *Profile) GetVlansMX() ([]*dproto.Vlan, error) {
	p.Log("Starting JunOS.GetVlansMX()")
	vlans := make([]*dproto.Vlan, 0)

	//p.Debug("GETTING MX SERIES VLANS")
	//reStr := `(?msi:Logical interface (?P<ifname>[^\n\s]+)[^\n]+ifindex[^\n]+\n[^\n]+VLAN-Tag \[ 0x\d+\.(?P<vid>\d+) )`
	reStr := `(?msi:Logical interface (?P<ifname>[^\n\s\.]+)\.\d+[^\n]+ifindex[^\n]+(\n\s+Desc[^\n]+)?\n[^\n]+VLAN-Tag \[ 0x\d+\.(?P<vid>\d+) )`
	re, err := regexp.Compile(reStr)
	if err != nil {
		return vlans, fmt.Errorf("Cannot compile MX vlan regex: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show interfaces")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'show interfaces': %s", err.Error())
	}
	p.Debug(result)

	vlanports := make(map[int64][]string)

	out := p.ParseMultiple(re, result)
	for i := range out {
		ifname := strings.Trim(out[i]["ifname"], " ")
		if ifname == "" {
			continue
		}
		vidStr := strings.Trim(out[i]["vid"], " ")
		if vidStr == "" || vidStr == "0" {
			continue
		}

		vid, err := strconv.ParseInt(vidStr, 10, 64)
		if err != nil {
			continue
		}

		if ports, ok := vlanports[vid]; ok {
			// check if there is no this ports already
			found := false
			for n := range ports {
				if ports[n] == ifname {
					found = true
					break
				}
			}
			if !found {
				ports = append(ports, ifname)
				vlanports[vid] = ports
			}
		} else {
			ports = []string{ifname}
			vlanports[vid] = ports
		}
		//vlan := dproto.Vlan{}
		//p.Debug("%+#v", out[i])
	}

	for vid, ports := range vlanports {
		vlan := dproto.Vlan{
			Name: fmt.Sprintf("%d", vid),
			AccessPorts:make([]string,0),
			ID:vid,
			TrunkPorts:ports,
		}
		vlans = append(vlans, &vlan)
	}


	return vlans, nil
}
