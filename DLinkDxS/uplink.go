package DLinkDxS

import (
	"fmt"
	"net"
	"github.com/ircop/discoverer/util/mac"
	"strings"
)

// GetUplink for DLinkDxS
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting DLinkDxS.GetUplink()")
	// 1: get default route
	// 2: get arp of this route
	// 3: get port with this mac in FDB

	patterns := make(map[string]string)
	patterns["iproute"] = `0.0.0.0\s+(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)`
	patterns["arp"] = `[a-zA-Z0-9]+\s+(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))`
	patterns["port"] = `(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+(?P<port>[^\s]+)\s+Dynamic`

	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return "", err
	}

	//result, err := p.Cli.Cmd("sh iproute 0.0.0.0/0")
	result, err := p.Cli.Cmd("show iproute")
	if err != nil {
		return "", fmt.Errorf("Cannot get default iproute: %s", err.Error())
	}
	p.Debug(result)

	// get iproute
	out := p.ParseSingle(regexps["iproute"], result)
	gw := net.ParseIP(out["ip"])
	if gw == nil {
		return "", fmt.Errorf("Cannot parse default iproute")
	}

	// get arp for this ip
	result, err = p.Cli.Cmd(fmt.Sprintf("sh arpe ipa %s", gw))
	if err != nil {
		return "", fmt.Errorf("Cannot get arpentry for gw '%s'", gw)
	}
	p.Debug(result)
	out = p.ParseSingle(regexps["arp"], result)
	mac := out["mac"]
	if !Mac.IsMac(mac) {
		return "", fmt.Errorf("Cannot parse gefault gw's macaddr")
	}

	// get port by macaddr from FDB
	result, err = p.Cli.Cmd(fmt.Sprintf("show fdb mac %s", mac))
	if gw == nil {
		return "", fmt.Errorf("Cannot get fdb for gw mac '%s'", mac)
	}
	p.Debug(result)
	out = p.ParseSingle(regexps["port"], result)
	port := strings.Trim(out["port"], " ")
	if port == "" {
		return "", fmt.Errorf("Cannot parse gw macaddr port (mac: '%s')", mac)
	}


	return port, nil
}
