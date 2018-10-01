package DLinkDxS

import (
	"fmt"
	"github.com/ircop/dproto"
	"net"
	"regexp"
)

// GetPlatform gathers all major info about DLink swithc (except 3100 platforms)
func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Log("Starting DLinkDxS.GetPlatform()")

	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	var reModelStr = `(?ms:[Dd]evice [Tt]ype\s+:\s*(?P<model>\S+).+[Ff]irmware [Vv]ersion(?: 1)?\s+:\s*(?:Build\s+)?(?P<version>\S+).+[Hh]ardware [Vv]ersion\s+:\s*(?P<revision>\S+))`
	var reSerialStr = `(?:[Ss]erial [Nn]umber|Device S/N)\s+:\s*(?P<serial>\S+)\s*\n`
	var reMacStr = `(?m:^MAC Address\s+:\s*(?P<mac>\S+))`

	// todo: check if cli/snmp available, etc...
	if p.Cli == nil {
		return platform, fmt.Errorf("No cli provided, giving up.")
	}

	result, err := p.Cli.Cmd("show switch")
	if err != nil {
		fmt.Printf(result)
		return platform, fmt.Errorf("Error GetPlatform: %s", err.Error())
	}
	p.Debug("---\n%s\n---\n",result)

	/*bts := []byte(result)
	for i := range bts {
		p.Debug("%s | %d | %x", string(bts[i]), bts[i], bts[i])
		if i > 250 {
			break
		}
	}
	return platform, nil*/

	reModel, err := regexp.Compile(reModelStr)
	if err != nil {
		return platform, fmt.Errorf("DLinkDxS: GetPlatform: Cannot compile model regex: %s", err.Error())
	}
	reSerial, err := regexp.Compile(reSerialStr)
	if err != nil {
		return platform, fmt.Errorf("DLinkDxS: GetPlatform: Cannot compile serial regex: %s", err.Error())
	}
	reMac, err := regexp.Compile(reMacStr)
	if err != nil {
		return platform, fmt.Errorf("DLinkDxS: GetPlatform: Cannot compile mac regex: %s", err.Error())
	}

	// get model, rev, fw from 1st regex
	out := p.ParseSingle(reModel, result)
	if model, ok := out["model"]; ok {
		platform.Model = model
	}
	if rev, ok := out["revision"]; ok {
		platform.Revision = rev
	}
	if fw, ok := out["version"]; ok {
		platform.Version = fw
	}


	// get serial from 2nd regex
	out = p.ParseSingle(reSerial, result)
	if serial, ok := out["serial"]; ok {
		platform.Serial = serial
	}

	// Get mac from 3rd regex
	out = p.ParseSingle(reMac, result)
	if mac, ok := out["mac"]; ok {
		m, err := net.ParseMAC(mac)
		if err == nil {
			platform.Macs = append(platform.Macs, m.String())
		}
	}

	return platform, nil
}
