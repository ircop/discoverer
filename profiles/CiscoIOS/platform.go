package CiscoIOS

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
	"github.com/ircop/discoverer/util/mac"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

// GetPlatform for IOS
func (p *Profile) GetPlatform() (dproto.Platform, error) {
	p.Log("Starting CiscoIOS.GetPlatform()")
	var platform dproto.Platform
	platform.Macs = make([]string, 0)

	result, err := p.Cli.Cmd("show version")
	if err != nil {
		return platform, err
	}
	p.Debug(strings.Replace(result, "%", "%%", -1))

	rePlatform, err := regexp.Compile(`(?:Cisco IOS Software( \[Everest\])?,.*?|IOS \(tm\)) (IOS[\-\s]XE Software,\s)?(?P<platform>.+?) Software \((?P<image>[^)]+)\), (Experimental )?(.+)?Version (?P<version>[^\s,]+)`)
	if err != nil {
		return platform, fmt.Errorf("Cannot compile rePlatform regex: %s", err.Error())
	}
	reSerial, err := regexp.Compile(`(?:Processor board ID (?P<serial>[^\n]+)\n)`)
	if err != nil {
		return platform, fmt.Errorf("Cannot compile reSerial regex: %s", err.Error())
	}

	out := p.ParseSingle(rePlatform, result)
	platform.Model = strings.Trim(out["platform"], " ")
	platform.Version = strings.Trim(out["version"], " ")

	out = p.ParseSingle(reSerial, result)
	platform.Serial = strings.Trim(out["serial"], " ")

	// get mac-addresses on platform
	macs, err := p.getMacs(platform.Version, platform.Model)
	if err != nil {
		p.Log("Cannot get platform macs: %s", err.Error())
	}
	platform.Macs = macs

	p.Model = platform.Model
	p.Version = platform.Version

	return platform, nil
}

func (p *Profile) getMacs(version string, model string) ([]string, error) {
	p.Debug("Starting CiscoIOS.GetMacs()")

	macs := make([]string,0)

	var count int64
	firstMac := ""

	// Command and mac discovery method depends on platform
	// todo: this is some kind of porno. Optimize to single function.
	if match, err := regexp.Match(`(SE|EA|EZ|FX|EX|EY|WC)`, []byte(version)); match && nil == err {
		p.Debug("Matched small chassis (%s)", version)
		// todo: cache results in cli class
		result, e := p.Cli.Cmd("show version")
		p.Debug(strings.Replace(result, "%", "%%", -1))

		if e != nil {
			p.Log("Error getting 'sh ver': %s", err.Error())
			return macs, nil
		}

		re, e := regexp.Compile(`(?ms:^Base ethernet MAC Address\s*:\s*(?P<mac>\S+))`)
		if e != nil {
			p.Log("Cannot compile 'base ethernet mac address' regex: %s", err.Error())
			return macs, nil
		}

		out := p.ParseSingle(re, result)
		m := Mac.New(out["mac"])
		if m == nil {
			p.Log("Failed to find first mac")
			return macs, nil
		}
		firstMac = m.String()
		count = 1
	} else if match, err := regexp.Match(`SG|\d\d\.\d\d\.\d\d\.E|EWA`, []byte(version)); match && err == nil {
		p.Debug("matched 4k chassis")
		result, e := p.Cli.Cmd("show idprom chassis")
		if e != nil {
			result, e = p.Cli.Cmd("show idprom supervisor")
			if e != nil {
				p.Log("Error getting macs: %s", err.Error())
				return macs, nil
			}
		}
		p.Debug(strings.Replace(result, "%", "%%", -1))

		re, err := regexp.Compile(`(?ms:MAC Base = (?P<mac>\S+).+MAC Count = (?P<count>\d+))`)
		if err != nil {
			return macs, errors.Wrap(err, "Cannot compile 4k chassis mac regex")
		}

		out := p.ParseSingle(re, result)
		m := Mac.New(out["mac"])
		if m == nil {
			return macs, fmt.Errorf("GetMac: cannot find first chassis mac (4k chassis)")
		}
		countStr, ok := out["count"]
		if !ok {
			return macs, fmt.Errorf("GetMac: cannot find mac count string (4k chassis)")
		}

		i, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			return macs, fmt.Errorf("GetMac: cannot parse mac count string (4k chassis) (%s)", countStr)
		}

		firstMac = m.String()
		count = i
	} else if match, err := regexp.Match(`S[YXR]`, []byte(version)); match && err == nil {
		p.Debug("matched 6k chassis")
		result, e := p.Cli.Cmd("show catalyst6000 chassis-mac-addresses")
		if e != nil {
			return macs, e
		}
		p.Debug(strings.Replace(result, "%", "%%", -1))

		re, e := regexp.Compile(`(?ms:chassis MAC addresses:.+from\s+(?P<from_id>\S+)\s+to\s+(?P<to_id>\S+))`)
		if e != nil {
			return macs, fmt.Errorf("Cannot compile 6k regex: %s", err.Error())
		}

		out := p.ParseSingle(re, result)
		mFrom := Mac.New(out["from_id"])
		mTo := Mac.New(out["to_id"])
		if mFrom == nil || mTo == nil {
			return macs, fmt.Errorf("Cannot find first or last mac (6k chassis): '%s'/'%s'", out["from_id"], out["to_id"])
		}

		fromUint := mFrom.Int64()
		toUint := mTo.Int64()
		if fromUint == 0 || toUint == 0 || fromUint > toUint {
			return macs, fmt.Errorf("6k chassis: wrong first/last mac (integers): %d/%d (%s/%s)", fromUint, toUint, mFrom.String(), mTo.String())
		}

		firstMac = mFrom.String()
		count = toUint - fromUint + 1
	} else if match, err := regexp.Match(`ASR100[0-6]`, []byte(model)); match && err == nil {
		p.Debug("Matched ASR1k chassis")
		result, e := p.Cli.Cmd("show diag chassis eeprom detail")
		if e != nil {
			return macs, e
		}
		p.Debug(strings.Replace(result, "%", "%%", -1))

		re, e := regexp.Compile(`(?m:Chassis MAC Address\s*:\s*(?P<mac>\S+)\s+MAC Address block size\s*:\s*(?P<count>\d+))`)
		if e != nil {
			return macs, fmt.Errorf("Cannot compile a1k mac regex: %s", err.Error())
		}

		out := p.ParseSingle(re, result)
		m := Mac.New(out["mac"])
		if m == nil {
			return macs, fmt.Errorf("Failed to find a1k chassis first mac (%s)", out["mac"])
		}
		countStr, ok := out["count"]
		if !ok {
			return macs, fmt.Errorf("Failed to find a1k chassis mac count")
		}

		i, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil || i < 1{
			return macs, fmt.Errorf("Failed to pare a1k chassis mac count")
		}

		firstMac = m.String()
		count = i
	} else if match, err := regexp.Match(`^C(19\d\d|29\d\d|39\d\d)$`, []byte(version)); match && err == nil {
		p.Debug("Matched 19/29/39 chassis")
		result, e := p.Cli.Cmd("show diag")
		if e != nil {
			return macs, fmt.Errorf("Failed to 'show diag' on 19/29/39 chassis: %s", e.Error())
		}
		p.Debug(strings.Replace(result, "%", "%%", -1))

		re, e := regexp.Compile(`(?m:Chassis MAC Address\s*:\s*(?P<mac>\S+)\s*\nMAC Address block size\s*:\s*(?P<count>\d+))`)
		if e != nil {
			return macs, fmt.Errorf("Failed to compile 19/29/39 chassis regex; %s", e.Error())
		}

		out := p.ParseSingle(re, result)
		m := Mac.New(out["mac"])
		if m == nil {
			return macs, fmt.Errorf("Unable to find mac on 19/29/39 chassis")
		}

		countString, ok := out["count"]
		if !ok {
			return macs, fmt.Errorf("Failed to find mac count string on 19/29/39 chassis")
		}

		i, err := strconv.ParseInt(countString, 10, 64)
		if err != nil || i < 1 {
			return macs, fmt.Errorf("Failed to parse 19/29/39 chassis mac count")
		}

		firstMac = m.String()
		count = i
	} else if match, err := regexp.Match(`(7200|7301)`, []byte(version)); match && err == nil {
		p.Debug("Matched 7200/7300 chassis")
		result, e := p.Cli.Cmd(fmt.Sprintf("show c%s | i MAC", version))
		if e != nil {
			return macs, fmt.Errorf("Failed to get 7200/7300 chassis info")
		}
		p.Debug(strings.Replace(result, "%", "%%", -1))

		re, e := regexp.Compile(`(?m:MAC Pool Size\s+(?P<count>\d+)\s+MAC Addr Base\s+(?P<mac>\S+))`)
		if e != nil {
			return macs, fmt.Errorf("Failed to compile 7200/7300 mac regex")
		}

		out := p.ParseSingle(re, result)
		m := Mac.New(out["mac"])
		if m == nil {
			return macs, fmt.Errorf("Failed to find mac for 7200/7300 chassis")
		}
		countStr, ok := out["count"]
		if !ok {
			return macs, fmt.Errorf("Failed to find mac count for 7200/7300 chassis")
		}


		i, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil || i < 1 {
			return macs, fmt.Errorf("Failed to parse 7200/7300 chassis mac count")
		}

		firstMac = m.String()
		count = i
	} else {
		return macs, fmt.Errorf("GetMac: Unsupported chassis '%s'", version)
	}

	// Now we should go from first mac, over mac count, and make macs slice
	baseMac := Mac.New(firstMac)
	if baseMac == nil {
		return macs, fmt.Errorf("GetMacs(): firstMac is nil at end of function")
	}
	if count < 1 {
		return macs, fmt.Errorf("GetMacs(): Mac count is 0 at end of function")
	}

	macs = baseMac.Range(count)

	return macs, nil
}
