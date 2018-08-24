package DLinkDGS3100

import (
	"fmt"
	"net"
	"github.com/ircop/discoverer/util/mac"
	"strings"
)

// GetUplink for DLinkDGS3100
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting DLinkDGS3100.GetUplink()")

	patterns := make(map[string]string)
	patterns["route"] = `0.0.0.0\s+(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+`
	patterns["mac"] = `(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+dynamic`
	patterns["port"] = `(?m:^\d+\s+[^\s]+\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+(?P<port>\d:\d+)\s+)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return "", fmt.Errorf("Cannot compile regexps: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show iproute")
	if err != nil {
		return "", fmt.Errorf("Cannot 'show iproute': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["route"], result)
	gw := net.ParseIP(out["gw"])
	if gw == nil {
		return "", fmt.Errorf("Cannot parse gateway (%s)", out["gw"])
	}

	result, err = p.Cli.Cmd("show arpentry ipaddress " + gw.String())
	if err != nil {
		return "", fmt.Errorf("Cannot 'show arpentry ipaddress %s': %s", gw.String(), err.Error())
	}
	p.Debug(result)
	out = p.ParseSingle(regexps["mac"], result)
	mac := out["mac"]
	if !Mac.IsMac(mac) {
		return "", fmt.Errorf("Cannot parse gw arpentry (%s)", gw.String())
	}

	result, err = p.Cli.Cmd("show fdb mac " + mac)
	if err != nil {
		return "", fmt.Errorf("Cannot run 'sh fdb mac %s': %s", mac, err.Error())
	}
	p.Debug(result)


	out = p.ParseSingle(regexps["port"], result)
	port := strings.Trim(out["port"], " ")
	if port == "" {
		return "", fmt.Errorf("Cannot parse uplink port")
	}

	return port, nil
}
