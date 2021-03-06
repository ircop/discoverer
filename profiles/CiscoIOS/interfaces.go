package CiscoIOS

import (
	"fmt"
	"github.com/ircop/dproto"
	"github.com/ircop/discoverer/util/text"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// GetInterfaces for CiscoIOS profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Log("Starting CiscoIOS.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	patterns := make(map[string]string)
	/*patterns["ports"] = `(?m:^(\s+)?(?P<ifname>.+?)\s+is(?:\s+administratively)?\s+`+
						`(?P<admin>up|down),\s+line\s+protocol\s+is\s+`+
						`(?P<oper>up|down)(\s+)?(?:\((?:connected|notconnect|disabled|monitoring|err-disabled)\)\s*)?\n\s+`+
						`(.*)address is (?P<mac>([0-9A-Fa-f]){4}\.([0-9A-Fa-f]){4}\.([0-9A-Fa-f]){4})(.*)\n`+
						`(?:\s+Description:\s(?P<desc>[^\n]+)\n)?(?:\s+Internet address ((is\s(?P<ip>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}))|([^\d]+))(\s+)?\n)?[^\n]+\n[^\n]+\n\s+`+
						`Encapsulation\s+(?P<encaps>[^\n]+))`*/
	patterns["ifname"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	patterns["name"] = `(?m:^(\s+)?(?P<ifname>.+?)\s+is(?:\s+administratively)?\s+(?P<admin>up|down),\s+line\s+protocol\s+is\s+(?P<oper>up|down)(\s+)?)`
	patterns["desc"] = `(?m:^\s+Description:\s+(?P<desc>[^\n]+)\n)`
	patterns["lldp"] = `(?mis:^(?P<iface>(?:Fa|Gi|Te)[^:]+?):)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	// First get portchannels
	pos, err := p.GetPortchannels()
	if err != nil {
		p.Log("Error getting portchannels: %s", err.Error())
		return interfaces, nil
	}
	//return interfaces, fmt.Errorf("test")

	// Then get interfaces
	result, err := p.Cli.Cmd("show interface")
	if err != nil {
		return interfaces, errors.Wrap(err, "Cannot get interfaces data")
	}
	p.Debug(result)

	parts, err := text.SplitByParts(result, `(?msi:^[^\s]+\s+is (administratively )?(up|down), line)`)
	if err != nil {
		return interfaces, errors.Wrap(err, "Cannot split output by parts")
	}

	for _, part := range parts {
		out := p.ParseSingle(regexps["name"], part)
		ifname := strings.Trim(out["ifname"], " ")
		if ifname == "" {
			p.Log("Error: empty interface name!")
			continue
		}

		shortname, err := p.ConvertIfname(ifname, regexps["ifname"])
		if err != nil {
			p.Log(err.Error())
		}

		out = p.ParseSingle(regexps["desc"], part)
		desc := strings.Trim(out["desc"], " ")

		iftype := p.GetInterfaceType(shortname)
		newInt := dproto.Interface{
			Name:ifname,
			Shortname:shortname,
			Type:iftype,
			LldpID:ifname,
			Description:desc,
		}
		if iftype == dproto.InterfaceType_AGGREGATED {
			if po, ok := pos[shortname]; ok {
				newInt.PoMembers = po
			} else {
				p.Log("WARNING! Cannot find port-channel details for '%s'!", shortname)
			}
		}
		/*if iftype == dproto.InterfaceType_TUNNEL {
			continue
		}
		p.Debug("%s", part)*/

		interfaces[ifname] = &newInt
	}

	/*
	ports := p.ParseMultiple(regexps["ports"], result)
	for _, port := range ports {
		ifname := port["ifname"]
		if ifname == "" {
			p.Log("Error: empty interface name!")
			continue
		}
		ifname, err = p.ConvertIfname(port["ifname"], regexps["ifname"])
		if err != nil {
			p.Log(err.Error())
		}

		iftype := p.GetInterfaceType(ifname)
		newInt := dproto.Interface{
			Name: strings.Trim(port["ifname"], " "),
			Shortname: ifname,
			Description: strings.Trim(port["desc"], " "),
			Type:iftype,
			LldpID:strings.Trim(port["ifname"], " "),
		}
		if iftype == dproto.InterfaceType_AGGREGATED {
			if po, ok := pos[ifname]; ok {
				newInt.PoMembers = po
			} else {
				p.Log("WARNING! Cannot find port-channel details for '%s'!", ifname)
			}
		}

		interfaces[ifname] = &newInt
	}
	*/

	return interfaces, nil
}


func (p *Profile) ConvertIfname(fullname string, re *regexp.Regexp) (string, error) {
	short := strings.Trim(fullname, " ")

	out := re.FindStringSubmatch(strings.ToLower(fullname))
	if len(out) < 3 {
		return short, fmt.Errorf("Failed to get interface short name: %s", fullname)
	}

	ifType := strings.Title(out[1])
	if ifType == "Et" {
		ifType = "Eth"
	}
	ifNum := out[2]
	if ifType == "" || ifNum == "" {
		return short, fmt.Errorf("Failed to get interface short name number: %s", fullname)
	}

	short = fmt.Sprintf("%s%s", ifType, ifNum)
	return short, nil
}
