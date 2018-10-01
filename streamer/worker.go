package streamer

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/ircop/dproto"

	"github.com/ircop/discoverer/logger"
	"github.com/ircop/discoverer/profiles/CiscoIOS"
	"github.com/ircop/discoverer/profiles/DLinkDGS3100"
	"github.com/ircop/discoverer/profiles/DLinkDxS"
	"github.com/ircop/discoverer/profiles/EltexMES"
	"github.com/ircop/discoverer/profiles/Huawei/SW"
	"github.com/ircop/discoverer/profiles/Juniper/JunOS"
	"github.com/ircop/discoverer/profiles/Mikrotik/RouterOS"
	"github.com/ircop/discoverer/profiles/base"
	"github.com/ircop/remote-cli"
	nats "github.com/nats-io/go-nats-streaming"
	"runtime/debug"
	"sync"
)

var sendlock sync.Mutex

// todo: dispatch messages based on the packet type (box/perioric/etc.)
func workerCallback(msg *nats.Msg, chanReplies string) {
	// recover on top of all our jobs
	defer func() {
		if r := recover(); r != nil {
			logger.Panic("Recovered in nats worker callback: %+v\ntrace:\n%s\n", r, debug.Stack())
		}
	}()
	defer msg.Ack()

	logger.Debug("NATS worker got message")

	// read dpacket ; read packet type ; run box if this is box
	var packet dproto.DPacket
	err := proto.Unmarshal(msg.Data, &packet)
	if err != nil {
		logger.Err("Failed to parse dproto packet: %s", err.Error())
		return
	}

	if packet.PacketType != dproto.PacketType_BOX_REQUEST {
		logger.Err("Unsuppoprted packet type")
		return
	}

	// unmarshal payload into box-request packet
	var req dproto.BoxRequest
	if err = proto.Unmarshal(packet.Payload.Value, &req); err != nil {
		logger.Err("Failed to unmarshal box request: %s", err.Error())
		return
	}

	RequestID := req.RequestID
	var cli *remote_cli.Cli
	if req.Proto == dproto.Protocol_SSH {
		cli = remote_cli.New(remote_cli.CliTypeSsh, req.Host, int(req.Port), req.Login, req.Password, ``, int(req.Timeout))
	} else {
		cli = remote_cli.New(remote_cli.CliTypeTelnet, req.Host, int(req.Port), req.Login, req.Password, ``, int(req.Timeout))
	}


	var sw discoverer.Profile
	switch req.Profile {
	case dproto.ProfileType_DXS:
		sw = &DLinkDxS.Profile{}
		break
	case dproto.ProfileType_DGS3100:
		sw = &DLinkDGS3100.Profile{}
		break
	case dproto.ProfileType_IOS:
		sw = &CiscoIOS.Profile{}
		break
	case dproto.ProfileType_HUA:
		sw = &HuaweiSW.Profile{}
		break
	case dproto.ProfileType_JUNOS:
		sw = &JunOS.Profile{}
		break
	case dproto.ProfileType_MES:
		sw = &EltexMES.Profile{}
		break
	case dproto.ProfileType_ROUTEROS:
		sw = &RouterOS.Profile{}
		break
	default:
		logger.Err("Failed to map device profile: '%v'", req.Profile.String())
		sendError(conn, chanReplies, RequestID, fmt.Sprintf("Failed to map device profile: %+#v", req.Profile.String()))
		return
	}

	sw.SetLogger(logger.Log)
	sw.SetDebugLogger(logger.Debug)

	logger.Log("Starting box discovery for %s...", req.Host)
	if err = sw.Init(cli, req.Enable, ""); err != nil {
		// todo
		sendError(conn, chanReplies, RequestID, err.Error())
		logger.Err("Failed to init profile: %s (%s)", err.Error(), req.Host)
		return
	}
	defer sw.Disconnect()

	response := dproto.BoxResponse{
		Errors: make(map[string]string, 0),
		ReplyID:RequestID,
	}

	platform, err := sw.GetPlatform()
	response.Platform = &platform
	if err != nil {
		response.Errors[dproto.TaskType_PLATFORM.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_PLATFORM.String(), err.Error())
	}
	interfaces, err := sw.GetInterfaces()
	response.Interfaces = interfaces
	if err != nil {
		response.Errors[dproto.TaskType_INTERFACES.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_INTERFACES.String(), err.Error())
	}
	ips, err := sw.GetIps()
	response.Ipifs = ips
	if err != nil {
		response.Errors[dproto.TaskType_IPS.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_IPS.String(), err.Error())
	}
	uplink, err := sw.GetUplink()
	response.Uplink = uplink
	if err != nil {
		response.Errors[dproto.TaskType_UPLINK.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_UPLINK.String(), err.Error())
	}
	lldp, err := sw.GetLldp()
	response.LldpNeighbors = lldp
	if err != nil {
		response.Errors[dproto.TaskType_LLDP.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_LLDP.String(), err.Error())
	}
	vlans, err := sw.GetVlans()
	response.Vlans = vlans
	if err != nil {
		response.Errors[dproto.TaskType_VLANS.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_VLANS.String(), err.Error())
	}
	config, err := sw.GetConfig()
	response.Config = config
	if err != nil {
		response.Errors[dproto.TaskType_CONFIG.String()] = err.Error()
		logger.Err("%s: %s: %s", req.Host, dproto.TaskType_CONFIG.String(), err.Error())
	}

	logger.Log("Done box for %s", req.Host)
	sendReply(conn, response, chanReplies)
	//logger.Debug("Should send reply: %+#v\n", response)
}
/*
func workerCallback2(msg *nats.Msg, chanReplies string) {
	// recover on top of all our jobs
	defer func() {
		if r := recover(); r != nil {
			logger.Panic("Recovered in nats worker callback: %+v\ntrace:\n%s\n", r, debug.Stack())
		}
	}()
	defer msg.Ack()

	logger.Debug("NATS worker got message")

	var task __dproto.TaskRequest
	err := proto.Unmarshal(msg.Data, &task)
	if err != nil {
		logger.Err("Cannot unmarshal nats task request: %s", err.Error())
		return
	}
	// debug task contents?

	RequestID := task.RequestID
	// create and init device profile first
	var cli *remote_cli.Cli
	if task.Proto == __dproto.Protocol_SSH {
		cli = remote_cli.New(remote_cli.CliTypeSsh, task.Host, int(task.Port), task.Login, task.Password, ``, int(task.Timeout))
	} else {
		cli = remote_cli.New(remote_cli.CliTypeTelnet, task.Host, int(task.Port), task.Login, task.Password, ``, int(task.Timeout))
	}

	//sendError(conn, chanReplies, RequestID, fmt.Sprintf("Failed to map device profile: %+#v", task.Profile.String()))
	//return

	var sw discoverer.Profile
	switch task.Profile {
	case __dproto.ProfileType_DXS:
		sw = &DLinkDxS.Profile{}
		break
	case __dproto.ProfileType_DGS3100:
		sw = &DLinkDGS3100.Profile{}
		break
	case __dproto.ProfileType_IOS:
		sw = &CiscoIOS.Profile{}
		break
	case __dproto.ProfileType_HUA:
		sw = &HuaweiSW.Profile{}
		break
	case __dproto.ProfileType_JUNOS:
		sw = &JunOS.Profile{}
		break
	case __dproto.ProfileType_MES:
		sw = &EltexMES.Profile{}
		break
	case __dproto.ProfileType_ROUTEROS:
		sw = &RouterOS.Profile{}
		break
	default:
		logger.Err("Failed to map device profile: '%v'", task.Profile.String())
		sendError(conn, chanReplies, RequestID, fmt.Sprintf("Failed to map device profile: %+#v", task.Profile.String()))
		return
	}

	sw.SetLogger(logger.Log)
	sw.SetDebugLogger(logger.Debug)

	logger.Log("Starting %s for %s...", task.Type.String(), task.Host)
	if err = sw.Init(cli, task.Enable, ""); err != nil {
		sendError(conn, chanReplies, RequestID, err.Error())
		logger.Err("Failed to init profile: %s (%s)", err.Error(), task.Host)
		return
	}
	defer sw.Disconnect()


	// Profile is ready, now run tasks
	response := __dproto.Response{
		Type:task.Type,
		Errors:make(map[string]string,0),
		ReplyID:RequestID,
	}


	// if this is single task, we return task + error
	// if this is 'all', we should return set of errors in some way...
	if task.Type == __dproto.PacketType_PLATFORM || task.Type == __dproto.PacketType_ALL {
		platform, err := sw.GetPlatform()
		if err != nil {
			response.Errors[__dproto.PacketType_PLATFORM.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Platform = &platform
	}
	if task.Type == __dproto.PacketType_CONFIG || task.Type == __dproto.PacketType_ALL {
		config, err := sw.GetConfig()
		if err != nil {
			response.Errors[__dproto.PacketType_CONFIG.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Config = config
	}
	if task.Type == __dproto.PacketType_INTERFACES || task.Type == __dproto.PacketType_ALL {
		interfaces, err := sw.GetInterfaces()
		if err != nil {
			response.Errors[__dproto.PacketType_INTERFACES.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Interfaces = interfaces
	}
	if task.Type == __dproto.PacketType_IPS || task.Type == __dproto.PacketType_ALL {
		ipifs, err := sw.GetIps()
		if err != nil {
			response.Errors[__dproto.PacketType_IPS.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Ipifs = ipifs
	}
	if task.Type == __dproto.PacketType_LLDP || task.Type == __dproto.PacketType_ALL {
		lldp, err := sw.GetLldp()
		if err != nil {
			response.Errors[__dproto.PacketType_LLDP.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.LldpNeighbors = lldp
	}
	if task.Type == __dproto.PacketType_UPLINK || task.Type == __dproto.PacketType_ALL {
		up, err := sw.GetUplink()
		if err != nil {
			response.Errors[__dproto.PacketType_UPLINK.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Uplink = up
	}
	if task.Type == __dproto.PacketType_VLANS || task.Type == __dproto.PacketType_ALL {
		vlans, err := sw.GetVlans()
		if err != nil {
			response.Errors[__dproto.PacketType_VLANS.String()] = err.Error()
			logger.Err("%s: %s: %s", task.Host, task.Type.String(), err.Error())
		}
		response.Vlans = vlans
	}

	logger.Log("Done %s for %s", task.Type.String(), task.Host)
	sendReply(conn, response, chanReplies)
	logger.Debug("Should send reply: %+#v\n", response)
}
*/
func sendReply(conn nats.Conn, response dproto.BoxResponse, topic string) {
	logger.Debug("- sending reply... -")
	bs, err := proto.Marshal(&response)
	if err != nil {
		logger.Err("Cannot marshal response for request: %s", err.Error())
		return
	}

	packet := dproto.DPacket{
		PacketType:dproto.PacketType_BOX_REPLY,
		Payload:&any.Any{
			TypeUrl:dproto.PacketType_BOX_REPLY.String(),
			Value:bs,
		},
	}
	packetBts, err := proto.Marshal(&packet)
	if err != nil {
		logger.Err("Cannot marshal dproto packet: %s", err.Error())
		return
	}

	_, err = conn.PublishAsync(topic, packetBts, func(lguid string, errX error) {
		if errX == nil {
			logger.Debug("Got ack for reply %s", response.ReplyID)
		} else {
			logger.Err("Error sending reply id %s: %s", response.ReplyID, err.Error())
		}
	})
	if err != nil {
		logger.Err("Cannot send response for request %s: %s", response.ReplyID, err.Error())
	}
}

func sendError(conn nats.Conn, topic string, reply string, message string) {
	msg := dproto.BoxResponse{
		Error:message,
		ReplyID:reply,
	}
	bs, err := proto.Marshal(&msg)
	if err != nil {
		logger.Err("Cannot marshall nats error message (req.): %s", reply, err.Error())
		return
	}

	packet := dproto.DPacket{
		PacketType:dproto.PacketType_BOX_REPLY,
		Payload:&any.Any{
			TypeUrl:dproto.PacketType_BOX_REPLY.String(),
			Value:bs,
		},
	}
	packetBts, err := proto.Marshal(&packet)
	if err != nil {
		logger.Err("Cannot marshal dproto packet: %s", err.Error())
		return
	}

	_, err = conn.PublishAsync(topic, packetBts, func(lguid string, errX error) {
		if errX == nil {
			logger.Debug("Got ack for reply %s", msg.ReplyID)
		} else {
			logger.Err("Error sending reply id %s: %s", msg.ReplyID, err.Error())
		}
	})
	if err != nil {
		logger.Err("Cannot publish nats error message (req.%s): %s", reply, err.Error())
	}
}
