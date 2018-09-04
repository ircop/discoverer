package CiscoIOS

import (
	"fmt"
	"strings"
	"net"
)

// GetUplink for CiscoIOS
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting CiscoIOS.GetUplink()")

	// 1: get default route
	// 2: get arp of this route
	// 3: get port with this mac in FDB

	patterns := make(map[string]string)
	patterns["gw"] = `\* (?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)`
	patterns["arp"] = `Internet\s+(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(\d+|-)\s+[^\s]+\s+ARPA\s+(?P<ifname>[^\n]+)\n`
	patterns["ifname"] = `^(?P<type>[a-z]{2})[a-z\-]*\s*(?P<number>\d+(/\d+(/\d+)?)?(\.\d+(/\d+)*(\.\d+)?)?(:\d+(\.\d+)*)?(/[a-z]+\d+(\.\d+)?)?(A|B)?)$`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return "", err
	}

	r, err := p.Cli.Cmd("sh ip route 0.0.0.0 0.0.0.0")
	if err != nil {
		return "", fmt.Errorf("Failed to get default route: %s", err.Error())
	}
	p.Debug(strings.Replace(r, "%", "%%", -1))

	out := p.ParseSingle(regexps["gw"], r)
	gw := strings.Trim(out["ip"], "")
	if net.ParseIP(gw) == nil {
		return "", fmt.Errorf("Cannot parse gw address (%s)", gw)
	}

	// place space after ip, to avoid unwanted matches like '10.10.10.1' => '10.10.10.11', '10.10.10.12', etc.
	r, err = p.Cli.Cmd(fmt.Sprintf("sh arp | in %s ", gw))
	if err != nil {
		return "", fmt.Errorf("Cannot get arp entry for gw: %s", err.Error())
	}
	p.Debug(strings.Replace(r, "%", "%%", -1))

	out = p.ParseSingle(regexps["arp"], r)
	ifname := strings.Trim(out["ifname"], " ")
	if ifname == "" {
		return "", fmt.Errorf("Cannot parse uplink ifname by arp record")
	}

	short, err := p.ConvertIfname(ifname, regexps["ifname"])
	if err != nil {
		return "", fmt.Errorf("Cannot convert ifname (%s) to short one", ifname, err.Error())
	}

	return short, nil
}
