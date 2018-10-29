package QtechQSW

import (
	"github.com/ircop/discoverer/util/mac"
	"github.com/ircop/discoverer/util/text"
	"github.com/ircop/dproto"
	"strings"
)

// GetLldp for QtechQSW profile
func (p *Profile) GetLldp() ([]*dproto.LldpNeighbor, error) {
	p.Log("starting QtechQSW.GetLldp()")
	neighbors := make([]*dproto.LldpNeighbor, 0)

	patterns := make(map[string]string, 0)
	//patterns["neis"] = `(?msi:^(?P<ifname>[^\s]+)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+\d+\s+(?P<port>[^\s]+)\s+)`
	patterns["item"] = `(?msi:Port name\s+:\s+(?P<local>[^\n\s]+)\n.*ChassisID\s+:(\s+)?(?P<cid>[^\n\s]+)\n.*PortId\s+:(\s+)?(?P<pid>[^\n\s]+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return neighbors, err
	}

	result, err := p.Cli.Cmd("show lldp neighbors interface")
	if err != nil {
		return neighbors, err
	}
	p.Debug(result)

	parts, err := text.SplitByParts(result, `^Port name :`)
	if err != nil {
		return neighbors, err
	}
	for i := range parts {
		part := parts[i]
		out := p.ParseSingle(regexps["item"], part)
		local := strings.Trim(out["local"], " ")
		cid := strings.Trim(out["cid"], " ")
		pid := strings.Trim(out["pid"], " ")
		if local == "" || cid == "" || pid == "" || !Mac.IsMac(cid) {
			continue
		}

		item := dproto.LldpNeighbor{
			LocalPort:local,
			ChassisID:cid,
			PortID:pid,
		}
		neighbors = append(neighbors, &item)
	}

	/*result, err := p.Cli.Cmd("show lldp neighbors brief")
	if err != nil {
		return neighbors, err
	}
	p.Debug(result)

	out := p.ParseMultiple(regexps["neis"], result)
	for i := range out {
		//p.Debug("%+#v", out[i])
		ifname := strings.Trim(out[i]["ifname"], " ")
		mac := strings.Trim(out[i]["mac"], " ")
		port := strings.Trim(out[i]["port"], " ")
		if ifname == "" || !Mac.IsMac(mac) || port == "" {
			continue
		}

		item := dproto.LldpNeighbor{
			LocalPort:ifname,
			PortID:port,
			ChassisID:mac,
		}
		neighbors = append(neighbors, &item)
	}
	*/
	return neighbors, nil
}
