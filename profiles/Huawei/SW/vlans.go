package HuaweiSW

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"strconv"
	"strings"
)

// GetVlans for HuaweiSW
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting HuaweiSW.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)

	patterns := make(map[string]string)
	patterns["cut"] = `(?mis:^(?P<content>VID\s+Type.+)VID\s+Status)`
	patterns["rmline"] = `(?mis:^-+(\s+)?\n)`
	patterns["untagsWithTags"] = `^UT:(?P<ifaces>.+)TG`
	patterns["untagsWithoutTags"] = `^UT:(?P<ifaces>.+)`
	patterns["tags"] = `TG:(?P<ifaces>.+)`
	patterns["ifname"] = `(?mis:(^|\s(\s+)?)(?P<ifname>[^\s\(\n]+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return vlans, nil
	}

	result, err := p.Cli.Cmd("display vlan")
	if err != nil {
		return vlans, fmt.Errorf("Cannot 'display vlan': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["cut"], result)
	result = strings.Trim(out["content"], " ")
	if result == "" {
		return vlans, fmt.Errorf("Cannot parse vlans content.")
	}
	result = regexps["rmline"].ReplaceAllString(result, "")
	p.Debug(result)


	rows := text.ParseTable(result, `^VID\s+Type\s+Ports`, "", false)
	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		vid, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			p.Log("Warning! Cannot parse vlan id '%s'", row[0])
			continue
		}
		portstring := row[2]
		//v59: []string{"59", "common", "UT:Eth-Trunk3(U)TG:XGE0/0/48(D)    Eth-Trunk0(U)"}

		// get tags first
		o := p.ParseSingle(regexps["tags"], portstring)
		tags := strings.Trim(o["ifaces"], " ")
		if tags != "" {
			o = p.ParseSingle(regexps["untagsWithTags"], portstring)
		} else {
			o = p.ParseSingle(regexps["untagsWithoutTags"], portstring)
		}
		untags := strings.Trim(o["ifaces"], " ")
		/*p.Debug("portstring: %s", portstring)
		p.Debug("untags: %s", untags)
		p.Debug("tags: %s", tags)*/

		// get 'UT:....'
		//o := p.ParseSingle(regexps["untags"], portstring)
		//untags := strings.Trim(o["ifaces"], " ")
		// get 'TG:...'
		//o = p.ParseSingle(regexps["tags"], portstring)
		//tags := strings.Trim(o["ifaces"], " ")

		vlan := dproto.Vlan{
			ID:vid,
			AccessPorts:make([]string,0),
			TrunkPorts:make([]string,0),
		}

		// parse port strings
		// string like `XGE0/0/11(D)    XGE0/0/16(D)    XGE0/0/17(D)    XGE0/0/18(D)   XGE0/0/19(D)`
		out := p.ParseMultiple(regexps["ifname"], untags)
		//fmt.Printf("%+v\n", out)
		for _, port := range out {
			ifname := strings.Trim(port["ifname"], " ")
			ifname = p.ConvertIfname(ifname)
			iftype := p.GetInterfaceType(ifname)
			if iftype == dproto.InterfaceType_UNKNOWN {
				p.Log("Warning! Unknown interface type (%s)", ifname)
				continue
			}
			vlan.AccessPorts = append(vlan.AccessPorts, ifname)
		}

		out = p.ParseMultiple(regexps["ifname"], tags)
		for _, port := range out {
			ifname := strings.Trim(port["ifname"], " ")
			ifname = p.ConvertIfname(ifname)
			iftype := p.GetInterfaceType(ifname)
			if iftype == dproto.InterfaceType_UNKNOWN {
				p.Log("Warning! Unknown interface type (%s)", ifname)
				continue
			}
			vlan.TrunkPorts = append(vlan.TrunkPorts, ifname)
		}
		vlans = append(vlans, &vlan)
	}

	return vlans, nil
}
