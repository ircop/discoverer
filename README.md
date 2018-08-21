# discoverer

Little framework for network equipment universal iteraction, in pre-pre-pre-pre-alpha stage :)

The point is networking automation. You can call same functions for any equipment and become structured result (or error :) )

Currently supported equipment are D-Link and Cisco.

Supported methods: see [Profile](https://github.com/ircop/discoverer/blob/master/base/profile.go#L9) interface:

```
	GetPlatform() (Platform, error)
	GetLldp() ([]LldpNeighborship, error)
	GetInterfaces() (map[string]Interface, error)
	GetVlans() ([]Vlan, error)
	GetIps() ([]IPInterface, error)
	GetConfig() (string, error)
	GetUplink() (string, error)
```

Underlying CLI communication software is [ircop/remote-cli](https://github.com/ircop/remote-cli)

SNMP support will be added later


# How to use

First of all you should initialize CLI instance and connect it to device. CLI is outside of profiles because we often need to defer Close() it.

**Initializing cli**

```
import "github.com/ircop/remote-cli"

....

// params:
// - CLI type (CliTypeTelnet | CliTypeSsh
// - IP addr
// - Port
// - Login
// - Password
// - CLI prompt (this is default prompt, you can pass empty string (""))
// - i/o Timeout (read/write) in seconds
cli := remote_cli.New(remote_cli.CliTypeTelnet, "10.110.120.130", 23, "script", "password", `(?msi:[\$%#>]$)`, 3)

// Connect cli
if err := cli.Connect(); err != nil {
  panic(err)
}
defer cli.Close()

// create instance of DLinkDxS profile (it's compatible with D(G|X|E)S series)
sw := DLinkDxS.Profile{}
// or sw := CiscoIOS.Profile{}

// You may set loggers - debug and regular - and make anything you want with discoverer logs
sw.SetLogger(func(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	fmt.Printf("\n")
})
sw.SetDebugLogger(func(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	fmt.Printf("\n")
})

// And finally you should init your profile
// params:
// - cli: it's cli instance, created earlier
// - enable password: you can pass empty string. If it's not empty, profile
//   will try to 'enable' (cisco) or 'enable admin' (dlink) during initialization.
// - snmp community: not used yet
if err = sw.Init(cli, "admin", ""); err != nil {
  panic(err)
}

// And now you can call required methods and become structured output.

// GetPlatform() gathers main device info: model, firmware, HW revision, chassis mac addresses, serial number
platform, err := sw.GetPlatform()
if err != nil {
  panic(err)
}
fmt.Printf("%+v\n", platform)

// Sample output for DLink: {Model:DXS-3350SR Revision:5A1.2A1 Version:4.40-B04 Macs:[00:17:9a:86:6c:00] Serial:}
// or {Model:DES-3028 Revision:A1 Version:2.70.B06 Macs:[00:22:b0:51:1a:e7] Serial:}
// or {Model:DES-3526 Revision:A4 Version:6.20.B18 Macs:[00:21:91:57:ee:c1] Serial:}

// GetVlans() returns vlans existing on platform, and trunk/access ports for this vlans
vlans, err := sw.GetVlans()
if err != nil {
  panic(err)
}
for _, v := range vlans {
  fmt.Printf("%+v\n", v)
}

// Sample output for DLink:
// vlan: {Name:default ID:1 TrunkPorts:[] AccessPorts:[]}
// vlan: {Name:bg ID:210 TrunkPorts:[25 26] AccessPorts:[5]}
// vlan: {Name:manag52 ID:1052 TrunkPorts:[25 26] AccessPorts:[]}
//
// Or for asr1k:
// {Name:Te0/2/0.425 ID:425 TrunkPorts:[Te0/2/0.425] AccessPorts:[]}
// {Name:Te0/2/0.430 ID:430 TrunkPorts:[Te0/2/0.430] AccessPorts:[]}
// {Name:Te0/3/0.25 ID:25 TrunkPorts:[Te0/3/0.25] AccessPorts:[]}
// ...
// Note: for cisco, you should run GetPlatform() first, because vlan gathering process differs for various platforms.
// If you didn't, method will return error with explaination

// And etc., etc. for other methods.

```

# Example

You can see example of very simple web-based tool in [examples](https://github.com/ircop/discoverer/tree/master/examples/web) section.

