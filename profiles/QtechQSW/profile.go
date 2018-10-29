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
		p.Cli.SetPrompt(`(?msi:^(?P<hostname>[a-zA-Z0-9]\S{0,19})(?:[\.\-_\d\w]+)?(?:\(config[^\)]*\))?#)`)
		p.SetPrompt()
		p.Cli.GlobalTimeout(60)

		if err := p.Cli.RegisterErrorPattern(`(Unrecognized |Invalid|Invalid input|Ambiguous)`, "Syntax error"); err != nil {
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
