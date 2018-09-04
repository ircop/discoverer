package DLinkDxS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"regexp"
	"strings"
)

// GetInterfaces for DLinkDxS profile
func (p *Profile) GetInterfaces() (map[string]dproto.Interface, error) {
	p.Debug("Starting DLinkDxS.GetInterfaces()")

	interfaces := make(map[string]dproto.Interface)

	result, err := p.Cli.Cmd("show ports description")
	if err != nil {
		p.Debug(result)
		return interfaces, fmt.Errorf("Error getting interfaces: %s", err.Error())
	}
	p.Debug(result)

	rePortsStr := `(?m:^\s*(?P<port>\d+(/|:)?\d*)\s*(\((?P<media_type>(C|F))\))?\s+(?P<admin_state>Enabled|Disabled)\s+(?P<admin_speed>Auto|10M|100M|1000M|10G)/((?P<admin_duplex>Half|Full)/)?(?P<admin_flowctrl>Enabled|Disabled)\s+(?P<status>LinkDown|Link\sDown|(?:Err|Loop)\-Disabled|Empty)?((?P<speed>10M|100M|1000M|10G)/(?P<duplex>Half|Full)/(?P<flowctrl>None|Disabled|802.3x))?\s+(?P<addr_learning>Enabled|Disabled)\s*((?P<trap_state>Enabled|Disabled)\s*)?((?P<asd>\-)\s*)?(\n\s+(?P<mdix>Auto|MDI|MDIX|Cross|Normal|\-)\s*)?(\n\s*Desc(ription)?:\s*?(?P<desc>.*?))?$)`
	rePorts, err := regexp.Compile(rePortsStr)
	if err != nil {
		return interfaces, err
	}

	// get lacp local info
	portIds := p.lldpLocalPorts()

	// get ports
	ports := p.ParseMultiple(rePorts, result)
	for _, port := range ports {
		name := strings.Trim(port["port"], " ")

		if name == "" {
			p.Log("Something wrong: Empty port name: %+v", port)
			continue
		}

		// todo: states and speeds will be collected in periodic discoveries, not here
		newInt := dproto.Interface{
			Name: strings.Trim(port["port"], " "),
			Shortname: strings.Trim(port["port"], " "),
			Description: strings.Trim(port["desc"], " "),
			Type: dproto.InterfaceType_PHISYCAL,
			LldpID: name,
		}
		if id, ok := portIds[name]; ok {
			newInt.LldpID = id
		}

		interfaces[name] = newInt
	}

	portchannels := p.getPortchannels()
	for name, portstring := range portchannels {
		portMembers := p.ExpandInterfaceRange(portstring)
		for _, n := range portMembers {
			if _, ok := interfaces[n]; !ok {
				p.Log("Something wrong: port-channel member '%s' doesnt exist in interface list.", n)
				continue
			}
		}

		newInt := dproto.Interface{
			Name: name,
			Shortname: name,
			Type: dproto.InterfaceType_AGGREGATED,
			PoMembers: portMembers,
		}
		interfaces[name] = newInt
	}

	return interfaces, nil
}

// Port-channel information
func (p *Profile) getPortchannels() map[string]string {
	p.Debug("Sending 'show link_aggregation'")

	// return: map[name(id)][]string
	result := make(map[string]string)

	out, err := p.Cli.Cmd("show link_aggregation")
	if nil != err {
		p.Log(err.Error())
		return result
	}
	p.Debug(out)

	re, err := regexp.Compile(`(?mis:Group ID\s+:\s+(T)?(?P<name>\d+).+?Type\s+:\s+(?P<type>\S+).+?Member Port\s+:(?P<members>[^\n]+\S+)?.+?Status\s+:\s+(?P<status>\S+))`)
	if err != nil {
		p.Log(err.Error())
		return result
	}

	groups := p.ParseMultiple(re, out)
	for _, group := range groups {
		poName, ok1 := group["name"]
		portstring := group["members"]

		if !ok1 {
			p.Log("Something wrong: matched port-challel without name")
			continue
		}
		poName = fmt.Sprintf("T%s", strings.Trim(poName, " "))
		portstring = strings.Trim(portstring, " ")
		result[poName] = portstring
	}

	return result
}

// lldpLocalPorts gathers all local portIDs
func (p *Profile) lldpLocalPorts() map[string]string {
	p.Debug("Sending 'show lldp local_ports'")

	data := make(map[string]string)

	result, err := p.Cli.Cmd("show lldp local_ports")
	if err != nil {
		p.Log("Cannot get lldp local ports: %s", err.Error())
		return data
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?mis:Port ID\s+:\s+(?P<port>\d+(?:[:/]\d+)?)\s*\n\-+\s*\nPort id subtype \s+:[^\n]+\nport id\s+:\s+?(?P<id>[^\n]+))`)
	if err != nil {
		p.Log("Cannot compile lldp local_ports regex")
		return data
	}

	lldp := p.ParseMultiple(re, result)
	p.Debug("Parsed %d lldp local entries", len(lldp))

	for _, l := range lldp {
		name := strings.Trim(l["port"], " ")
		id := strings.Trim(l["id"], " ")
		if name == "" || id == "" {
			continue
		}

		data[name] = id
	}

	return data
}
