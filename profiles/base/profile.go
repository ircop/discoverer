package discoverer

import (
	"github.com/ircop/dproto"
	"github.com/ircop/remote-cli"
	"fmt"
)

// Profile interface, will be implemented in per-OS profiles
type Profile interface {
	Init(remote_cli.CliInterface, string, string) error
	InitShared(remote_cli.CliInterface, string, string) error
	SetCallbacks() (error)
	SetPrompt()
	SetAddr(string)

	GetPlatform() (dproto.Platform, error)                // dlink|cisco|mes|3100|jun|hua
	GetInterfaces() (map[string]*dproto.Interface, error) // dlink|cisco|mes|3100|jun|hua
	GetLldp() ([]*dproto.LldpNeighbor, error)             // dlink|cisco|mes|3100|jun|hua
	GetVlans() ([]*dproto.Vlan, error)                    // dlink|cisco|mes|3100|jun|hua
	GetIps() ([]*dproto.Ipif, error)                      // dlink|cisco|mes|3100|jun|hua
	GetUplink() (string, error)                             // dlink|cisco|mes|3100|jun|hua-
	GetConfig() (string, error)                             // dlink|cisco|mes|3100|jun|hua

	SetLogger(func(string, ...interface{}))
	SetDebugLogger(func(string, ...interface{}))

	Disconnect()
}

var ErrNotImplemented = fmt.Errorf("Method not implemented")

// Generic profile realization. Used for dummy functions, like 'not implemented, sorry'
// Cli in connected state should be passed in Init
type Generic struct {
	Profile
	Addr			string
	//device			discoverer.DeviceProfile
	Cli				remote_cli.CliInterface
	Community		string

	logger			func(string, ...interface{})
	loggerSet		bool
	loggerDebug		func(string, ...interface{})
	loggerDebugSet	bool

	// Model may be needed for some functions
	Model			string
	// Version is currently running firmware
	Version			string
	// Enable passwors
	EnablePassword	string
	enabled			bool
}

// Init parses device profile contents, stores them, checks them.
// 'enable' is enable password ; 'community' is community string - both are optional
func (p *Generic) InitShared(cli remote_cli.CliInterface, enable string, community string) error {
	//p.device = device
	p.Cli = cli
	p.EnablePassword = enable
	p.enabled = false
	p.Community = community

	if p.Cli == nil && community == "" {
		return fmt.Errorf("Both CLI type and SNMP community are not set!")
	}


	return nil
}

// SetEnable sets enable password
func (p *Generic) SetEnable(pw string) {
	p.EnablePassword = pw
}

// Disconnect cli
func (p *Generic) Disconnect() {
	if p.Cli != nil {
		p.Cli.Close()
	}
}

// Init dummy func
// 'enable' is enable password ; 'community' is community string - both are optional
func (p *Generic) Init(cli remote_cli.CliInterface, enable string, community string) error {
	return p.InitShared(cli, enable, community)
}

func (p *Generic) SetLogger(cb func(string, ...interface{})) {
	p.loggerSet = true
	p.logger = cb
}

func (p *Generic) SetDebugLogger(cb func(string, ...interface{})) {
	p.loggerDebugSet = true
	p.loggerDebug = cb
}

// Log writes normal log (via callback)
func (p *Generic) Log(msg string, args ...interface{}) {
	if p.loggerSet {
		//msg = strings.Replace(msg, "%", "%%", -1)
		p.logger(msg, args...)
	}
}

// Debug writes debug log (via callback)
func (p *Generic) Debug(msg string, args ...interface{}) {
	if p.loggerDebugSet {
		//msg = strings.Replace(msg, "%", "%%", -1)
		p.loggerDebug(msg, args...)
	}
}

func (p *Generic) SetAddr(addr string) {
	p.Addr = addr
}
