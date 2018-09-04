package discoverer

import (
	"fmt"
	"github.com/ircop/discoverer/dproto"
)

// LldpNeighbor contains neighbor Platform ID and Interface ID
/*type LldpNeighbor struct {
	LocalPort	string
	ChassisID	string
	PortID		string
}*/

/*type LldpNeighborship struct {
	PortName		string
	Members			[]LldpNeighbor
}*/

// GetLldp dummy
func (p *Generic) GetLldp() ([]dproto.LldpNeighbor, error) {
	return make([]dproto.LldpNeighbor, 0), fmt.Errorf("Sorry, GetLldp not implemented in current profile")
}
