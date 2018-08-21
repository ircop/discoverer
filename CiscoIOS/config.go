package CiscoIOS

import (
	"fmt"
	"strings"
)

// GetConfig for CiscoIOS
func (p *Profile) GetConfig() (string, error) {
	p.Log("Starting CiscoIOS.GetConfig()")

	result, err := p.Cli.Cmd("show running-config")
	if err != nil {
		return "", fmt.Errorf("Cannot get config: %s", strings.Replace(err.Error(), "%", "%%", -1))
	}


	return result, nil
}
