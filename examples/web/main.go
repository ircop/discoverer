package main

import (
	"net/http"
	"fmt"
	"net"
	"strconv"
	"github.com/ircop/discoverer/profiles/base"
	"github.com/ircop/discoverer/profiles/DLinkDxS"
	"github.com/ircop/discoverer/profiles/CiscoIOS"
	"encoding/json"
	"github.com/ircop/remote-cli"
)

func main() {
	http.HandleFunc("/execute", handle)
	fmt.Printf("Listening on :8765\n")
	http.ListenAndServe(":8765", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	profile := r.URL.Query().Get("profile")
	login := r.URL.Query().Get("login")
	password := r.URL.Query().Get("password")
	enable := r.URL.Query().Get("enable")
	protoStr := r.URL.Query().Get("dproto")
	ip := r.URL.Query().Get("ip")
	portStr := r.URL.Query().Get("port")
	method := r.URL.Query().Get("method")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// we need at lease login/pw/dproto/ip/profile
	if login == "" || password == "" || protoStr == "" || ip == "" || portStr == "" || profile == "" || method == "" {
		returnError(w, fmt.Sprintf("Some of login/password/dproto/ip/profile is/are empty (%s/%s/%s/%s/%s/%s/%s)", login, password, protoStr, ip, portStr, profile, method) )
		return
	}

	// check some vars
	if net.ParseIP(ip) == nil {
		returnError(w, fmt.Sprintf("'%s' is not valid IP-address", ip))
		return
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		returnError(w, fmt.Sprintf("Cannot parse port '%s': %s", portStr, err.Error()))
		return
	}
	if port < 1 || port > 65535 {
		returnError(w, fmt.Sprintf("Port '%d' is not valid port number", portStr, err.Error()))
		return
	}

	proto := 0
	switch protoStr{
	case "telnet":
		proto = discoverer.CliTypeTelnet
		break
	case "ssh":
		proto = discoverer.CliTypeSsh
		break
	default:
		returnError(w, fmt.Sprintf("'dproto' should be 'telnet' or 'ssh', not '%s'", protoStr))
		return
	}

	var prof discoverer.Profile
	switch profile {
	case "dlink":
		prof = &DLinkDxS.Profile{}
		break
	case "cisco":
		prof = &CiscoIOS.Profile{}
		break
	default:
		returnError(w, fmt.Sprintf("'profile' should be 'cisco' or 'dlink', not '%s'", profile))
		return
	}

	// init cli
	cli := remote_cli.New(proto, ip, int(port), login, password, "", 10)

	execute(prof, cli, enable, method, w)
}

func execute(prof discoverer.Profile, cli *remote_cli.Cli, enable string, method string, w http.ResponseWriter) {
	// connect cli
	err := cli.Connect()
	if err != nil {
		returnError(w, fmt.Sprintf("Cannot connect via cli: %s", err.Error()))
		return
	}
	defer cli.Close()

	// init profile
	err = prof.Init(cli, enable, "")
	if err != nil {
		returnError(w, fmt.Sprintf("Failed to init profile: %s", err.Error()))
		return
	}

	// log output
	prof.SetLogger(func(msg string, args ...interface{}) {
		fmt.Printf("%s\n", fmt.Sprintf(msg, args...))
	})
	prof.SetDebugLogger(func(msg string, args ...interface{}) {
		fmt.Printf("%s\n", fmt.Sprintf(msg, args...))
	})

	switch method {
	case "get_platform":
		platform, err := prof.GetPlatform()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		fmt.Printf("PLATFORM:\n%+v\n", platform)
		writeJson(w, platform)
		return
	case "get_lldp":
		lldp, err := prof.GetLldp()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		m := make(map[string]interface{})
		m["lldp"] = lldp
		writeJson(w, m)
	case "get_interfaces":
		ifs, err := prof.GetInterfaces()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		if err != nil {
			returnError(w, err.Error())
			return
		}
		m := make(map[string]interface{})
		m["ifs"] = ifs
		writeJson(w, m)
	case "get_vlans":
		vlans, err := prof.GetVlans()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		m := make(map[string]interface{})
		m["vlans"] = vlans
		writeJson(w, m)
	case "get_ips":
		ips, err := prof.GetIps()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		m := make(map[string]interface{})
		m["ips"] = ips
		writeJson(w, m)
	case "get_config":
		cfg, err := prof.GetConfig()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		m := make(map[string]string)
		m["config"] = cfg
		writeJson(w, m)
	case "get_uplink":
		up, err := prof.GetUplink()
		if err != nil {
			returnError(w, err.Error())
			return
		}
		// jsonify
		m := make(map[string]string)
		m["uplink"] = up
		writeJson(w, m)
	default:
		returnError(w, "Unknown method")
	}
}

func writeJson(w http.ResponseWriter, value interface{}) {
	bytes, err := json.Marshal(value)
	if err != nil {
		returnError(w, fmt.Sprintf("Error marshaling into json: %s", err.Error()))
		return
	}

	fmt.Fprintf(w, "%s", bytes)
}


func returnError(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, `{"error":"%s"}`, message)
}
