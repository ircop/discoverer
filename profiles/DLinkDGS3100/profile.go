package DLinkDGS3100

import (
	"github.com/ircop/discoverer/profiles/base"
	"github.com/ircop/remote-cli"
)

// Profile instance
type Profile struct {
	discoverer.Generic
}


// SetPrompt sets CLI prompt for current profile (only if CLI is active)
func (p *Profile) SetPrompt() {
	if p.Cli != nil {
		//p.Cli.SetPrompt(`(?msi:[\$%#>]$)`)
		p.Cli.SetPrompt(`([a-zA-Z0-9\-]+?)#(\s)?$`)
		//p.Cli.SetPrompt(`(?msi:^D(E|G|X)S.*#$)`)
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
		p.Cli.DlinkPagination()

		if err := p.Cli.RegisterErrorPattern(`(Command: |Incomplete|Available commands|Next possible completions|Ambiguous token)`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil{
			return err
		}

		// Enable admin first if needed
		/*if enable != "" {
			p.Cli.Write([]byte("enable admin"))
			r, err := p.Cli.ReadUntil(`([Pp]ass[Ww]ord:$)`)
			if err != nil {
				//fmt.Printf(r)
				if !strings.Contains(r, "already have") {
					return fmt.Errorf("Cannot 'enable admin': %s", err.Error())
				}
			}

			r, err = p.Cli.Cmd(enable)
			if err != nil {
				return fmt.Errorf("Cannot 'enable admin': %s", err.Error())
			}
		}*/
	}

	return nil
}
