package JunOS

import (
	"fmt"
	"github.com/ircop/dproto"

	//"github.com/ircop/discoverer/util/text"
	"github.com/ircop/discoverer/util/text"
	"strings"
)

// GetInterfaces for JunOS profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Log("Starting JunOS.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	// run TERSE: get all active interfaces ; cut %.0 ; remember them
	// run DESCRIPTIONS: get descriptions for collected interfaces

	patterns := make(map[string]string, 0)
	patterns["po"] = `(?msi:(\s+)?Link:(\s+)?\n(?P<ifaces>.+)(\n+)\s+Aggregate)`
	patterns["ifstring"] = `(\s+)?(?P<iface>[^\s]+)`
	patterns["replace"] = `(?msi:^\s+(Input|Output)([^\n]+))`
	patterns["replace_unit"] = `(?msi:(\.\d+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	// -- terse --
	result, err := p.Cli.Cmd("show interfaces terse")
	if err != nil {
		return interfaces, fmt.Errorf("Cannot 'show interfaces terse': %s", err.Error())
	}
	p.Debug(result)

/*	rePo, err := regexp.Compile(`(?msi:(\s+)?Link:(\s+)?\n(?P<ifaces>.+)(\n+)\s+Aggregate)`)
	if err != nil {
		return interfaces, fmt.Errorf("Cannot compile port-channel regex: %s", err.Error())
	}
	reIfstring, err := regexp.Compile(`(\s+)?(?P<iface>[^\s]+)`)
	if err != nil {
		return interfaces, fmt.Errorf("Cannot compile ifstring regex: %s", err.Error())
	}*/

	rows := text.ParseTable(result, `Interface\s+`, "", true, false)
	for _, row := range rows {
		if len(row) < 4 {
			p.Log("Warning! Interfaces row len is %d", len(row))
			continue
		}
		ifname := strings.Trim(row[0], " ")
		if ifname == "" || ifname == "vlan" {
			continue
		}
		if strings.HasSuffix(ifname, ".0") {
			continue
		}
		iftype := p.GetInterfaceType(ifname)

		if iftype == dproto.InterfaceType_UNKNOWN {
			p.Debug("Unknown interface type: '%s'", ifname)
			continue
		}

		iface := dproto.Interface {
			Name:ifname,
			Shortname:ifname,
			LldpID:ifname,
			Description: "",
			Type:iftype,
		}

		if iftype == dproto.InterfaceType_AGGREGATED {
			p.Debug("show interfaces %s detail", ifname)
			r, e := p.Cli.Cmd("show interfaces "+ifname+" detail")
			if e != nil {
				p.Log("Error: cannot get interface %s details: %s", ifname, e.Error())
				continue
			}
			p.Debug(r)
			out := p.ParseSingle(regexps["po"], r)
			ifaces := strings.Trim(out["ifaces"], " ")
			ifaces = regexps["replace"].ReplaceAllString(ifaces, "")
			ifaces = regexps["replace_unit"].ReplaceAllString(ifaces, "")
			//ifaces = strings.Replace(ifaces, ".0", "", -1)
			ifaces = strings.Replace(ifaces, "\n", "", -1)
			ifaces = strings.Trim(ifaces, " ")

			if ifaces != "" {
				out2 := p.ParseMultiple(regexps["ifstring"], ifaces)
				members := make([]string, 0)
				for _, part := range out2 {
					iface := strings.Trim(part["iface"], " ")
					if iface != "" {
						members = append(members, iface)
					}
				}

				iface.PoMembers = members
			}
		}

		interfaces[ifname] = &iface
	}

	// -- descriptions --
	result, err = p.Cli.Cmd("show interfaces descriptions")
	if err != nil {
		return interfaces, fmt.Errorf("Cannot 'show interfaces descriptions': %s", err.Error())
	}
	p.Debug(result)
	rows = text.ParseTable(result, `Interface\s+`, "", true, false)
	for _, row := range rows {
		if len(row) < 4 {
			p.Log("Warning: row len < 4")
			continue
		}
		ifname := strings.Trim(row[0],  " ")
		if strings.HasSuffix(ifname, ".0") {
			continue
		}
		descr := strings.Trim(row[3], " ")
		if iface, ok := interfaces[ifname]; ok && descr != "" {
			iface.Description = descr
			interfaces[ifname] = iface
		}
	}

	return interfaces, nil
}
