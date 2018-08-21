package EltexMES

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"strings"
)

// GetInterfaces for EltexMES profile
func (p *Profile) GetInterfaces() (map[string]discoverer.Interface, error) {
	p.Debug("Starting EltexMES.GetInterfaces()")
	interfaces := make(map[string]discoverer.Interface)

	patterns := make(map[string]string)
	patterns["ifs"] = `(?ms:^(?P<ifname>[a-zA-Z]\S+)\s+(?P<oper>Up|Down)\s+(?P<admin>Up|Down|Not Present)\s(?:(?P<desc>.*?)?)?)`
	patterns["members"] = `(?msi:^(?P<ifname>Po\d+)\s+(?P<type1>\S+):\s+(?P<interfaces1>\S+)+(\s+(?P<type2>\S+):\s+(?P<interfaces2>\S+)$|$))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	result, err := p.Cli.Cmd("show interfaces description")
	if err != nil {
		return interfaces, fmt.Errorf("Unable to execute 'show int descr': %s", err.Error())
	}
	p.Debug(result)

	poDescriptions := make(map[string]string)
	out := p.ParseMultiple(regexps["ifs"], result)
	for _, part := range out {
		ifname := strings.Trim(part["ifname"], " ")
		desc := strings.Trim(part["desc"], " ")
		if ifname != "" {
			if strings.Contains(ifname, "Po") {
				poDescriptions[ifname] = desc
				continue
			}
			iface := discoverer.Interface{
				Name:ifname,
				Description:desc,
				Type:discoverer.IntTypePhisycal,
				Shortname:ifname,
				LldpID:ifname,
			}
			interfaces[ifname] = iface
		}
	}

	/////////// port-channels ////////////
	result, err = p.Cli.Cmd("show interfaces port-channel")
	if err != nil {
		result, err = p.Cli.Cmd("show interfaces channel-group")
		if err != nil {
			return interfaces, fmt.Errorf("Unable to get port-channel info: %s", err.Error())
		}
	}
	p.Debug(result)

	// parse
	out = p.ParseMultiple(regexps["members"], result)
	for _, part := range out {
		ifname := strings.Trim(part["ifname"], " ")
		ifaces := strings.Trim(part["interfaces1"], " ")
		if ifaces == "" {
			ifaces = strings.Trim(part["interfaces2"], " ")
		}
		if ifaces == "" || ifname == "" {
			continue
		}
		iflist := p.ExpandInterfaceRange(ifaces)
		iface := discoverer.Interface{
			Name:ifname,
			Description: poDescriptions[ifname],
			LldpID:ifname,
			Shortname:ifname,
			Type:discoverer.IntTypeAggregated,
			PoMembers:iflist,
		}
		interfaces[ifname] = iface
	}

	return interfaces, nil
}
