package discoverer

import "fmt"

// GetConfig dummy
func (p *Generic) GetConfig() (string, error) {
	//return "", ErrNotImplemented
	return "", fmt.Errorf("Sorry, GetConfig() not supported for current profile")
}
