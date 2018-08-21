package CiscoIOS

import (
	"fmt"
	"strings"
	"regexp"
)

// return map[poName][]slaveName
func (p *Profile) GetPortchannels() (map[string][]string, error) {
	p.Debug("Starting CiscoIOS.GetPortchannels()")

	pos := make(map[string][]string)

	result, err := p.Cli.Cmd("show etherchannel summary")
	if err != nil {
		return pos, fmt.Errorf("Cannot get etherchannel info: %s", err.Error())
	}
	p.Debug(result)

	patterns := make(map[string]string)
	patterns["split"] = `\d+\s+Po`
	patterns["ifstring"] = `(?msi:^(?P<name>[^\(]+)(?P<flags>\(([^\)]+)\))?\s+(?P<proto>([^\s]+))\s+(?P<ifstring>(.*)))`
	patterns["po"] = `(?mi:^(?P<group_id>\d+)\s+(?P<name>[^\(]+)(?P<flags>\(([^\)]+)\))?\s+(?P<proto>([^\s]+))\s+(?P<ifstring>(.+)(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?(\n\s+(.+))?))`
	patterns["ifname"] = `(^|\s)(?P<if>\S+)\(`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return pos, err
	}

	parts := regexps["split"].Split(result, -1)
	for _, part := range parts {
		if match, _ := regexp.Match(`^\d`, []byte(part)); !match {
			continue
		}
		part = "Po" + part
		out := p.ParseSingle(regexps["ifstring"], part)
		name := strings.Trim(out["name"], " ")
		ifstring := strings.Trim(out["ifstring"], " ")
		ifstring = strings.Replace(ifstring, "\n", "", -1)
		//if name == "" || ifstring == "" {
		if name == "" {
			p.Log("Cannot find all required PO parameters (name/ifstring: %s/%s)", name, ifstring)
			continue
		}

		ifnames := make([]string, 0)
		if ifstring != "" {
			ifs := p.ParseMultiple(regexps["ifname"], ifstring)
			for _, child := range ifs {
				if ifname, ok := child["if"]; ok {
					ifnames = append(ifnames, ifname)
				}
			}
		}
		pos[name] = ifnames
	}
	//fmt.Printf("POS: \n %+v\n", pos)

	/*strs := p.ParseMultiple(regexps["po"], result)
	for _, str := range strs {
		poName := str["name"]
		proto := str["proto"]
		ifString := str["ifstring"]

		if poName == "" || proto == "" {
			p.Log("Cannot find all required PO parameters (name/proto/ifs: %s/%s/%s)", poName, proto, ifString)
			continue
		}

		ifsOut := p.ParseMultiple(regexps["ifname"], ifString)
		ifnames := make([]string,0)
		for _, ifNamesResult := range ifsOut {
			if ifname, ok := ifNamesResult["if"]; ok {
				ifnames = append(ifnames, ifname)
			}
		}

		pos[poName] = ifnames
	}
	fmt.Printf("POS: \n %+v\n", pos)*/

	return pos, nil
}
