package main

import (
	"flag"
	"fmt"
	"github.com/ircop/discoverer/cfg"
	"github.com/ircop/discoverer/logger"
	"github.com/ircop/discoverer/streamer"
)

func main() {
	configPath := flag.String("c", "./discoverer.toml", "Config file location")
	flag.Parse()

	config, err := cfg.NewCfg(*configPath)
	if err != nil {
		fmt.Printf("[FATAL]: Cannot read config: %s\n", err.Error())
		return
	}

	// log
	if err := logger.InitLogger(config.LogDebug, config.LogDir); err != nil {
		fmt.Printf("[FATAL]: %s", err.Error())
		return
	}

	// run nats and/or rest daemons
	// nats first:
	if err := streamer.Run(config.NatsURL, config.NatsTasks, config.NatsReplies); err != nil {
		logger.Err("Cannot run nats listener: %s", err.Error())
		return
	}

	select{}
}
