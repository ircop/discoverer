package EltexMES

import (
	"github.com/ircop/discoverer/base"
	"github.com/ircop/remote-cli"
)

// Profile instance
type Profile struct {
	discoverer.Generic
}


// SetPrompt sets CLI prompt for current profile (only if CLI is active)
func (p *Profile) SetPrompt() {
	if p.Cli != nil {
		p.Cli.SetPrompt(`(?msi:^[\.a-zA-Z0-9\-_]+[\$%#>](\s+)?$)`)
	}
}

func (p *Profile) Init(cli remote_cli.CliInterface, enable string, community string) error {
	err := p.InitShared(cli, enable, community)
	if err != nil {
		return err
	}

	if cli != nil {
		p.SetPrompt()
		p.Cli.GlobalTimeout(60)

		if err := p.Cli.RegisterErrorPattern(`(% Unrecognized | % Invalid|% Ambiguous|% Incomplete|% Unknown)`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
			return err
		}

		//if err := p.Cli.RegisterCallback(`^More[^\n]+<return>(\s+)?\n`, func() { p.Cli.WriteRaw([]byte{' '}) }); err != nil {
		if err := p.Cli.RegisterCallback(`^(All|More)[^\n]+(<return>|<ctrl>\+z)(\s+)?`, func() { p.Cli.WriteRaw([]byte{' '}) }); err != nil {
			return err
		}
	}

	return nil
}
