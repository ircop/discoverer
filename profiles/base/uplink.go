package discoverer

import "fmt"

// GetUplink dummy
func (p *Generic) GetUplink() (string, error) {
	//return "", ErrNotImplemented
	return "", fmt.Errorf("sorry, GetUplink not implemented in current profile")
}
