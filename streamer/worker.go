package streamer

import (
	"github.com/golang/protobuf/proto"
	"github.com/ircop/discoverer/dproto"
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
	"github.com/nats-io/go-nats"
	"runtime/debug"
)

func workerCallback(msg *nats.Msg, chanReplies string) {
	// recover on top of all our jobs
	defer func() {
		if r := recover(); r != nil {
			logger.Panic("Recovered in nats worker callback: %+v\ntrace:\n%s\n", r, debug.Stack())
		}
	}()

	logger.Debug("NATS worker got message")

	var task dproto.TaskRequest
	err := proto.Unmarshal(msg.Data, &task)
	if err != nil {
		logger.Err("Cannot unmarshal nats task request: %s", err.Error())
		return
	}
	// debug task contents?

	RequestID := task.RequestID
	// create and init device profile first
	var cli *remote_cli.Cli
	if task.Proto == dproto.Protocol_SSH {
		cli = remote_cli.New(remote_cli.CliTypeSsh, task.Host, int(task.Port), task.Login, task.Password, ``, int(task.Timeout))
	} else {
		cli = remote_cli.New(remote_cli.CliTypeTelnet, task.Host, int(task.Port), task.Login, task.Password, ``, int(task.Timeout))
	}

	var sw discoverer.Profile
	switch task.Profile {
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
		logger.Err("Failed to map device profile: '%v'", task.Profile)
		sendError(chanReplies, RequestID, "Failed to map device profile")
		return
	}

	sw.SetLogger(logger.Log)
	sw.SetDebugLogger(logger.Debug)

	if err = sw.Init(cli, "", ""); err != nil {
		sendError(chanReplies, RequestID, err.Error())
		logger.Err("Failed to init profile: %s", err.Error())
		return
	}
	defer sw.Disconnect()

	// Profile is ready, now run tasks
	response := dproto.Response{
		Type:task.Type,
		Errors:make(map[string]string,0),
	}


	// if this is single task, we return task + error
	// if this is 'all', we should return set of errors in some way...
	if task.Type == dproto.PacketType_PLATFORM || task.Type == dproto.PacketType_ALL {
		platform, err := sw.GetPlatform()
		if err != nil { response.Errors[dproto.PacketType_PLATFORM.String()] = err.Error() }
		response.Platform = &platform
	}
	if task.Type == dproto.PacketType_CONFIG || task.Type == dproto.PacketType_ALL {
		config, err := sw.GetConfig()
		if err != nil { response.Errors[dproto.PacketType_CONFIG.String()] = err.Error() }
		response.Config = config
	}
	if task.Type == dproto.PacketType_INTERFACES || task.Type == dproto.PacketType_ALL {
		interfaces, err := sw.GetInterfaces()
		if err != nil { response.Errors[dproto.PacketType_INTERFACES.String()] = err.Error() }
		response.Interfaces = interfaces
	}
	if task.Type == dproto.PacketType_IPS || task.Type == dproto.PacketType_ALL {
		ipifs, err := sw.GetIps()
		if err != nil { response.Errors[dproto.PacketType_IPS.String()] = err.Error() }
		response.Ipifs = ipifs
	}
	if task.Type == dproto.PacketType_LLDP || task.Type == dproto.PacketType_ALL {
		lldp, err := sw.GetLldp()
		if err != nil { response.Errors[dproto.PacketType_LLDP.String()] = err.Error()}
		response.LldpNeighbors = lldp
	}
	if task.Type == dproto.PacketType_UPLINK || task.Type == dproto.PacketType_ALL {
		up, err := sw.GetUplink()
		if err != nil { response.Errors[dproto.PacketType_UPLINK.String()] = err.Error() }
		response.Uplink = up
	}
	if task.Type == dproto.PacketType_VLANS || task.Type == dproto.PacketType_ALL {
		vlans, err := sw.GetVlans()
		if err != nil { response.Errors[dproto.PacketType_VLANS.String()] = err.Error()}
		response.Vlans = vlans
	}

	sendReply(response, chanReplies, RequestID)
}

func sendReply(response dproto.Response, topic string, reply string) {
	logger.Debug("- sending reply... -")
	bs, err := proto.Marshal(&response)
	if err != nil {
		logger.Err("Cannot marshal response for request %s: %s", reply, err.Error())
		return
	}

	msg := nats.Msg{
		Data:bs,
		Subject:topic,
		Reply:reply,
	}

	err = conn.PublishMsg(&msg)
	if err != nil {
		logger.Err("Cannot send response for request %s: %s", reply, err.Error())
	}
}

func sendError(topic string, reply string, message string) {
	msg := dproto.Response{
		Type:dproto.PacketType_ERROR,
		Error:message,
	}
	bs, err := proto.Marshal(&msg)
	if err != nil {
		logger.Err("Cannot marshall nats error message (req.): %s", reply, err.Error())
		return
	}

	natsMsg := nats.Msg{
		Data:bs,
		Subject:topic,
		Reply:reply,
	}
	err = conn.PublishMsg(&natsMsg)
	if err != nil {
		logger.Err("Cannot publish nats error message (req.%s): %s", reply, err.Error())
	}
}