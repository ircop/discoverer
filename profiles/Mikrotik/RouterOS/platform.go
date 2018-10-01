package RouterOS

import (
	"fmt"
	"github.com/ircop/dproto"
	"github.com/ircop/discoverer/util/mac"
	"strings"
)

// GetPlatform for RouterOS
func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Debug("Starting RouterOS.GetPlatform()")
	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	patterns := make(map[string]string)
	patterns["model"] = `(?ms:version: (?P<q>\"?)(?P<firmware>\d+\.\d+(\.\d+)?).+board-name: (?P<qp>\"?)(?P<model>\D+?.\S+?)\n)`
	patterns["macs"] = `(?m:(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})))`
	patterns["serial"] = `serial-number:\s+(?P<serial>[^\s\n]+)(\s+)?\n`
	//patterns["ports"] = `(?m:^(?m:(.+);;;\s(?P<desc>[^\n]+)\n)?(.+)\sname=\"(?P<ifname>([^\"]+))\"\s+(.+)?type=\"(?P<type>([a-zA-Z0-9]+))\"\s+(.+)?mac-address=(?P<mac>[a-fA-F0-9\:]+)?)`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return platform, err
	}

	result, err := p.Cli.Cmd("system resource print")
	if err != nil {
		return platform, err
	}
	p.Debug(result)

	out := p.ParseSingle(regexps["model"], result)
	platform.Version = strings.Trim(out["firmware"], " ")
	platform.Model = strings.Trim(out["model"], " ")

	// macs: parse `interfaces print`. Macs could be non-unique, so collect macs as map keys.
	result, err = p.Cli.Cmd("/interface print")
	if err != nil {
		return platform, fmt.Errorf("Cannot '/interface print': %s", err.Error())
	}
	p.Debug(result)

	macs := make(map[string]int)
	parts := p.ParseMultiple(regexps["macs"], result)
	for _, part := range parts {
		mac := Mac.New(part["mac"])
		if mac == nil {
			continue
		}
		macs[mac.String()] = 1
	}

	for mac, _ := range macs {
		platform.Macs = append(platform.Macs, mac)
	}

	// serial: not for cxr/x86
	if !strings.Contains(platform.Model, "x86") && !strings.Contains(platform.Model, "CHR") {
		result, err = p.Cli.Cmd("/system routerboard print")
		if err != nil {
			return platform, err
		}
		p.Debug(result)

		out := p.ParseSingle(regexps["serial"], result)
		s := strings.Trim(out["serial"], " ")
		platform.Serial = s
	}

	return platform, nil
}
