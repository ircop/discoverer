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

func (p *Profile) Init(cli remote_cli.CliInterface, enable string, community string) error {
	err := p.InitShared(cli, enable, community)
	if err != nil {
		return err
	}

	if p.Cli != nil {
		p.SetPrompt()
		p.Cli.GlobalTimeout(60)
		p.Cli.SetPrompt(`(?msi:^[a-zA-Z0-9\-_]+[\$%#>]$)`)

		if err := p.Cli.RegisterErrorPattern(`(% Invalid input detected at|% Ambiguous command|% Incomplete command|% Unknown command)`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
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
	}
	return nil
}
