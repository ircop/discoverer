package JunOS

import (
	"fmt"
	"regexp"
	"strings"
)

// GetConfig for JunOS
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting JunOS.GetConfig()")

	reConfig, err := regexp.Compile(`(?ms:## Last[^\n]+\n(?P<config>.+)\{master)`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile config regex: %s", err.Error())
	}

	result, err := p.Cli.Cmd("conf")
	if err != nil {
		return "", fmt.Errorf("Cannot enter configuration mode: %s", err.Error())
	}
	p.Debug(result)

	result, err = p.Cli.Cmd("show")
	if err != nil {
		return "", fmt.Errorf("Cannot 'show': %s", err.Error())
	}

	out := p.ParseSingle(reConfig, result)
	config := strings.Trim(out["config"], " ")
	if config == "" {
		return "", fmt.Errorf("Cannot parse config")
	}

	p.Cli.Cmd("q")

	return config, nil
}