package streamer

import (
	"fmt"
	"github.com/ircop/discoverer/logger"
	nats "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

var conn nats.Conn

// Run streaming nats worker
func Run(natsURL string, chanTasks string, chanReplies string) error {
	logger.Log("Starting NATS worker")

	var err error

	hostname, err := os.Hostname()
	hostname = strings.Replace(hostname, ".", "-", -1)
	if err != nil {
		return fmt.Errorf("Cannot discover hostname: %s", err.Error())
	}

	// todo: set uniq discoverer ID
	conn, err = nats.Connect("test-cluster", hostname, nats.NatsURL(natsURL))
	if err != nil {
		return err
	}
	//defer conn.Close()

	// subscribe
	_, err = conn.QueueSubscribe(chanTasks, "", func(msg *nats.Msg) {
			//fmt.Printf(" - GOT REQUEST -\n")
			go workerCallback(msg, chanReplies)
		},
		//nats.DurableName("tasks"),
		nats.MaxInflight(100),				// this is how mutch THIS WORKER can handle one-time events
		nats.SetManualAckMode(),
		nats.AckWait(time.Minute * 15),
	)
	if err != nil {
		return errors.Wrap(err, "Cannot subsctibe to NATS tasks channel")
	}
	logger.Debug("Subscribed to '%s' channel", chanTasks)

	return nil
}
