package Mac

import (
	"net"
	"strings"
	"regexp"
	"strconv"
	"fmt"
)

type Mac struct {
	value		string
}

// New macaddr instance
func New(val string) *Mac {
	val = strings.Trim(val, " ")
	if val == "" {
		return nil
	}

	hw, err := net.ParseMAC(val)
	if err != nil {
		return nil
	}

	m := new(Mac)
	m.value = hw.String()

	return m
}

func (m *Mac) String() string {
	return m.value
}

func IsMac(mac string) bool {
	_, err := net.ParseMAC(mac)
	return err == nil
}

func (m *Mac) Int64() int64 {
	if m.value == "" {
		return 0
	}

	split, err := regexp.Compile(`[^a-fA-F0-9]`)
	if err != nil {
		return 0
	}

	hex := split.ReplaceAllLiteralString(m.value, "")

	intMac, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0
	}

	return intMac
}

func (m *Mac) Range(count int64) []string {
	macs := make([]string, 0)
	intFirst := m.Int64()
	for i := intFirst ; i < (intFirst + count) ; i++ {
		// make macaddr from int64
		cur := i
		var n int64
		r := make([]string, 0)
		for n=0; n<6; n++ {
			r = append(r, fmt.Sprintf("%02x", cur & 0xFF))
			cur = cur >> 8
		}
		reversed := make([]string, 0)
		for j := len(r)-1; j >= 0; j-- {
			reversed = append(reversed, r[j])
		}
		macs = append(macs, strings.ToLower(strings.Join(reversed, ":")))
	}
	return macs
}
