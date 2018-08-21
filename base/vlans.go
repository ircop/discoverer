package discoverer

import (
	"fmt"
	"strings"
	"strconv"
)

// Vlan structure
type Vlan struct {
	Name		string
	ID			int64
	TrunkPorts	[]string
	AccessPorts	[]string
}

// GetVlans dummy
func (p *Generic) GetVlans() ([]Vlan, error) {
	return make([]Vlan, 0), fmt.Errorf("Sorry, GetVlans not implemented in current profile")
}

// ExpandVlanRange
// 1-5,10,15 -> [1,2,3,4,5,10,15], etc.
func (p *Generic) ExpandVlanRange(vstring string) []string {
	result := make([]string,0)

	vstring = strings.Replace(vstring, "\n", "", -1)

	list1 := strings.Split(vstring, ",")
	for _, x := range list1 {
		x = strings.Trim(x, " ")
		if x == "" {
			continue
		}

		//if match, _ := regexp.Match("-", []byte(x)); match {
		if strings.Contains(x, "-") {
			list2 := strings.Split(x, "-")
			if len(list2) != 2 {
				continue
			}

			from := strings.Trim(list2[0], " ")
			to := strings.Trim(list2[1], " ")
			//fmt.Printf("FROM %s TO %s\n", from, to)
			start, e1 := strconv.ParseInt(from, 10, 64)
			end, e2 := strconv.ParseInt(to, 10, 64)
			if e1 != nil || e2 != nil {
				continue
			}

			for i := start; i <= end; i++ {
				result = append(result, fmt.Sprintf("%d", i))
			}
		} else {
			result = append(result, x)
		}
	}

	return result
}