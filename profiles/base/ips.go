package discoverer

import (
	"fmt"
	"github.com/ircop/dproto"
)

/*type IPInterface struct {
	Interface	string
	IP			net.IP
	Mask		net.IP
}*/

// GetIps dummy
func (p *Generic) GetIps() ([]*dproto.Ipif, error) {
	//return make([]dproto.Ipif,0), ErrNotImplemented
	return make([]*dproto.Ipif,0), fmt.Errorf("Sorry, GetIps() not implemented in current profile")
}