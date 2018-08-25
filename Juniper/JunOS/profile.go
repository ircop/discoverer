package JunOS

import (
	"github.com/ircop/discoverer/base"
	"github.com/ircop/remote-cli"
	"fmt"
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
	if cli != nil {
		p.Cli.GlobalTimeout(60)
		p.Cli.SetPrompt(`(?msi:^(({master(?::\d+)}\n)?\S+(>|#))\s*$)`)
		if err := p.Cli.RegisterErrorPattern(`(is ambiguous\.|syntax error, expecting|unknown command)`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
			return err
		}

		_, err = p.Cli.Cmd("set cli screen-length 0")
		if err != nil {
			return fmt.Errorf("Cannot disable paging: %s", err.Error())
		}
	}
	return nil
}

