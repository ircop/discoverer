package HuaweiSW

import (
	"regexp"
	"strings"
)

// GetConfig for HuaweiSW
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting HuaweiSW.GetConfig()")

	result, err := p.Cli.Cmd("display current-configuration")
	if err != nil {
		return "", err
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?msi:display current-configuration\n(?P<cfg>.*)return)`)
	if err != nil {
		return "", err
	}

	out := p.ParseSingle(re, result)
	cfg := strings.Trim(out["cfg"], " ") + "return\n"

	return cfg, nil
}