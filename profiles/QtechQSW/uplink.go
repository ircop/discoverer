package QtechQSW

import (
	"fmt"
	"strings"
)

// GetUplink for QtechQSW
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting QtechQSW.GetUplink()")

	patterns := make(map[string]string, 0)
	patterns["route"] = `(?msi:\* (?P<ip>[^,]+), via)`
	patterns["arp"] = `(?msi:(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+(?P<if>[^\s]+)\s+(?P<port>[^\s]+)\s+)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return "", err
	}

	result, err := p.Cli.Cmd("sh ip route 0.0.0.0")
	if err != nil {
		return "", err
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["route"], result)
	gw := strings.Trim(out["ip"], " ")
	if gw == "" {
		return "", nil
	}

	result, err = p.Cli.Cmd(fmt.Sprintf("show arp %s", gw))
	if err != nil {
		return "", err
	}
	out = p.ParseSingle(regexps["arp"], result)
	port := strings.Trim(out["port"], " ")

	return port, nil
}