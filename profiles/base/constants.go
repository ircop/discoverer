package discoverer

const (
	CliTypeTelnet	= 1
	CliTypeSsh		= 2
	CliTypeNone		= 0
)

// DeviceProfile explains how to communicate with network device: cli type (telnet/ssh),
// login/password, snmp, snmp community
//type DeviceProfile struct {
	//IP				string
	//Port			int
	//Login			string
	//Password		string
	//CliType			int
	//Community		string
	//LoginPrompt		string
	//PasswordPrompt	string
	//Prompt			string
	//Timeout			int
//}

