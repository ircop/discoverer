package DLinkDxS

import (
	"github.com/ircop/discoverer/dproto"
	"github.com/pkg/errors"
	"fmt"
	//"regexp"
	"regexp"
	"strings"
)

// GetLldp for DLinkDxS profile
func (p *Profile) GetLldp() ([]*dproto.LldpNeighbor, error) {
	p.Log("starting DLinkDxS.GetLldp()")
	neighbors := make([]*dproto.LldpNeighbor, 0)

	out, err := p.Cli.Cmd("show lldp remote_ports mode normal")
	if err != nil {
		return neighbors, errors.Wrap(err, "Cannot get lldp neighborship")
	}
	p.Debug(out)

	reParts, err := regexp.Compile(`Port ID\s:`)
	if err != nil {
		return neighbors, fmt.Errorf("Cannot compile PortID regex")
	}
	rePort, err := regexp.Compile(`(?ms:^Port ID : (?P<port>(\d:)?\d+)\s*\n^-+\s*\n^[Rr]emote [Ee]ntities [Cc]ount : (?P<count>[0-9]+)\s*\n(?P<entities>.+))`)
	if err != nil {
		return neighbors, fmt.Errorf("Cannot compile rePort regex")
	}
	reEntity, err := regexp.Compile(`(?mis:^Entity \d+\s*\n^\s+Chassis ID Subtype\s+:(\s+)?(?P<chassis_id_subtype>.*?)\s*\n^\s+Chassis ID\s+:(\s+)?(?P<chassis_id>.*?)\s*\n^\s+Port ID Subtype\s+:(\s+)?(?P<port_id_subtype>.*?)\s*\n^\s+Port ID\s+:(\s+)?(?P<port_id>.*?)\s*\n^\s*Port Description\s+:(?P<port_description>.*?)^\s+System Name\s+:(?P<system_name>.*?)^\s+System Description\s+:(?P<system_description>.*?)^\s+System Capabilities\s+:(?P<system_capabilities>.*?)\s*^\s+Management Address Count)`)
	if err != nil {
		return neighbors, fmt.Errorf("Cannot compile entity regex")
	}

	// split output by 'Port ID'
	parts := reParts.Split(out, -1)
	p.Debug("lldp: found %d port entities, splitted by Port ID", len(parts))

	for _, part := range parts {
		// each port info
		// recover splitted data
		port := fmt.Sprintf("Port ID :%s", part)
		out := p.ParseSingle(rePort, port)
		portName, portOk := out["port"]
		portName = strings.Trim(portName, " ")
		if !portOk || portName == "" {
			continue
		}

		// parse entities for current port
		entities := p.ParseMultiple(reEntity, port)
		p.Debug("Found %d entities on port %s", len(entities), portName)
		//members := make([]discoverer.LldpNeighbor, 0)
		for _, ent := range entities {
			cid := strings.Trim(ent["chassis_id"], " ")
			pid := strings.Trim(ent["port_id"], " ")
			if cid == "" || pid == "" {
				p.Log("Error: %s: no chassis id or port id (%s/%s)", portName, cid, pid)
				continue
			}
			item := dproto.LldpNeighbor{
				LocalPort:portName,
				ChassisID:cid,
				PortID:pid,
			}
			neighbors = append(neighbors, &item)
		}
	}

	return neighbors, nil
}
