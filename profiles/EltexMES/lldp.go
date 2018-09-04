package EltexMES

import (
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/text"
	"github.com/pkg/errors"
	"strings"
)

// GetLldp for EltexMES profile
func (p *Profile) GetLldp() ([]dproto.LldpNeighbor, error) {
	p.Debug("starting EltexMES.GetLldp()")
	neighbors := make([]dproto.LldpNeighbor, 0)

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

		item := dproto.LldpNeighbor{
			LocalPort:ifname,
			ChassisID:cid,
			PortID:pid,
		}
		neighbors = append(neighbors, item)
	}

	return neighbors, nil
}
