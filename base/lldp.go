package discoverer

import "fmt"

// LldpNeighbor contains neighbor Platform ID and Interface ID
type LldpNeighbor struct {
	LocalPort	string
	ChassisID	string
	PortID		string
}

/*type LldpNeighborship struct {
	PortName		string
	Members			[]LldpNeighbor
}*/

// GetLldp dummy
func (p *Generic) GetLldp() ([]LldpNeighbor, error) {
	return make([]LldpNeighbor, 0), fmt.Errorf("Sorry, GetLldp not implemented in current profile")
}
