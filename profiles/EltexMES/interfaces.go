package EltexMES

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"regexp"
	"strings"
)

// GetInterfaces for EltexMES profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Debug("Starting EltexMES.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

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
			iface := dproto.Interface{
				Name:ifname,
				Description:desc,
				Type:dproto.InterfaceType_PHISYCAL,
				Shortname:ifname,
				LldpID:ifname,
			}
			interfaces[ifname] = &iface
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
		iface := dproto.Interface{
			Name:ifname,
			Description: poDescriptions[ifname],
			LldpID:ifname,
			Shortname:ifname,
			Type:dproto.InterfaceType_AGGREGATED,
			PoMembers:iflist,
		}
		interfaces[ifname] = &iface
	}

	// And L3 interfaces (svi's)
	result, err = p.Cli.Cmd("sh ip interface")
	if err != nil {
		return interfaces, fmt.Errorf("Cannot 'sh ip interface': %s", err.Error())
	}
	p.Debug(result)

	// some models shows both IP and GW tables; GW first. Cut it
	reIfHeader, err := regexp.Compile(`IP Address\s+I\/F\s+`)
	if err != nil {
		return interfaces, fmt.Errorf("Cannot compile ip/gw split regex: %s", reIfHeader)
	}
	parts := reIfHeader.Split(result, -1)
	if len(parts) < 2 {
		p.Log("Cannot split ip/gw tables: split result: %d", len(parts))
	}
	result = parts[1]

	rows := text.ParseTable(result, "^--", "", false)
	reSvi, err := regexp.Compile(`vlan\s\d+`)
	if err != nil {
		return interfaces, fmt.Errorf("Cannot compile svi regex")
	}

	for _, row := range rows {
		if len(row) < 4 {
			continue
		}

		ifname := strings.Trim(row[1], " ")
		if !reSvi.Match([]byte(ifname)) {
			continue
		}
		ifname = strings.Replace(ifname, " ", "", -1)
		iface := dproto.Interface{
			Name: ifname,
			Shortname: ifname,
			Type:dproto.InterfaceType_AGGREGATED,
		}
		interfaces[ifname] = &iface
	}

	return interfaces, nil
}
