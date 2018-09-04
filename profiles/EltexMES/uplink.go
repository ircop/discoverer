package EltexMES

import (
	"fmt"
	"strings"
	"regexp"
)

// GetUplink for DLinkDxS
func (p *Profile) GetUplink() (string, error) {
	p.Log("Starting DLinkDxS.GetUplink()")

	// Some models normaly shows 'sh ip route', but some - show routes in ip interface
	if p.Version == "" {
		return "", fmt.Errorf("Empty version. Please run GetPlatform() first.")
	}


	if strings.HasPrefix(p.Version, "1.") {
		return p.routeIpInterface()
	}
	return p.routeIproute()
}

func (p *Profile) routeIpInterface() (string, error) {
	result, err := p.Cli.Cmd("show ip interface")
	if err != nil {
		return result, fmt.Errorf("Cannot 'show ip interface': %s", err.Error())
	}
	p.Debug(result)

	reGw, err := regexp.Compile(`(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+Active\s+`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile gw regex: %s", err.Error())
	}
	out := p.ParseSingle(reGw, result)
	gw := strings.Trim(out["gw"],  " ")
	if gw == "" {
		return "", nil
	}

	result, err = p.Cli.Cmd(fmt.Sprintf("sh arp ip-address %s", gw))
	if err != nil {
		return "", fmt.Errorf("Cannot 'sh arp': %s", err.Error())
	}
	p.Debug(result)
	reArp, err := regexp.Compile(`\s+(?P<iface>[^\s]+)\s+(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile arp regex: %s", err.Error())
	}
	out = p.ParseSingle(reArp, result)
	iface := out["iface"]

	return iface, nil
}

func (p *Profile) routeIproute() (string, error) {
	result, err := p.Cli.Cmd("sh ip route address 0.0.0.0")
	if err != nil {
		return result, fmt.Errorf("Cannot 'show ip rouote 0.0.0.0': %s", err.Error())
	}
	p.Debug(result)

	reGw, err := regexp.Compile(`0.0.0.0\/0(.*)via\s(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile gw regex: %s", err.Error())
	}
	out := p.ParseSingle(reGw, result)
	gw := strings.Trim(out["gw"], " ")
	if gw == "" {
		return "", nil
	}

	// arp
	result, err = p.Cli.Cmd(fmt.Sprintf("sh arp ip-address %s", gw))
	if err != nil {
		return "", fmt.Errorf("Cannot 'sh arp': %s", err.Error())
	}
	p.Debug(result)
	reArp, err := regexp.Compile(`\s+(?P<iface>[^\s]+)\s+(?P<gw>\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)\s+(?P<mac>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))`)
	if err != nil {
		return "", fmt.Errorf("Cannot compile arp regex: %s", err.Error())
	}
	out = p.ParseSingle(reArp, result)
	iface := out["iface"]

	return iface, nil
}