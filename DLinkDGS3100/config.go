package DLinkDGS3100

import (
	"fmt"
	"strings"
)

// GetConfig for DLinkDGS3100
func (p *Profile) GetConfig() (string, error) {
	p.Debug("Starting DLinkDGS3100.GetConfig()")

	result, err := p.Cli.Cmd("sh configuration running")
	if err != nil {
		return "", fmt.Errorf("Cannot 'sh configuration running': %s", err.Error())
	}

	// strip line with command ('sh configuration running')
	lines := strings.Split(result, "\n")
	r2 := ""
	for _, line := range lines {
		if strings.Contains(line, "sh configuration running") {
			continue
		}
		r2 += line + "\n"
	}

	return r2, nil
}
