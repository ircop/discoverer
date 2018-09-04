package DLinkDGS3100

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"regexp"
	"strings"
)

// GetLldp for DLinkDGS3100 profile
func (p *Profile) GetLldp() ([]*dproto.LldpNeighbor, error) {
	p.Debug("starting DLinkDGS3100.GetLldp()")
	neighbors := make([]*dproto.LldpNeighbor, 0)

	result, err := p.Cli.Cmd("show lldp remote_ports")
	if err != nil {
		return neighbors, fmt.Errorf("Cannot 'show lldp remote_ports': %s", err.Error())
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?msi:Port ID\s+:\s+(?P<port>\d:\d+))(\s+)?\n[-]+\n\nEntity[^\n]+\nChas[^\n]+\nChassis ID\s+:\s+(?P<cid>[^\n]+)\nPort[^\n]+\nPort ID\s+:\s+(?P<pid>[^\n]+)\n`)
	if err != nil {
		return neighbors, fmt.Errorf("Cannot compile lldp regex")
	}

	out := p.ParseMultiple(re, result)
	for _, part := range out {
		port := strings.Trim(part["port"], " ")
		cid := strings.Trim(part["cid"], " ")
		pid := strings.Trim(part["pid"], " ")
		nei := dproto.LldpNeighbor{
			LocalPort:port,
			ChassisID:cid,
			PortID:pid,
		}
		neighbors = append(neighbors, &nei)
	}

	return neighbors, nil
}
