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

func (p *Profile) Init(cli *remote_cli.Cli, enable string, community string) error {
	err := p.InitShared(cli, enable, community)
	if err != nil {
		return err
	}

	p.Cli.GlobalTimeout(60)

	if err := p.Cli.RegisterErrorPattern(`(is ambiguous\.|syntax error, expecting|unknown command)`, "Syntax error"); err != nil {
		return err
	}

	_, err = p.Cli.Cmd("set cli screen-length 0")
	if err != nil {
		return fmt.Errorf("Cannot disable paging: %s", err.Error())
	}
	return nil
}

