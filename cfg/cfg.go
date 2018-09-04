package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"net"
)

// Cfg is struct for handling config parameters
type Cfg struct {
	Nats			bool
	NatsURL			string
	NatsTasks		string
	NatsReplies		string

	RPC				bool
	RPCPort			int64
	RPCHost			string
	RPCSsl			bool
	RPCCert			string
	RPCKey			string

	LogDir			string
	LogDebug		bool
}

// NewCfg reads config with given path
func NewCfg(path string) (*Cfg, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := new(Cfg)

	c.Nats = viper.GetBool("nats.nats")
	c.NatsURL = viper.GetString("nats.url")
	c.NatsTasks = viper.GetString("nats.tasks-chan")
	c.NatsReplies = viper.GetString("nats.reply-chan")

	c.RPC = viper.GetBool("rpc.rpc")
	c.RPCPort = viper.GetInt64("rpc.listen-port")
	c.RPCHost = viper.GetString("rpc.listen-host")
	c.RPCSsl = viper.GetBool("rpc.ssl")
	c.RPCCert = viper.GetString("rpc.ssl-cert")
	c.RPCKey = viper.GetString("rpc.ssl-key")

	c.LogDir = viper.GetString("log.dir")
	c.LogDebug = viper.GetBool("log.debug")

	// perform some checks
	if !c.Nats && !c.RPC {
		return nil, fmt.Errorf("Nothing to do: both rpc and nats are disabled in config")
	}

	if c.Nats && (c.NatsURL == "" || c.NatsTasks == "" || c.NatsReplies == "") {
		return nil, fmt.Errorf("Wrong NATS configuration: url/tasks-chan/reply-chan not set")
	}

	if c.RPC {
		if c.RPCSsl && (c.RPCCert == "" || c.RPCKey == "") {
			return nil, fmt.Errorf("SSL is enabled, but cert/key not set")
		}
		if c.RPCPort == 0 || c.RPCHost == "" {
			return nil, fmt.Errorf("RPC listen host/port not set")
		}

		if net.ParseIP(c.RPCHost) == nil {
			return nil, fmt.Errorf("Wrong RPC host given")
		}
	}

	if c.LogDir == "" {
		return nil, fmt.Errorf("Log dir is not set")
	}

	return c, nil
}
/*
// CheckParams checking parameters for various things and returns error if something is wrong
func (c *Cfg) CheckParams() error { // nolint
	// validate ip address
	if !util.IsIP(c.ListenIP) {
		return fmt.Errorf("listen.ip is not ip address")
	}

	// Validate port
	if c.ListenPort < 1024 || c.ListenPort > 65535 {
		return fmt.Errorf("listen.port should be between 1024 and 65535")
	}

	// validate db params
	if "" == c.DBHost || "" == c.DBUser || "" == c.DBPassword || "" == c.DBName || 0 == c.DBPort {
		return fmt.Errorf("db parameters not set in config")
	}

	// validate ssl
	if c.Ssl && ("" == c.SslKey || "" == c.SslCert) {
		return fmt.Errorf("ssl enabled, but no cert/key provided")
	}

	if c.Workers == 0 {
		return fmt.Errorf("Workers count is not set (nms.workers)")
	}

	return nil
}
*/