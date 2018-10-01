package discoverer

import (
	"fmt"
	"github.com/ircop/dproto"
	"strings"
	"regexp"
	"strconv"
)

/*const (
	// Unknown type of interface
	IntTypeUnknown		= 0
	// Physical port
	IntTypePhisycal		= 1
	// Port-channel
	IntTypeAggregated	= 2
	// SVI
	IntTypeSvi			= 3
	// Tunnel interface
	IntTypeTunnel		= 4
	// Lo
	IntTypeLoopback		= 5
	// Mgmt
	IntTypeManagement	= 6
	// Nukk
	IntTypeNull			= 7

	// Interface is enabled
	IntAdminStateUP		= 1
	// Interface is disabled
	IntAdminStateDown	= 2
)*/

/*// Interface template
type Interface struct {
	// Interface state
	Type		int
	// Operational state
	//State		int
	// Administrative state
	//AdminState	int
	// Full name (like GigabitEthernet 0/1)
	Name		string
	// Shortname (like Gi0/1)
	Shortname	string
	// Interface description
	Description	string
	// LLDP Id of interface. Usually name or mac.
	LldpID		string
	// Children: names of port-channel members
	PoMembers	[]string
}*/

// GetInterfaces gathers interfaces data on the device
func (p *Generic) GetInterfaces() (map[string]*dproto.Interface, error) {
	//return make(map[string]dproto.Interface), ErrNotImplemented
	return make(map[string]*dproto.Interface), fmt.Errorf("Sorry, GetInterfaces() not implemented for this profile")
}

// GetInterfaceType determines interface type by interface short name
func (p *Generic) GetInterfaceType(ifname string) dproto.InterfaceType {

	if match, _ := regexp.Match(`^(ae|eth-trunk|po|bond|t\d+$)`, []byte(strings.ToLower(ifname))); match{
		return dproto.InterfaceType_AGGREGATED
	} else if match, _ := regexp.Match(`^(40|fa|xe|xg|ge|gi|te|et|wlan|sfp|ether|g)`, []byte(strings.ToLower(ifname))); match {
		if strings.Contains(ifname, ".") {
			return dproto.InterfaceType_SVI
		}
		return dproto.InterfaceType_PHISYCAL
	} else if match, _ := regexp.Match(`^(vl|irb|bridg|vlan)`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_SVI
	} else if match, _ := regexp.Match(`^(lo)`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_LOOPBACK
	} else if match, _ := regexp.Match(`^(fxp|mg|meth)`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_MANAGEMENT
	} else if match, _ := regexp.Match(`^(tu|ppp|l2t|pptp|ovpn|sstp|gre|ipip|eoip)`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_TUNNEL
	}  else if match, _ := regexp.Match(`^(null)`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_NULL
	} else if match, _ := regexp.Match(`^(\d+(:\d+)?)$`, []byte(strings.ToLower(ifname))); match {
		return dproto.InterfaceType_PHISYCAL
	}

	return dproto.InterfaceType_UNKNOWN
}

/*
ExpandInterfaceRange func
 * Convert interface range to list:
 * "Gi 1/1-3,Gi 1/7" -> ["Gi 1/1", "Gi 1/2", "Gi 1/3", "Gi 1/7"]
 * "1:1-3" -> ["1:1", "1:2", "1:3"]
 * "1:1-1:3" -> ["1:1", "1:2", "1:3"]
 * todo: something like 1:(1,3-24)
 */
func (p *Generic) ExpandInterfaceRange(ifstring string) []string {
	result := make([]string,0)

	rePrefix, err := regexp.Compile(`^(?P<prefix>.*?)(?P<num>\d+)$`)
	if err != nil {
		p.Log("ExpandInterfaceRange: cannot compile interface prefix regexp: %s", err.Error())
		return result
	}

	// dgs3100 style: 1:(3,4,5-8)
	prefix2 := ""
	rePrefix2, err := regexp.Compile(`^(?P<prefix>\d+):\((?P<ports>[^\)]+)\)$`)
	if err != nil {
		p.Log("ExpandInterfaceRange: cannot compile interface prefix2 regexp: %s", err.Error())
		return result
	}

	if rePrefix2.Match([]byte(ifstring)) {
		out := p.ParseSingle(rePrefix2, ifstring)
		ifstring = strings.Trim(out["ports"], " ")
		prefix2 = strings.Trim(out["prefix"], " ")
	}


	list1 := strings.Split(ifstring, ",")
	for _, x := range list1 {
		x = strings.Trim(x, " ")
		if "" == x {
			continue
		}

		if match, _ := regexp.Match("-", []byte(x)); match {

			var prefix, startStr, stopStr string

			rePrefixed, err := regexp.Compile(`(?P<prefix>\d+:)\((?P<start>\d+)-(?P<stop>\d+)\)`)
			if err != nil {
				p.Log("Cannot compile rePrefixed regexp")
				continue
			}
			prefixed := p.ParseSingle(rePrefixed, x)
			if "" != prefixed["start"] && "" != prefixed["stop"] && "" != prefixed["prefix"] {
				prefix = prefixed["prefix"]
				startStr = prefixed["start"]
				stopStr = prefixed["stop"]
			} else {
				// expand range
				list2 := strings.Split(x, "-")
				if len(list2) < 2 {
					continue
				}

				from := strings.Trim(list2[0], " ")
				to := strings.Trim(list2[1], " ")

				// detect common prefix
				out := p.ParseSingle(rePrefix, from)
				prefix = out["prefix"]
				startStr = out["num"]
				out = p.ParseSingle(rePrefix, to)
				if prefix != out["prefix"] && "" != out["prefix"] {
					p.Log("ExpandInterfaceRange: start prefix doesnt equals to stop prefix ('%s'|'%s')", prefix, out["prefix"])
					continue
				}
				stopStr = out["num"]
				if "" == startStr || "" == stopStr {
					p.Log("ExpandInterfaceRange: start/stop interfaces are empty ('%s'|'%s') (string given: '%s')", startStr, stopStr, ifstring)
					continue
				}
			}

			start, err := strconv.ParseInt(startStr, 10, 64)
			if err != nil {
				p.Log("ExpandInterfaceRange: start interfaces is not integer (%s)", startStr)
				continue
			}
			stop, err := strconv.ParseInt(stopStr, 10, 64)
			if err != nil {
				p.Log("ExpandInterfaceRange: stop interfaces is not integer (%s)", startStr)
				continue
			}

			for i := start; i <= stop; i++ {
				ifname := fmt.Sprintf("%s%d", prefix, i)
				result = append(result, ifname)
			}
		} else {
			result = append(result, x)
		}
	}

	if prefix2 != "" {
		for i := range result {
			result[i] = prefix2 + ":" + result[i]
		}
	}

	return result
}
