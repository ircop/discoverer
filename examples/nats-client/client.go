package main

import (
	"sync"
	"github.com/ircop/discoverer/dproto"
	"time"
	"fmt"
	nats "github.com/nats-io/go-nats-streaming"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

var TasksPool sync.Map
type WaitingTask struct {
	RequestID string
	Type      __dproto.PacketType
	Timer     *time.Timer
}

var glock sync.Mutex

func main() {
	fmt.Println("Starting go-nats client")

	natsConn, err := nats.Connect("test-cluster", "example1", nats.NatsURL("nats://discoverer:password@127.0.0.1"))
	if err != nil {
		panic(err)
	}
	defer natsConn.Close()

	fmt.Println("Connected to nats streaming server")

	_, err = natsConn.Subscribe("replies", func(msg *nats.Msg) {
			go handleResponse(msg)
		},
		nats.DurableName("replies"),
		nats.MaxInflight(200),			// this is how mutch THIS CLIENT can handle one-time events
		nats.SetManualAckMode(),
		nats.AckWait(time.Minute * 15),
	)
	if err != nil {
		fmt.Printf("Cannot subscribe to responses channel: %s\n", err.Error())
		return
	}

	go sendTask(natsConn, __dproto.PacketType_PLATFORM, "10.170.3.99", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_INTERFACES, "10.170.3.99", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_PLATFORM, "10.10.10.248", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_IPS, "10.10.10.248", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_INTERFACES, "10.10.10.150", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_PLATFORM, "10.170.3.99", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_INTERFACES, "10.170.3.99", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_PLATFORM, "10.10.10.248", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_IPS, "10.10.10.248", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)
	go sendTask(natsConn, __dproto.PacketType_INTERFACES, "10.10.10.150", __dproto.Protocol_TELNET, __dproto.ProfileType_DXS)


	select{}
}

func handleResponse(msg *nats.Msg) {
	defer msg.Ack()

	var response __dproto.Response
	err := proto.Unmarshal(msg.Data, &response)
	if err != nil {
		fmt.Printf("Cannot unmarshal response: %s\n", err.Error())
		return
	}

	id := response.ReplyID
	waitingTask, ok := TasksPool.Load(id)
	if !ok {
		fmt.Printf("Ignoring unknown reply id '%s'\n", id)
		return
	}
	wt := waitingTask.(*WaitingTask)
	wt.Timer.Stop()

	TasksPool.Delete(id)

	if response.Type == __dproto.PacketType_ERROR {
		fmt.Printf("Error (%s): %s\n", id, response.Error)
		return
	}

	switch response.Type {
	case __dproto.PacketType_PLATFORM:
		if err, ok := response.Errors[__dproto.PacketType_PLATFORM.String()]; ok {
			fmt.Printf("Platform request resulted with error: %s", err)
		}
		fmt.Printf("Platform: %+#v\n", response.Platform)
		break
	case __dproto.PacketType_INTERFACES:
		if err, ok := response.Errors[__dproto.PacketType_INTERFACES.String()]; ok {
			fmt.Printf("Interfaces request resulted with error: %s", err)
		}
		//for _, iface := range response.Interfaces {
		//	fmt.Printf("iface: %+#v\n", iface)
		//}
		fmt.Printf("got %d ifaces\n", len(response.Interfaces))
		break
	case __dproto.PacketType_CONFIG:
		if err, ok := response.Errors[__dproto.PacketType_CONFIG.String()]; ok {
			fmt.Printf("Config request resulted with error: %s", err)
		}
		fmt.Printf("Config: %+#v\n", response.Config)
		break
	case __dproto.PacketType_IPS:
		if err, ok := response.Errors[__dproto.PacketType_IPS.String()]; ok {
			fmt.Printf("Ips request resulted with error: %s", err)
		}
		//for _, ipif := range response.Ipifs {
		//	fmt.Printf("ipif: %+#v\n", ipif)
		//}
		fmt.Printf("got %d ipifs\n", len(response.Ipifs))
		break
	case __dproto.PacketType_LLDP:
		if err, ok := response.Errors[__dproto.PacketType_LLDP.String()]; ok {
			fmt.Printf("LLDP request resulted with error: %s", err)
		}
		//for _, lldp := range response.LldpNeighbors {
		//	fmt.Printf("lldp: %+#v\n", lldp)
		//}
		fmt.Printf("got %d neighbors\n", len(response.LldpNeighbors))
		break
	case __dproto.PacketType_VLANS:
		if err, ok := response.Errors[__dproto.PacketType_VLANS.String()]; ok {
			fmt.Printf("Vlan request resulted with error: %s", err)
		}
		for _, v := range response.Vlans {
			fmt.Printf("vlan: %+#v\n", v)
		}
		break
	case __dproto.PacketType_UPLINK:
		if err, ok := response.Errors[__dproto.PacketType_UPLINK.String()]; ok {
			fmt.Printf("Uplink request resulted with error: %s", err)
		}
		fmt.Printf("Uplink: '%s'\n", response.Uplink)
		break
	default:
		fmt.Printf("Error: unknown response type!")
		return
	}
}

//func sendTask(conn *nats.Conn, taskType discoverproto.TaskType, host string, protocol discoverproto.Proto, profile discoverproto.ProfileType) {
func sendTask(conn nats.Conn, taskType __dproto.PacketType, host string, protocol __dproto.Protocol, profile __dproto.ProfileType) {
	id, _ := uuid.NewRandom()
	port := 22
	if protocol == __dproto.Protocol_TELNET {
		port = 23
	}

	message := __dproto.TaskRequest{
		RequestID:id.String(),
		Timeout: 60,
		Login:"login",
		Password:"pw",
		Profile:profile,
		Type:taskType,
		Enable:"",
		Proto:protocol,
		Host:host,
		Port:int32(port),
	}

	bs, err := proto.Marshal(&message)
	if err != nil {
		fmt.Printf("Error marshaling message: %s", err.Error())
		return
	}

	// Put this task into waiting pool.
	// Set timeout function, that will remove this task from waiting pool
	wt := WaitingTask{
		RequestID:id.String(),
		Type:taskType,
	}
	wt.Timer = time.AfterFunc(time.Minute * 15, func(){
		fmt.Printf("Removing waiting task '%s' from pool due to timeout\n", id.String())
		TasksPool.Delete(id.String())
	})
	TasksPool.Store(id.String(), &wt)


	// send
	glock.Lock()
	guid, err := conn.PublishAsync("tasks", bs, asyncCallback)
	if err != nil {
		fmt.Printf("error publishing task: %s\n", err.Error())
	}
	fmt.Printf("Published [%s]\n", guid)
	defer glock.Unlock()
}

func asyncCallback(lguid string, err error) {
	fmt.Printf("Got ACK for guid %s\n", lguid)
	if err != nil {
		fmt.Printf("Error in ack: %s\n", err.Error())
	}
}
