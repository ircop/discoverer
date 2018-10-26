package QtechQSW

import (
	"fmt"
	"github.com/ircop/dproto"
	"net"
)

// GetIps for QtechQSW
func (p *Profile) GetIps() ([]*dproto.Ipif, error) {
	ipifs := make([]*dproto.Ipif, 0)
	p.Log("Starting QtechQSW.GetIps()")

	patterns := make(map[string]string, 0)
	patterns["brief"] = `(?msi:^\d+\s+(?P<ifname>[^\s]+)\s+(?P<ip>[^\s]+)\s+(up|down))`
	patterns["ips"] = `(?msi:(?P<ip>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mask>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+\((Primary|Secondary)\))`
	regexps, err := p.CompileRegexps(patterns)
	if err != nil {
		return ipifs, err
	}

	result, err := p.Cli.Cmd("sh ip interface brief")
	if err != nil {
		return ipifs, fmt.Errorf("Cannot 'sh ip interface brief': %s", err.Error())
	}
	p.Debug(result)

	out := p.ParseMultiple(regexps["brief"], result)
	for i := range out {
		ip := out[i]["ip"]
		ifname := out[i]["ifname"]
		if ip == "127.0.0.1" || ifname == "" || net.ParseIP(ip) == nil {
			continue
		}

		//p.Debug("ip: %s, if: %s", ip, ifname)
		r, err := p.Cli.Cmd(fmt.Sprintf("show int %s", ifname))
		if err != nil {
			p.Debug("Failed to 'show int %s': %s", ifname, err.Error())
			continue
		}
		p.Debug(r)
		o := p.ParseMultiple(regexps["ips"], r)
		for n := range o {
			ipStr := o[n]["ip"]
			maskStr := o[n]["ip"]
			ip := net.ParseIP(ipStr)
			mask := net.ParseIP(maskStr)
			if ip == nil || mask == nil {
				continue
			}

			ipif := dproto.Ipif{
				Interface:ifname,
				IP:ip.String(),
				Mask:mask.String(),
			}
			ipifs = append(ipifs, &ipif)
		}
	}

	return ipifs, nil
}