package HuaweiSW

import (
	"fmt"
	"github.com/ircop/dproto"
	"regexp"
	"strings"
)

// GetInterfaces for HuaweiSW profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Debug("Starting HuaweiSW.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	patterns := make(map[string]string, 0)
	patterns["ifname"] = `(?m)^(?P<ifname>[^\s]+) current state`
	patterns["desc"] = `(?m)^Description:(?P<desc>[^\n]+)`
	patterns["trunk"] = `(?ms)PortName\s+Status\s+Weight(\s+)?\n-+(\s+)?\n(?P<ports>.+)\n-+\n`
	patterns["members"] = `(?m)^(?P<ifname>[^\s]+)\s+[^\s]+\s+\d+`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	result, err := p.Cli.Cmd("display interface")
	if err != nil {
		return interfaces, fmt.Errorf("Cannot 'display interface': %s", err.Error())
	}
	//p.Debug(result)

	parts, err := p.getOutputParts(result)
	if err != nil {
		return interfaces, fmt.Errorf("Error parsing 'show interface' output: %s", err.Error())
	}

	for _, part := range parts {
		out := p.ParseSingle(regexps["ifname"], part)
		ifname := strings.Trim(out["ifname"], " ")
		out = p.ParseSingle(regexps["desc"], part)
		desc := strings.Trim(out["desc"], " ")
		if ifname == "" {
			p.Log("Warning! Empty ifname: '%s'", part)
			continue
		}
		short := p.ConvertIfname(ifname)
		intType := p.GetInterfaceType(ifname)
		if intType == dproto.InterfaceType_UNKNOWN {
			p.Log("Warning! Unknown interface type (%s/%s)", ifname, short)
			continue
		}

		iface := dproto.Interface{
			Type:intType,
			Shortname:short,
			Name:ifname,
			Description:desc,
			LldpID:ifname,
		}

		// if this is port-channel, parse member ports
		if intType == dproto.InterfaceType_AGGREGATED {
			members := make([]string,0)
			out := p.ParseSingle(regexps["trunk"], part)
			portinfo := strings.Trim(out["ports"], " ")
			mout := p.ParseMultiple(regexps["members"], portinfo)
			for _, memberPart := range mout {
				member := strings.Trim(memberPart["ifname"], " ")
				members = append(members, member)
			}
			iface.PoMembers = members
		}

		interfaces[ifname] = &iface
	}

	return interfaces, nil
}

// Since stupid golang regexps does not have negative lookahead, we should parse output manually =\
func (p *Profile) getOutputParts(output string) ([]string, error) {
	result := make([]string, 0)

	reFirst, err := regexp.Compile(`^[^\n]+$`)
	reSecond, err2 := regexp.Compile(`^Line protocol current state`)
	if err != nil {
		return result, err
	}
	if err2 != nil {
		return result, err2
	}

	lines := strings.Split(output, "\n")
	curPart := make([]string,0)
	prev := ""
	for _, line := range lines {
		// 1. Part start is 2 lines like `"*" current state.+\n` + 'Line protocol current state'
		if reSecond.Match([]byte(line)) && reFirst.Match([]byte(prev)) {
			// we have new part.
			curPart = curPart[:len(curPart)-1]
			result = append(result, strings.Join(curPart, "\n"))
			curPart = []string{prev, line}
			continue
		}
		// this is not new-line. Append this to current part ; set 'prev' to current line
		curPart = append(curPart, line)
		prev = line
	}

	// and add last part
	result = append(result, strings.Join(curPart, "\n"))

	return result, nil
}

// ConvertIfname takes interface full name and returns shortname
func (p *Profile) ConvertIfname(fullname string) string {

	if strings.HasPrefix(fullname, "40GE") || strings.HasPrefix(fullname, "Eth-Trunk") || strings.HasPrefix(fullname, "Vlanif") {
		return fullname
	}
	if strings.HasPrefix(fullname, "XGigabit") {
		return strings.Replace(fullname, "XGigabitEthernet", "xg", -1)
	}
	if strings.HasPrefix(fullname, "XGE") {
		return strings.Replace(fullname, "XGE", "xg", -1)
	}
	if strings.HasPrefix(fullname, "GigabitEthernet") {
		return strings.Replace(fullname, "GigabitEthernet", "g", -1)
	}
	if strings.HasPrefix(fullname, "GE") {
		return strings.Replace(fullname, "GE", "g", -1)
	}
	if strings.HasPrefix(fullname, "LoopBack") {
		return strings.Replace(fullname, "LoopBack", "Lo", -1)
	}

	// MEth, null, etc.
	return fullname
}
