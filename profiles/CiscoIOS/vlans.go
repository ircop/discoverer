package CiscoIOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"regexp"
	"strconv"
	"strings"
)

// GetVlans for CiscoIOS
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting CiscoIOS.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)
	vlanmap := make(map[string]dproto.Vlan)

	// here we need filled model/version
	if p.Model == "" {
		return vlans, fmt.Errorf("Model field is empty, you should run GetPlatform() first.")
	}

	if match, err := regexp.Match(`ASR100[0-6]`, []byte(p.Model)); match && err == nil {
		// router does not have vlans
		// update: todo: check dot1q encapsulation on interfaces
		return p.getRouterVlans()
		//return vlans, nil
	}


	patterns := make(map[string]string)
	patterns["vlanids"] = `(?P<vid>\d+)\s+(?P<name>[^\s]+)\s+active`
	//patterns["switchport"] = `(?msi:^Name:\s+(?P<ifname>[^\n]+)\n.+\nAdministr.+\nOperational Mode:\s+(?P<mode>[^\s]+).+\n(Admi.+\n)?(Oper.+\n)?(Nego.+\n)?Access Mode VLAN:\s+(?P<access_vlan>\d+).+\nTrunking Native Mode VLAN:\s+(?P<native_vlan>\d+).+\n(.*)Trunking VLANS Enabled:\s+(?P<trunking_vlans>.+))Pruning`
	patterns["switchport"] = `(?msi:^Name:\s+(?P<ifname>[^\n]+)\n.+\nAdministrative Mode:\s+(?P<mode>[^\s]+)\nOperational Mode:\s+(?P<mode2>[^\s]+).+\n(Admi.+\n)?(Oper.+\n)?(Nego.+\n)?Access Mode VLAN:\s+(?P<access_vlan>\d+).+\nTrunking Native Mode VLAN:\s+(?P<native_vlan>\d+).+\n(.*)Trunking VLANS Enabled:\s+(?P<trunking_vlans>.+))Pruning`
	patterns["ifname"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return vlans, err
	}

	// First: get vlan list
	// Second: get switchport, parse vlans, compare
	result, err := p.Cli.Cmd("show vlan")
	if err != nil {
		return vlans, fmt.Errorf("Error 'show vlan': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseMultiple(regexps["vlanids"], result)
	//fmt.Printf("PARTS: %d\n", len(out))
	for _, part := range out {
		name := strings.Trim(part["name"], " ")
		vidStr := strings.Trim(part["vid"], " ")
		if name == "" || vidStr == "" {
			continue
		}
		vid, err := strconv.ParseInt(vidStr, 10, 60)
		if err != nil {
			p.Log("Error: cannot parse vlan id '%s'", vidStr)
			continue
		}

		// todo: make map[vid]Vlan
		vlan := dproto.Vlan{
			TrunkPorts: make([]string, 0),
			AccessPorts: make([]string, 0),
			ID: vid,
			Name: vidStr,
		}
		vlanmap[vidStr] = vlan
	}
	fmt.Printf("VLANMAP LEN %d\n", len(vlanmap))

	// Got vlan IDs. Now parse switchports.
	result, err = p.Cli.Cmd("show int switchport")
	if err != nil {
		return vlans, fmt.Errorf("Cannot run 'show int switchport': %s", err.Error())
	}
	p.Debug(result)
	//parts := regexps["split"].Split(result, -1)
	parts := strings.Split(result, "Name:")
	//fmt.Printf("PARTS: %d\n", len(parts))
	for _, part := range parts {
		part = "Name:" + part
		out := p.ParseSingle(regexps["switchport"], part)
		ifname := strings.Trim(out["ifname"], " ")
		mode := strings.Trim(out["mode"], " ")
		if ifname == "" || mode == "" {
			continue
		}

		ifname, err = p.ConvertIfname(ifname, regexps["ifname"])
		if err != nil {
			p.Log("Error: cannot convert cisco port name '%s'", ifname)
		}
		if mode != "access" && mode != "trunk" {
			p.Log("Port '%s' is not trunk nor access (%s)", ifname, mode)
			continue
		}
		native := strings.Trim(out["native_vlan"], " ")
		access := strings.Trim(out["access_vlan"], " ")
		trunk := strings.Trim(out["trunking_vlans"], " ")

		if mode == "access" {
			// we need only access vlan
			if access == "" {
				continue
			}
			_, err := strconv.ParseInt(access, 10, 64)
			if err != nil {
				p.Log("Error: cannot parse vid '%s'", access)
				continue
			}

			vlan, ok := vlanmap[access]
			if !ok {
				continue
			}
			vlan.AccessPorts = append(vlan.AccessPorts, ifname)
			vlanmap[access] = vlan
			continue
		}

		if mode == "trunk" {
			// parse trunk vlans on each port and add them into vlan's trunk array
			// parse native vlans on each ports and add them into vlan's access array
			if native != "" {
				_, err := strconv.ParseInt(native, 10, 64)
				if err != nil {
					p.Log("Error: failed to parse native vlan (%s) for port '%s'", native, ifname)
				} else {
					// add vid to arr
					vlan, ok := vlanmap[native]
					if ok {
						vlan.AccessPorts = append(vlan.AccessPorts, ifname)
						vlanmap[native] = vlan
					}
				}
			}

			if trunk != "" {
				// parse trunking ports
				trunk = strings.Replace(trunk, "\n", "", -1)
				trunk = strings.Replace(trunk, " ", "", -1)
				if trunk == "ALL" {
					for k, v := range vlanmap {
						v.TrunkPorts = append(v.TrunkPorts, ifname)
						vlanmap[k] = v
					}
				}
				for _, v := range p.ExpandVlanRange(trunk) {
					if vlan, ok := vlanmap[v]; ok {
						vlan.TrunkPorts = append(vlan.TrunkPorts, ifname)
						vlanmap[v] = vlan
					}
				}
			}
		}
	}

	// vlanmap to vlan arr
	for _, v := range vlanmap {
		vlans = append(vlans, &v)
	}

	return vlans, nil
}

// Find vlans on ASR*
func (p *Profile) getRouterVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting CiscoIOS.getRouterVlans()")
	vlans := make([]*dproto.Vlan, 0)
	vlanmap := make(map[string]dproto.Vlan)

	patterns := make(map[string]string)
	patterns["ports"] = `(?m:^(\s+)?(?P<ifname>.+?)\s+is(?:\s+administratively)?\s+`+
		`(?P<admin>up|down),\s+line\s+protocol\s+is\s+`+
		`(?P<oper>up|down)(\s+)?(?:\((?:connected|notconnect|disabled|monitoring|err-disabled)\)\s*)?\n\s+`+
		`(.*)address is (?P<mac>([0-9A-Fa-f]){4}\.([0-9A-Fa-f]){4}\.([0-9A-Fa-f]){4})(.*)\n`+
		`(?:\s+Description:\s(?P<desc>[^\n]+)\n)?(?:\s+Internet address ((is\s(?P<ip>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}))|([^\d]+))(\s+)?\n)?[^\n]+\n[^\n]+\n\s+`+
		`Encapsulation\s+(?P<encaps>[^\n]+))`
	patterns["ifname"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	patterns["encaps"] = `802\.1Q\s+Virtual\s+LAN,\s+Vlan\s+ID\s+(?P<vid>\d+)\.`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return vlans, err
	}

	result, err := p.Cli.Cmd("show interfaces")
	if err != nil {
		panic(err)
	}
	p.Debug(result)

	out := p.ParseMultiple(regexps["ports"], result)
	for _, part := range out {
		ifname := strings.Trim(part["ifname"], " ")
		encaps := strings.Trim(part["encaps"], " ")
		if ifname == "" || encaps == "" {
			continue
		}

		ifname, err = p.ConvertIfname(ifname, regexps["ifname"])
		if err != nil {
			p.Log("Cannot convert ifname '%s': %s", ifname, err.Error())
		}

		o := p.ParseSingle(regexps["encaps"], encaps)
		vid := strings.Trim(o["vid"], " ")
		if vid == "" {
			//fmt.Printf("VID EMPT; o = %+v, ENCAPS = '%s'\n", o, encaps)
			continue
		}

		vidInt, err := strconv.ParseInt(vid, 10, 64)
		if err != nil {
			p.Log("Cannot convert vlan id '%s' to integer", vid)
			continue
		}

		if v, ok := vlanmap[vid]; ok {
			v.TrunkPorts = append(v.TrunkPorts, ifname)
			vlanmap[vid] = v
		} else {
			v := dproto.Vlan{
				TrunkPorts: []string{ifname},
				Name:ifname,
				ID:vidInt,
			}
			vlanmap[vid] = v
		}
	}

	for _, v := range vlanmap {
		vlans = append(vlans, &v)
	}

	return vlans, nil
}