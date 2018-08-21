package discoverer

import "fmt"

// Platform struct contains unified results of GetPlatform() method
type Platform struct {
	// Device model
	Model		string
	// HW revision, if exist
	Revision	string
	// FW version
	Version		string
	// Mac addresses on this platform
	Macs		[]string
	// Serial number
	Serial		string
}

// GetPlatform gathers main summary data from device and return Platform struct and/or error
func (p *Generic) GetPlatform() (Platform, error) {
	return Platform{}, fmt.Errorf("Sorry, GetPlatform() is not implemented in this profile")
}
