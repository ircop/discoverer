package EltexMES

import (
	"github.com/ircop/discoverer/base"
	"github.com/pkg/errors"
	"github.com/ircop/discoverer/util/text"
	"strings"
)

// GetLldp for EltexMES profile
func (p *Profile) GetLldp() ([]discoverer.LldpNeighbor, error) {
	p.Debug("starting EltexMES.GetLldp()")
	neighbors := make([]discoverer.LldpNeighbor, 0)

	result, err := p.Cli.Cmd("sh lldp neighbors")
	if err != nil {
		return neighbors, errors.Wrap(err, "Cannot get lldp neighborship")
	}
	p.Debug(result)

	rows := text.ParseTable(result, "^-----", "", false)
	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		ifname := strings.Trim(row[0], " ")
		cid := strings.Trim(row[1], " ")
		pid := strings.Trim(row[2], " ")
		if ifname == "" || cid == "" || pid == "" {
			p.Log("LLDP: Warning: empty ifname/cid/pid:(%s/%s/%s)", ifname, cid, pid)
			continue
		}

		item := discoverer.LldpNeighbor{
			LocalPort:ifname,
			ChassisID:cid,
			PortID:pid,
		}
		neighbors = append(neighbors, item)
	}

	return neighbors, nil
}
