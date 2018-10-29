package QtechQSW

import (
	"github.com/ircop/discoverer/util/text"
	"github.com/ircop/dproto"
	"strconv"
	"strings"
)

// GetVlans for QtechQSW
func (p *Profile) GetVlans() ([]*dproto.Vlan, error) {
	p.Log("Starting QtechQSW.GetVlans()")
	vlans := make([]*dproto.Vlan, 0)

	result, err := p.Cli.Cmd("show vlan")
	if err != nil {
		return vlans, err
	}
	p.Debug(result)

	rows := text.ParseTable(result, `^---`, ``, false, true)
	for i := range rows {
		row := rows[i]
		vidString := strings.Trim(row[0], " ")
		vid, err := strconv.ParseInt(vidString, 10, 64)
		if err != nil {
			continue
		}
		portstring := row[4]
		ports := strings.Split(portstring, " ")
		access := make([]string, 0)
		trunk := make([]string, 0)
		for j := range ports {
			port := strings.Trim(ports[j], " ")
			if port == "" {
				continue
			}

			if strings.Contains(port, "(T)") {
				port = strings.Replace(port, "(T)", "", -1)
				trunk = append(trunk, port)
			} else {
				access = append(access, port)
			}
		}

		vlan := dproto.Vlan{
			Name:vidString,
			ID:vid,
			TrunkPorts:trunk,
			AccessPorts:access,
		}
		vlans = append(vlans, &vlan)
	}


	return vlans, nil
}