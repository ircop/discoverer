package streamer

import (
	"github.com/ircop/discoverer/logger"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"time"
	"fmt"
)

var conn *nats.Conn

// Run streaming nats worker
func Run(natsURL string, chanTasks string, chanReplies string) error {
	logger.Log("Starting NATS worker")

	var err error
	conn, err = nats.Connect(natsURL, nats.ReconnectWait(time.Second * 2),
			nats.DisconnectHandler(func(nc *nats.Conn) {
				logger.Err("NATS server disconnected")
			}),
			nats.ReconnectHandler(func(nc *nats.Conn) {
				logger.Log("NATS server got reconnected")
			}),
			nats.ClosedHandler(func(nc *nats.Conn) {
				logger.Err("NATS server connection closed: %s", nc.LastError())
			}),
		)
	if err != nil {
		return err
	}
	//defer conn.Close()

	// subscribe
	_, err = conn.QueueSubscribe(chanTasks, chanTasks, func(msg *nats.Msg) {
		fmt.Printf(" - GOT REQUEST -\n")
		go workerCallback(msg, chanReplies)
	})
	if err != nil {
		return errors.Wrap(err, "Cannot subsctibe to NATS tasks channel")
	}

	return nil
}
