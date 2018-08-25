package JunOS

import (
	"fmt"
	"regexp"
	"strings"
)

// GetUplink for JunOS
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting JunOS.GetUplink()")

	reVia, err := regexp.Compile(`(?m:\svia\s(?P<ifname>[^\s\n]+))`)
	if err != nil {
		return "", fmt.Errorf("Failed to compile route regex: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show route 0.0.0.0")
	if err != nil {
		return "", fmt.Errorf("Failed to 'show route 0.0.0.0': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseMultiple(reVia, result)
	if len(out) != 1 {
		p.Log("Cannot find uplink: routes count != 1 (%d)", len(out))
		return "", nil
	}
	part := out[0]
	ifname := strings.Trim(part["ifname"], " ")

	if ifname != "" {
		return ifname, nil
	}

	return "", nil
}
