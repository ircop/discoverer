package cfg

import (
	"github.com/franela/goblin"
	"testing"
)

func TestCfg_CheckParams(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#NewCfg", func() {
		g.It("Should return error if wrong config path given", func() {
			c, err := NewCfg("")
			g.Assert(err == nil).IsFalse("Should return error with wrong config path")
			g.Assert(c == nil).IsTrue("Should return empty Cfg pointer with wrong config path")
		})
	})

	g.Describe("#CheckParams", func() {
		g.It("Should return error with empty parameters", func() {
			c := new(Cfg)
			err := c.CheckParams()
			g.Assert(err == nil).IsFalse()
		})

		g.It("Should return error with wrong ip and right port", func() {
			c := new(Cfg)
			c.ListenIP = "abrakadabra"
			c.ListenPort = 2345
			err := c.CheckParams()
			g.Assert(err == nil).IsFalse()
		})

		g.It("Should return error with correct ip and wrong port", func() {
			c := new(Cfg)
			c.ListenIP = "0.0.0.0"
			c.ListenPort = 10
			err := c.CheckParams()
			g.Assert(err == nil).IsFalse()

			c.ListenPort = 66666
			err = c.CheckParams()
			g.Assert(err == nil).IsFalse()
		})
	})
}
