package DLinkDxS

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
		p.Cli.SetPrompt(`(:[a-zA-Z0-9]|:([a-zA-Z0-9]+?))#$`)
		p.Cli.GlobalTimeout(60)
		p.Cli.DlinkPagination()

		if err := p.Cli.RegisterErrorPattern(`(already have|ecord drive name|Available commands|Next possible completions|Ambiguous token)\:`, "Syntax error"); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
			return err
		}
		p.Log("Connected")

		// Enable admin first if needed
		if enable != "" {
			p.Cli.Write([]byte("enable admin"))
			p.Log("sent enable admin, reading until password...")
			r, err := p.Cli.ReadUntil(`((:[a-zA-Z0-9]|:([a-zA-Z0-9]+?))#$|([Pp]ass[Ww]ord:$))`)
			if err != nil {
				p.Log("got err")
				//fmt.Printf(r)
				if !strings.Contains(r, "already have") && !strings.Contains(r, "commands:") {
					return fmt.Errorf("Cannot 'enable admin': %s", err.Error())
				}
			} else {
				p.Log(r)
				if !strings.Contains(strings.ToLower(r), "already have") && !strings.Contains(strings.ToLower(r),"next possible") &&
					!strings.Contains(r, "commands:") && !strings.Contains(strings.ToLower(r), "syntax error") {
						r, err = p.Cli.Cmd(enable)
						if err != nil {
							return fmt.Errorf("Cannot 'enable admin' (2): %s", err.Error())
						}
				}
			}
		}
	}
	//p.Log("p.Init: done")
	return nil
}
