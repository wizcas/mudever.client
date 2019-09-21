package protocol

var optionNames = map[OptByte]string{
	0:   "BINARY-TRANSMISSION",
	1:   "ECHO",
	3:   "SUPPRESS-GO-AHEAD",
	5:   "STATUS",
	6:   "TIMING-MARK",
	24:  "TERMINAL-TYPE",
	31:  "NEGOTIATE-ABOUT-WND-SIZE",
	32:  "TERMINAL-SPEED",
	33:  "REMOTE-FLOW-CONTROL",
	34:  "LINE-MODE",
	36:  "ENVIRONMENT-VARIABLES",
	37:  "AUTHENTICATION",
	38:  "ENCRYPTION",
	39:  "NEW-ENVIRONMENT",
	40:  "TN3270E",
	41:  "XAUTH",
	42:  "CHARSET",
	43:  "RSP",
	44:  "COM-PORT-CONTROL",
	45:  "TELNET-SUPPRESS-LOCAL-ECHO",
	46:  "TELNET-START-TLS",
	47:  "KERMIT",
	48:  "SEND-URL",
	49:  "FORWARD-X",
	69:  "MSDP",  // MUD Server Data Protocol
	70:  "MSSP",  // MUD Server Status Protocol
	85:  "MCCP1", // Mud Client Compression Protocol version 1
	86:  "MCCP2", // Mud Client Compression Protocol version 2
	90:  "MSP",   // MUD Sound Protocol
	91:  "MXP",   // MUD eXtension Protocol
	93:  "ZMP",   // Zenith MUD Protocol
	138: "TEL-OPT-PRAGMA-LOGON",
	139: "TEL-OPT-SSPI-LOGON",
	140: "TEL-OPT-PRAGMA-HEARTBEAT",
	200: "ATCP", // Achaea Telnet Client Protocol
	201: "GMCP", // Generic MUD Communication Protocol
}

var cmdNames = map[CmdByte]string{
	// CONTROL FUNCTIONS
	240: "SE",
	241: "NOP",
	242: "DM",
	243: "BRK",
	244: "IP",
	245: "AO",
	246: "AYT",
	247: "EC",
	248: "EL",
	249: "GA",
	250: "SB",
	// COMMANDS
	251: "WILL",
	252: "WONT",
	253: "DO",
	254: "DONT",
	255: "IAC",
}
