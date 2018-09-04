package RouterOS

import (
	"github.com/ircop/discoverer/profiles/base"
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

	if p.Cli != nil {
		p.Cli.SetLogin(fmt.Sprintf("%s+ct255w255hoffc", p.Cli.GetLogin()))
		p.Cli.GlobalTimeout(60)
		p.Cli.SetPrompt(`\[(?P<prompt>[^\]@]+@.+?)\] >\s+$`)
		p.Cli.SetLoginPrompt(`^Login\:$`)
		p.Cli.SetPasswordPrompt(`^ ?[Pp]ass[Ww]ord: ?$`)

		if err := p.Cli.RegisterErrorPattern(`(bad command|unknown command)`, "Syntax error"); err != nil {
			return err
		}

		if err = p.Cli.Connect(); err != nil {
			return err
		}
	}

	return nil
}

