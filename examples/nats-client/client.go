package main

import (
	"sync"
	"github.com/ircop/discoverer/dproto"
	"time"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

var TasksPool sync.Map
type WaitingTask struct {
	RequestID		string
	Type			dproto.PacketType
	Timer			*time.Timer
}

func main() {
	fmt.Println("Starting go-nats client")

	natsConn, err := nats.Connect("nats://discoverer:password@127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer natsConn.Close()

	fmt.Println("Connected to nats streaming server")

	_, err = natsConn.Subscribe("replies", func(msg *nats.Msg) {
		go handleResponse(msg)
	})
	if err != nil {
		fmt.Printf("Cannot subscribe to responses channel: %s\n", err.Error())
		return
	}


	sendTask(natsConn, dproto.PacketType_PLATFORM, "10.170.3.99", dproto.Protocol_TELNET, dproto.ProfileType_DXS)
	sendTask(natsConn, dproto.PacketType_INTERFACES, "10.170.3.99", dproto.Protocol_TELNET, dproto.ProfileType_DXS)
	sendTask(natsConn, dproto.PacketType_PLATFORM, "10.10.10.248", dproto.Protocol_TELNET, dproto.ProfileType_DXS)
	sendTask(natsConn, dproto.PacketType_IPS, "10.10.10.248", dproto.Protocol_TELNET, dproto.ProfileType_DXS)
	sendTask(natsConn, dproto.PacketType_INTERFACES, "10.10.10.150", dproto.Protocol_TELNET, dproto.ProfileType_DXS)


	select{}
}

func handleResponse(msg *nats.Msg) {
	fmt.Println(" - got response -")
	id := msg.Reply
	waitingTask, ok := TasksPool.Load(id)
	if !ok {
		fmt.Printf("Ignoring unknown reply id '%s'\n", id)
		return
	}
	wt := waitingTask.(*WaitingTask)
	wt.Timer.Stop()

	TasksPool.Delete(id)

	var response dproto.Response
	err := proto.Unmarshal(msg.Data, &response)
	if err != nil {
		fmt.Printf("Cannot unmarshal response %s: %s\n", id, err.Error())
		return
	}

	if response.Type == dproto.PacketType_ERROR {
		fmt.Printf("Error (%s): %s\n", id, response.Error)
		return
	}

	switch response.Type {
	case dproto.PacketType_PLATFORM:
		if err, ok := response.Errors[dproto.PacketType_PLATFORM.String()]; ok {
			fmt.Printf("Platform request resulted with error: %s", err)
		}
		fmt.Printf("Platform: %+#v\n", response.Platform)
		break
	case dproto.PacketType_INTERFACES:
		if err, ok := response.Errors[dproto.PacketType_INTERFACES.String()]; ok {
			fmt.Printf("Interfaces request resulted with error: %s", err)
		}
		for _, iface := range response.Interfaces {
			fmt.Printf("iface: %+#v\n", iface)
		}
		break
	case dproto.PacketType_CONFIG:
		if err, ok := response.Errors[dproto.PacketType_CONFIG.String()]; ok {
			fmt.Printf("Config request resulted with error: %s", err)
		}
		fmt.Printf("Config: %+#v\n", response.Config)
		break
	case dproto.PacketType_IPS:
		if err, ok := response.Errors[dproto.PacketType_IPS.String()]; ok {
			fmt.Printf("Ips request resulted with error: %s", err)
		}
		for _, ipif := range response.Ipifs {
			fmt.Printf("ipif: %+#v\n", ipif)
		}
		break
	case dproto.PacketType_LLDP:
		if err, ok := response.Errors[dproto.PacketType_LLDP.String()]; ok {
			fmt.Printf("LLDP request resulted with error: %s", err)
		}
		for _, lldp := range response.LldpNeighbors {
			fmt.Printf("lldp: %+#v\n", lldp)
		}
		break
	case dproto.PacketType_VLANS:
		if err, ok := response.Errors[dproto.PacketType_VLANS.String()]; ok {
			fmt.Printf("Vlan request resulted with error: %s", err)
		}
		for _, v := range response.Vlans {
			fmt.Printf("vlan: %+#v\n", v)
		}
		break
	case dproto.PacketType_UPLINK:
		if err, ok := response.Errors[dproto.PacketType_UPLINK.String()]; ok {
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
func sendTask(conn *nats.Conn, taskType dproto.PacketType, host string, protocol dproto.Protocol, profile dproto.ProfileType) {
	id, _ := uuid.NewRandom()
	port := 22
	if protocol == dproto.Protocol_TELNET {
		port = 23
	}

	message := dproto.TaskRequest{
		Timeout: 60,
		Login:"login",
		Password:"password",
		Profile:profile,
		Type:taskType,
		RequestID:id.String(),
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

	msg := nats.Msg{
		Data:bs,
		Subject:"tasks",
	}

	// Put this task into waiting pool.
	// Set timeout function, that will remove this task from waiting pool
	wt := WaitingTask{
		RequestID:id.String(),
		Type:taskType,
	}
	wt.Timer = time.AfterFunc(time.Second * 120, func(){
		fmt.Printf("Removing waiting task '%s' from pool due to timeout\n", id.String())
		TasksPool.Delete(id.String())
	})
	TasksPool.Store(id.String(), &wt)


	err = conn.PublishMsg(&msg)
	if err != nil {
		fmt.Printf("Error publishing msg: %s", err.Error())
		wt.Timer.Stop()
		TasksPool.Delete(id.String())
		return
	}

	fmt.Printf("Sent task request: %s (%s, %s)\n", id.String(), taskType.String(), host)
}

