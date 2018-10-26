package QtechQSW

import (
	"github.com/ircop/dproto"
	"strings"
)

// GetInterfaces for HuaweiSW profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Debug("Starting QtechQSW.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	patterns := make(map[string]string, 0)
	patterns["ifaces"] = `(?msi:^\s*(?P<ifname>\S+) is (?P<admin_status>\S*\s*\S+), line protocol is (?P<oper_status>\S+))`
	patterns["descr"] = `(?msi:^\s*(?P<ifname>\S+) is (?P<admin_status>\S*\s*\S+), line protocol is (?P<oper_status>\S+)\n(\s+[^\n]+alias name is (?P<descr>[^,]+), index)?)`
	patterns["po"] = `(?msi:^\s*(?P<ifname>\S+) is LAG member port, LAG port:(\s+)?(?P<po>[^\n]+)\n)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	result, err := p.Cli.Cmd("show interface")
	if err != nil {
		return interfaces, err
	}
	//p.Debug(result)

	out := p.ParseMultiple(regexps["ifaces"], result)
	for i := range out {
		//p.Debug("%+#v", out[i]["ifname"])
		ifname := strings.Trim(out[i]["ifname"], " ")
		if ifname == "" {
			continue
		}

		t := p.GetInterfaceType(ifname)


		iface := dproto.Interface{
			Shortname:ifname,
			Name:ifname,
			Type:t,
			LldpID:ifname,
		}
		if t == dproto.InterfaceType_AGGREGATED {
			iface.PoMembers = make([]string, 0)
		}

		// todo: po members (if po)

		interfaces[ifname] = &iface
	}

	// parse descriptions
	// separate regex, because vlanifs has no descr
	out = p.ParseMultiple(regexps["descr"], result)
	for i := range out {
		ifname := strings.Trim(out[i]["ifname"], " ")
		if ifname == "" {
			continue
		}

		d := strings.Trim(out[i]["descr"], " ")
		if d == "(null)" || d == "" {
			continue
		}

		iface, ok := interfaces[ifname]
		if ok {
			iface.Description = d
			interfaces[ifname] = iface
		}
	}

	// scan port-channel members
	out = p.ParseMultiple(regexps["po"], result)
	for i := range out {
		ifname := strings.Trim(out[i]["ifname"], " ")
		if ifname == "" {
			continue
		}

		po := strings.Trim(out[i]["po"], " ")
		if po == "" {
			continue
		}

		// check if and po existance
		pchan, ok := interfaces[po]
		if !ok {
			p.Log("Error: no port-channel found: '%s'", po)
			continue
		}

		if _, ok := interfaces[ifname]; !ok {
			p.Log("Error: no port-channel member found: '%s'", ifname)
			continue
		}

		pchan.PoMembers = append(pchan.PoMembers, ifname)
		interfaces[po] = pchan
	}


	return interfaces, nil
}
