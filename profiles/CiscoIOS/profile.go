package CiscoIOS

import (
	"github.com/ircop/discoverer/profiles/base"
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
		p.Cli.SetPrompt(`(?msi:^[a-zA-Z0-9\-_.]+[\$%#>]$)`)

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
			//fmt.Printf("-EN-\n%s\n-EN-\n", r)
			if err != nil {
				str := fmt.Sprintf(strings.Replace(r, "%", "%%", -1))
				return fmt.Errorf("Error during enable: %s", str)
				//panic(err)
			} else {
				r, err = p.Cli.Cmd(enable)
				//fmt.Printf("-EN2-\n%s\n-EN2-\n", r)
				if err != nil {
					return fmt.Errorf("Error during enable (2): %s", err.Error())
				}
			}
			//p.Log("ENABLED")
		}
	}
	return nil
}
