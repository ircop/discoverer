package HuaweiSW

import (
	"fmt"
	"github.com/ircop/discoverer/base"
	"github.com/ircop/remote-cli"
	"regexp"
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
		p.Cli.GlobalTimeout(60)
		//p.Cli.SetPrompt(`(?msi:(?P<hostname>\S+)(?:\(.*)?(#|]|>)$)`)
		p.Cli.SetPrompt(`(?m:^[<#\[](?P<hostname>[a-zA-Z0-9-_\\\.\[\(/'\"\|\s:,=]+)(?:-[a-zA-Z0-9/\_]+)*[>#\]\)])`)

		if err := p.Cli.RegisterErrorPattern(`(Error: .*\^)`, "Syntax error"); err != nil {
			return err
		}
		// confirmation of long-time output
		if err = p.Cli.RegisterCallback(`long time to ex.+\[Y/N]:`, func() {
			p.Cli.Write([]byte{'Y'})
		}); err != nil {
			return err
		}
		// pagination
		if err = p.Cli.RegisterCallback(`(?msi:^\s+ ---- More ---)`, func() {
			p.Cli.WriteRaw([]byte{' '})
		}); err != nil {
			return err
		}

		if err := p.Cli.Connect(); err != nil {
			return err
		}


		// this is so stupid :( Hua prompt is <xx> or [xx], AND some hua ouput is [someshit] :(
		re, err := regexp.Compile(`(?msi:^sysname\s+(?P<sysname>[^\s\n]+))(\s+)?\n`)
		if err != nil {
			return fmt.Errorf("Cannot compile sysname regex; %s", err.Error())
		}
		res, err := p.Cli.Cmd("display current-configuration configuration system  | in sysname")
		if err != nil {
			return fmt.Errorf("Cannot find sysname: %s", err.Error())
		}
		p.Debug(res)

		out := p.ParseSingle(re, res)
		sysname := out["sysname"]
		sysname = strings.Replace(sysname, "-", `\-`, -1)
		sysname = strings.Replace(sysname, ".", `\.`, -1)
		sysname = strings.Replace(sysname, "$", `\$`, -1)
		sysname = strings.Replace(sysname, "^", `\^`, -1)
		sysname = strings.Replace(sysname, "%", `\%`, -1)
		prompt := fmt.Sprintf(`(?m:^(\[|<)%s(>|]))`, sysname)
		p.Log("Setting prompt to '%s'", prompt)
		p.Cli.SetPrompt(prompt)
	}
	return nil
}

