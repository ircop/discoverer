package EltexMES

import (
	"github.com/ircop/discoverer/base"
	"github.com/pkg/errors"
	"github.com/ircop/discoverer/util/text"
)

// GetLldp for EltexMES profile
func (p *Profile) GetLldp() ([]discoverer.LldpNeighborship, error) {
	p.Debug("starting EltexMES.GetLldp()")
	neighbors := make([]discoverer.LldpNeighborship, 0)

	result, err := p.Cli.Cmd("sh lldp neighbors")
	if err != nil {
		return neighbors, errors.Wrap(err, "Cannot get lldp neighborship")
	}
	p.Debug(result)

	rows := text.ParseTable(result, "^-----", "")
	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		ifname := row[0]
		cid := row[1]
		pid := row[2]
		if ifname == "" || cid == "" || pid == "" {
			p.Log("LLDP: Warning: empty ifname/cid/pid:(%s/%s/%s)", ifname, cid, pid)
			continue
		}

		member := discoverer.LldpNeighbor{
			PortID:pid,
			ChassisID:cid,
		}
		item := discoverer.LldpNeighborship{
			PortName:ifname,
			Members: []discoverer.LldpNeighbor{member},
		}
		neighbors = append(neighbors, item)
	}

	return neighbors, nil
}
