package QtechQSW

import (
	"github.com/ircop/discoverer/profiles/base"
	"github.com/ircop/remote-cli"
)

type Profile struct {
	discoverer.Generic
}

func (p *Profile) Init(cli remote_cli.CliInterface, enable string, community string) error {
	err := p.InitShared(cli, enable, community)
	if err != nil {
		return err
	}

	if cli != nil {
		p.Cli.SetLoginPrompt(`([Uu]ser(\s)?[Nn]ame\:(\s+)|login:)?$`)
		p.SetPrompt()
		p.Cli.GlobalTimeout(60)

		if err := p.Cli.RegisterErrorPattern(`(% Unrecognized | % Invalid|% Ambiguous|% Incomplete|% Unknown)`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
			return err
		}

		if err := p.Cli.RegisterCallback(`^(All|More)[^\n]+(<return>|<ctrl>\+z)(\s+)?`, func() { p.Cli.WriteRaw([]byte{' '}) }); err != nil {
			return err
		}
	}

	return nil
}
