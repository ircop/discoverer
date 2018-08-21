package discoverer

import "fmt"

// LldpNeighbor contains neighbor Platform ID and Interface ID
type LldpNeighbor struct {
	ChassisID	string
	PortID		string
}

type LldpNeighborship struct {
	PortName		string
	Members			[]LldpNeighbor
}

// GetLldp dummy
func (p *Generic) GetLldp() ([]LldpNeighborship, error) {
	return make([]LldpNeighborship, 0), fmt.Errorf("Sorry, GetLldp not implemented in current profile")
}
