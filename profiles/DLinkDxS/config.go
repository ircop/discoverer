package DLinkDxS

import (
	"fmt"
	"regexp"
	"strings"
)

// GetConfig for DLinkDxS
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting DLinkDxS.GetConfig()")


	result, err := p.Cli.Cmd("show config current_config")
	//result, err := p.Cli.Cmd("show switch")
	if err != nil {
		return "", fmt.Errorf("Failed to get config: %s", err.Error())
	}
	if strings.Contains(result, "please input drive name first") {
		result, err = p.Cli.Cmd("show config active")
	}
	//p.Debug(result)

	// we should strip wrong data from this result. Starting from `Command: ....`, ending with prompt
	re, err := regexp.Compile(`(?msi:(Command\: [^\n]+\n)(?P<config>.*)$)`)
	out := p.ParseSingle(re, result)
	config := strings.Trim(out["config"], " ")

	if len(config) < 100 {
		return config, fmt.Errorf("Something wrong: too short config")
	}

	return config, nil
}
