package RouterOS

import (
	"fmt"
	"regexp"
	"strings"
)

// GetUplink for RouterOS
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting RouterOS.GetUplink()")

	result, err := p.Cli.Cmd(`/ip route print detail without-paging  where dst-address="0.0.0.0/0"`)
	if err != nil {
		return "", err
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?ms:reachable\s+via\s+(?P<ifname>[^\s\n]+))`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile uplink regex: %s", err.Error())
	}
	out := p.ParseSingle(re, result)
	iface := strings.Trim(out["ifname"], " ")

	return iface, nil
}
