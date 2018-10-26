package QtechQSW

import (
	"fmt"
	"regexp"
	"strings"
)

// GetConfig for QtechQSW
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting QtechQSW.GetConfig()")

	result, err := p.Cli.Cmd("show running-config")
	if err != nil {
		return "", err
	}

	reConfig, err := regexp.Compile(`(?ms:running-config\n(?P<config>.+end))`)
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