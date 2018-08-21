package discoverer

import (
	"net"
	"fmt"
)

type IPInterface struct {
	Interface	string
	IP			net.IP
	Mask		net.IP
}

// GetIps dummy
func (p *Generic) GetIps() ([]IPInterface, error) {
	return make([]IPInterface,0), fmt.Errorf("Sorry, GetIps() not implemented in current profile")
}