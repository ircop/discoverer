package EltexMES

import (
	"fmt"
	"strings"
)

// GetConfig for EltexMES
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting EltexMES.GetConfig()")


	result, err := p.Cli.Cmd("show running-config")
	if err != nil {
		return "", fmt.Errorf("Cannot get config: %s", err.Error())
	}
	p.Debug(result)

	result = strings.Replace(result, "show running-config\n", "", -1)

	return result, nil
}
