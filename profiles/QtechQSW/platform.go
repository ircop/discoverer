package QtechQSW

import (
	"github.com/ircop/discoverer/util/mac"
	"github.com/ircop/dproto"
	"strings"
)

func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Debug("Starting QtechQSW.GetPlatform()")
	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	result, err := p.Cli.Cmd("show version")
	if err != nil {
		return platform, err
	}
	p.Debug(result)

	patterns := make(map[string]string,0)
	patterns["model"] = `^\s*(?:Device: )?(?P<model>\S+)(?: Device|, sysLocation\:).+\n`
	patterns["version"] = `(?msi:^\s*SoftWare(?: Package)? Version\s+(?P<version>[^\n^\(]+)(?:\(\S+\))?\n)`
	patterns["revision"] = `(?msi:^\s*Hardware(?: Package)? Version\s+(?P<revision>[^\n^\(]+)(?:\(\S+\))?\n)`
	patterns["mac"] = `Vlan MAC\s+(?P<mac>[^\n]+)\n`
	patterns["serial"] = `Serial No\.:(\s+)?(?P<serial>[^\n]+)\n`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return platform, err
	}

	out := p.ParseSingle(regexps["model"], result)
	platform.Model = strings.Trim(out["model"], " ")
	out = p.ParseSingle(regexps["revision"], result)
	platform.Revision = strings.Trim(out["revision"], " ")
	out = p.ParseSingle(regexps["version"], result)
	platform.Version = strings.Trim(out["version"], " ")
	out = p.ParseSingle(regexps["serial"], result)
	platform.Serial = strings.Trim(out["serial"], " ")
	out = p.ParseSingle(regexps["mac"], result)
	m := Mac.New(out["mac"])
	if m != nil {
		platform.Macs = append(platform.Macs, m.String())
	}

	return platform, nil
}
