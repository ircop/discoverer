package RouterOS

import (
	"github.com/ircop/dproto"
	"strings"
)

// GetInterfaces for RouterOS profile
func (p *Profile) GetInterfaces() (map[string]*dproto.Interface, error) {
	p.Debug("Starting RouterOS.GetInterfaces()")
	interfaces := make(map[string]*dproto.Interface)

	patterns := make(map[string]string)
	patterns["ports"] = `(?m:\s+\d+\s+[DXRS]+\s+(?P<comment>;;;[^\n]+\n\s+)?name="(?P<name>[^"]+)")\s+(default-name="[^"]+"\s+)?type="(?P<type>[^"]+)"`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return interfaces, err
	}

	result, err := p.Cli.Cmd("/interface print detail without-paging")
	if err != nil {
		return interfaces, err
	}
	p.Debug(result)

	parts := p.ParseMultiple(regexps["ports"], result)
	for _, part := range parts {
		intName := part["name"]
		intType := part["type"]
		comment := strings.Trim(part["comment"], ";;;")
		comment = strings.Replace(comment, "\n", " ", -1)
		comment = strings.Trim(comment, " ")
		if intName == "" || intType == "" {
			p.Log("Warning: cannot parse interface name/type (%s/%s)", intName, intType)
			continue
		}
		tp := p.GetInterfaceType(intType)
		if tp == dproto.InterfaceType_UNKNOWN {
			p.Log("Warning: unknown interface type ('%s')", intType)
		}
		iface := dproto.Interface{
			Description:comment,
			Name:intName,
			Shortname:intName,
			Type:tp,
		}
		interfaces[intName] = &iface
	}

	return interfaces, nil
}
