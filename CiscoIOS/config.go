package CiscoIOS

import (
	"fmt"
	"strings"
	"regexp"
)

// GetConfig for CiscoIOS
func (p *Profile) GetConfig() (string, error) {
	p.Log("Starting CiscoIOS.GetConfig()")

	result, err := p.Cli.Cmd("show running-config")
	if err != nil {
		return "", fmt.Errorf("Cannot get config: %s", strings.Replace(err.Error(), "%", "%%", -1))
	}
	p.Debug(result)

	reConfig, err := regexp.Compile(`(?ms:Current configuration(\s+)?:\s+\d+\s+bytes\n(?P<config>.+end))`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile config regex: %s", err.Error())
	}

	out := p.ParseSingle(reConfig, result)
	config := strings.Trim(out["config"], " ")
	if config == "" {
		return "", fmt.Errorf("Cannot parse config")
	}


	return config, nil
}
