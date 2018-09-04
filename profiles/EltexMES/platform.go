package EltexMES

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/mac"
	"strings"
)

// GetPlatform for EltexMES
func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Debug("Starting EltexMES.GetPlatform()")
	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	patterns := make(map[string]string)
	//patterns["ver"] = `Active-image[^\n]+\n\s+Version:\s+(?P<ver>[^\n]+)\n`
	patterns["ver"] = `(Active-image[^\n]+\n\s+Version:\s+(?P<ver>[^\n]+)\n|SW version\s+(?P<ver2>\d+(\.\d+)?(\.\d+)?))`
	patterns["sys"] = `(?ms:System Description(\s+)?:\s*(?P<model>\S+).+System MAC Address(\s+)?:\s*(?P<mac>\S+)$)`
	patterns["sysid"] = `(----(\s+)?\n\s+\d+\s+(?P<serial>[^\s]+)(\s+)?|Serial number\s+:\s+(?P<serial2>[^\s]+))`
	patterns["rev"] = `HW version\s+(?P<rev>[^\s]+)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return platform, fmt.Errorf("Cannot compile regexps: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show version")
	if err != nil {
		return platform, err
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["ver"], result)
	ver := strings.Trim(out["ver"], " ")
	platform.Version = ver
	if platform.Version == "" {
		platform.Version = strings.Trim(out["ver2"], " ")
	}

	out = p.ParseSingle(regexps["rev"], " ")
	platform.Revision = strings.Trim(out["rev"], " ")

	result, err = p.Cli.Cmd("show system")
	if err != nil {
		return platform, fmt.Errorf("Cannot execute 'show system': %s", err.Error())
	}
	p.Debug(result)

	out = p.ParseSingle(regexps["sys"], result)
	platform.Model = strings.Trim(out["model"], " ")
	if mac := Mac.New(out["mac"]); mac != nil {
		platform.Macs = append(platform.Macs, mac.String())
	}

	result, err = p.Cli.Cmd("show system id")
	if err != nil {
		return platform, fmt.Errorf("Cannot execute 'show system id': %s", err.Error())
	}
	p.Debug(result)

	out = p.ParseSingle(regexps["sysid"], result)
	platform.Serial = strings.Trim(out["serial"], " ")
	if platform.Serial == "" {
		platform.Serial = strings.Trim(out["serial2"], " ")
	}

	p.Version = platform.Version
	p.Model = platform.Model

	return platform, nil
}
