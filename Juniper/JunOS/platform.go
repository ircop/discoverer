package JunOS

import (
	"github.com/ircop/discoverer/base"
	"fmt"
	"strings"
	"regexp"
	"github.com/ircop/discoverer/util/mac"
	"strconv"
)

// GetPlatform for JunOS
func (p *Profile) GetPlatform() (discoverer.Platform, error) {
	p.Debug("Starting Junos.GetPlatform()")
	var platform discoverer.Platform
	platform.Macs = make([]string, 0)

	patterns := make(map[string]string)
	patterns["platform"] = `(?ms:Model:\s+(?P<platform>\S+)\nJunos:\s+(?P<version>[^\n]+))`
	patterns["serial"] = `(?m:^Chassis\s+(?P<revision>REV \d+)?\s+(?P<serial>\S+)\s+(?P<rest>.+)$)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return platform, fmt.Errorf("Cannot compile regexps: %s", err.Error())
	}

	result, err := p.Cli.Cmd("show version")
	if err != nil {
		return platform, fmt.Errorf("Cannot 'show version': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["platform"], result)
	model := strings.Trim(out["platform"], " ")
	ver := strings.Trim(out["version"], " ")
	if model == "" || ver == "" {
		return platform, fmt.Errorf("Cannot parse model/version (%s/%s)", model, ver)
	}

	result, err = p.Cli.Cmd("show chassis hardware")
	if err != nil {
		return platform, fmt.Errorf("Cannot 'show chassis hardware': %s", err.Error())
	}
	p.Debug(result)

	out = p.ParseSingle(regexps["serial"], result)
	platform.Serial = strings.Trim(out["serial"], " ")

	platform.Model = model
	platform.Version = ver
	macs, err := p.getMacs()
	if err != nil {
		return platform, fmt.Errorf("Cannot get platform macs: %s", err.Error())
	}
	platform.Macs = macs

	return platform, nil
}


func (p *Profile) getMacs() ([]string, error) {
	macs := make([]string, 0)

	result, err := p.Cli.Cmd("show chassis mac-addresses")
	if err != nil {
		return macs, fmt.Errorf("Cannot run 'show chassis mac-addresses': %s", err.Error())
	}
	p.Debug(result)

	re, err := regexp.Compile(`(?msi:(\s+)?(?P<type>Public|Private) base address\s+(?P<mac>[^\n]+)\n\s+(Public|Private) count\s+(?P<count>\d+)([^\n]+)?)`)
	if err != nil {
		return macs, fmt.Errorf("Cannot compile mac regex: %s", err.Error())
	}

	out := p.ParseMultiple(re, result)
	for _, part := range out {
		first := strings.Trim(part["mac"], " ")
		cntStr := strings.Trim(part["count"], " ")
		mac := Mac.New(first)
		if mac == nil {
			return macs, fmt.Errorf("Failed to parse macaddr (%s)", first)
		}
		cnt, err := strconv.ParseInt(cntStr, 10, 64)
		if err != nil {
			return macs, fmt.Errorf("Failed to parse mac count (%s)", cntStr)
		}

		macrange := mac.Range(cnt)
		for _, m := range macrange {
			macs = append(macs, m)
		}
	}

	return macs, nil
}