package DLinkDGS3100

import (
	"fmt"
	"github.com/ircop/dproto"
	"github.com/ircop/discoverer/util/mac"
	"regexp"
	"strings"
)

// GetPlatform for DLinkDGS3100 profile
func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Debug("Starting DLinkDGS3100.GetInterfaces()")
	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	re, err := regexp.Compile(`(?ms:[Dd]evice [Tt]ype\s+:\s*(?P<model>\S+).+MAC Address\s+:\s+(?P<mac>[^\s]+).+[Ff]irmware [Vv]ersion(?: 1)?\s+:\s*(?:Build\s+)?(?P<version>\S+).+[Hh]ardware [Vv]ersion\s+:\s*(?P<revision>\S+).+Serial Number\s+:\s+(?P<serial>[^\s(]+))`)
	if err != nil {
		return platform, fmt.Errorf("Cannot compile regexp: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show switch")
	if err != nil {
		return platform, fmt.Errorf("Cannot 'show switch': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseSingle(re, result)
	platform.Model = strings.Trim(out["model"], " ")
	platform.Version = strings.Trim(out["version"], " ")
	platform.Revision = strings.Trim(out["revision"], " ")
	platform.Serial = strings.Trim(out["serial"], " ")
	if Mac.IsMac(out["mac"]) {
		platform.Macs = append(platform.Macs, strings.Trim(out["mac"], " "))
	}

	return platform, nil
}