package discoverer

import (
	"regexp"
)

// CliCallback is a struct, handling regexps/callbacks for CLI iteractions
type CliCallback struct {
	re			regexp.Regexp
	callback	func()
}

// CliCallbacks dummy function, used for steting CLI callbacks
func (p *Generic) SetCallbacks() error {
	// nothing here int dummy
	return nil
}

// SetPrompt sets cli prompt after initialization
func (p *Generic) SetPrompt() {
	// nothing here int dummy
}

// Command sends your data to remote via CLI.
// It is wrapper around github.com/ircop/remote-cli, used for CLI initialization
func (p *Generic) Command(cmd string) (string, error) {
	return "", nil
}
