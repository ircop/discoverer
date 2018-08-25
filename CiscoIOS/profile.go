package CiscoIOS

import (
	"github.com/ircop/discoverer/base"
	"github.com/ircop/remote-cli"
	"fmt"
	"strings"
)

// Profile instance
type Profile struct {
	discoverer.Generic
}

// SetPrompt sets CLI prompt for current profile (only if CLI is active)
func (p *Profile) SetPrompt() {
	if p.Cli != nil {
		p.Cli.SetPrompt(`(?msi:^[a-zA-Z0-9\-_]+[\$%#>]$)`)
	}
}

func (p *Profile) Init(cli *remote_cli.Cli, enable string, community string) error {
	err := p.InitShared(cli, enable, community)
	if err != nil {
		return err
	}

	p.SetPrompt()
	p.Cli.GlobalTimeout(60)

	if err := p.Cli.RegisterErrorPattern(`(% Invalid input detected at|% Ambiguous command|% Incomplete command|% Unknown command)`, "Syntax error"); err != nil {
		return err
	}

	// Enable first if needed
	if enable != "" {
		p.Cli.Write([]byte("enable"))
		// next will be password prompt or cisco prompt
		r, err := p.Cli.ReadUntil(`([Pp]ass[Ww]ord:|[a-zA-Z0-9\-]+[\%$%#>]$)`)
		if err != nil {
			fmt.Printf(strings.Replace(r, "%", "%%", -1))
			panic(err)
		} else {
			p.Cli.Cmd(enable)
		}
	}

	return nil
}
