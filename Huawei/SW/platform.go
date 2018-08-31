package HuaweiSW

import (
	"fmt"
	"github.com/ircop/discoverer/base"
	"strings"
)

// GetPlatform for HuaweiSW
func (p *Profile) GetPlatform() (discoverer.Platform, error) {
	p.Debug("Starting HuaweiSW.GetPlatform()")
	var platform discoverer.Platform
	platform.Macs = make([]string, 0)

	patterns := make(map[string]string,0)
	patterns["ver"] = `(?ms:Huawei Versatile Routing Platform .+ \((?P<model>[A-Z0-9]+) (?P<version>[A-Z0-9]+))\)`
	patterns["mainboard"] = `(?msi:\[(?:Main_Board|BackPlane_0)\].+?\n\n\[Board\sProperties\](?P<body>.*?)\n\n)`
	patterns["serial"] = `(?msi:BarCode=(?P<serial>[^\n]+))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return platform, err
	}

	result, err := p.Cli.Cmd("display version")
	if err != nil {
		return platform, err
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["ver"], result)
	model := strings.Trim(out["model"], " ")
	ver := strings.Trim(out["version"], " ")
	if model == "" || ver == "" {
		return platform, fmt.Errorf("Cannot parse model/version (%s/%s)", model, ver)
	}
	platform.Model = model
	platform.Version = ver

	result, err = p.Cli.Cmd("display elabel")
	if err != nil {
		return platform, err
	}
	p.Debug(result)

	out = p.ParseSingle(regexps["mainboard"], result)
	body := out["body"]
	out = p.ParseSingle(regexps["serial"], body)
	serial := strings.Trim(out["serial"], " ")
	if serial == "" {
		p.Log("Warning! Cannot parse serial.")
	}
	platform.Serial = serial

	return platform, nil
}
