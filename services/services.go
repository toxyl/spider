package services

type Service struct {
	Ports  []int  `yaml:"ports"`
	Banner string `yaml:"banner"`
}

type Services map[string]Service

var services = &Services{
	"Apache Kafka": Service{
		Ports:  []int{9092},
		Banner: "This server is Kafka",
	},
	"BitTorrent": Service{
		Ports:  []int{6881},
		Banner: "BitTorrent protocol",
	},
	"Cassandra": Service{
		Ports:  []int{9042},
		Banner: "CQL_VERSION",
	},
	"DirectAdmin": Service{
		Ports:  []int{2222},
		Banner: "220 DirectAdmin Login",
	},
	"DNS": Service{
		Ports:  []int{53},
		Banner: "BIND",
	},
	"Docker": Service{
		Ports:  []int{2375, 2376},
		Banner: "HTTP/1.1 404 Not Found",
	},
	"Elasticsearch": Service{
		Ports:  []int{9200, 9300},
		Banner: "{",
	},
	"FTP": Service{
		Ports:  []int{21},
		Banner: "220 [tapeworm,bookworm,[str:12],[mix:12]] FTP Server Ready",
	},
	"FTPS": Service{
		Ports:  []int{990},
		Banner: "220 FTPS server ready",
	},
	"HTTP": Service{
		Ports:  []int{80, 8000, 8080},
		Banner: "HTTP/1.1 404 Not Found",
	},
	"HTTPS": Service{
		Ports:  []int{443, 8443},
		Banner: "HTTP/1.1 404 Not Found",
	},
	"IKEv2": Service{
		Ports:  []int{4500},
		Banner: "IKEv2",
	},
	"IMAP": Service{
		Ports:  []int{143},
		Banner: "* OK [CAPABILITY IMAP4rev1 LITERAL+ SASL-IR LOGIN-REFERRALS ID ENABLE STARTTLS AUTH=PLAIN AUTH=LOGIN] Dovecot ready.",
	},
	"LDAP": Service{
		Ports:  []int{389},
		Banner: "220-OpenLDAP",
	},
	"LDAP-over-SSL": Service{
		Ports:  []int{636},
		Banner: "0x55 0x53 0x45 0x52 0x20 0x43 0x6f 0x[int:1]d",
	},
	"Matrix": Service{
		Ports:  []int{8448},
		Banner: "MATRIX",
	},
	"Memcached": Service{
		Ports:  []int{11211},
		Banner: "ERROR\r\n",
	},
	"Microsoft Exchange Server": Service{
		Ports:  []int{135},
		Banner: "ncacn_http/1.0",
	},
	"Microsoft SQL Server": Service{
		Ports:  []int{1433},
		Banner: "Microsoft SQL Server",
	},
	"Microsoft SQL Server Browser": Service{
		Ports:  []int{1434},
		Banner: "SQL Server Browser",
	},
	"MongoDB": Service{
		Ports:  []int{27017},
		Banner: "MongoDB",
	},
	"MQTT": Service{
		Ports:  []int{1883, 8883},
		Banner: "MQTT",
	},
	"MySQL": Service{
		Ports:  []int{3306},
		Banner: "[int:1].[int:1].[int:2]-0ubuntu0.18.04.1",
	},
	"NETBIOS": Service{
		Ports:  []int{137, 138, 139},
		Banner: "NBTSTAT",
	},
	"NTP": Service{
		Ports:  []int{123},
		Banner: "NTP",
	},
	"OpenVPN": Service{
		Ports:  []int{1194},
		Banner: "[OpenVPN,ClosedVPN]",
	},
	"POP3": Service{
		Ports:  []int{110},
		Banner: "+OK Microsoft Exchange Server 200[int:1] POP3 server ready",
	},
	"Postfix": Service{
		Ports:  []int{587},
		Banner: "220 mail.[tapeworm,bookworm,[str:12],[mix:12]].[f.u,not.here,go.away] ESMTP",
	},
	"PostgresSQL": Service{
		Ports:  []int{5432},
		Banner: "PostgreSQL",
	},
	"PPTP": Service{
		Ports:  []int{1723},
		Banner: "1",
	},
	"RDP": Service{
		Ports:  []int{3389},
		Banner: "\\x03\\x00\\x00\x0b\\x06\xd0\\x00\\x00\\x[0-255]\\x00",
	},
	"Redis": Service{
		Ports:  []int{6379},
		Banner: "Redis",
	},
	"Redis Sentinel": Service{
		Ports:  []int{26379},
		Banner: "+sentinel",
	},
	"RPC": Service{
		Ports:  []int{111, 1025},
		Banner: "-ERR unsupported command ''",
	},
	"Rsync": Service{
		Ports:  []int{873},
		Banner: "@RSYNCD: [1-32].0",
	},
	"SaltStack": Service{
		Ports:  []int{4433, 4505, 4506},
		Banner: "salt",
	},
	"SIP": Service{
		Ports:  []int{5060},
		Banner: "SIP/[int:1].0 200 OK",
	},
	"SMB": Service{
		Ports:  []int{445},
		Banner: "Windows for Workgroups",
	},
	"SMTP": Service{
		Ports:  []int{25},
		Banner: "220 [tapeworm,bookworm,[str:12],[mix:12]] ESMTP Postfix",
	},
	"SMTPS": Service{
		Ports:  []int{465},
		Banner: "220 mail.[tapeworm,bookworm,[str:12],[mix:12]].[f.u,not.here,go.away] ESMTP",
	},
	"SNMP": Service{
		Ports:  []int{161},
		Banner: "SNMP",
	},
	"SNMP Trap": Service{
		Ports:  []int{162},
		Banner: "TRAP",
	},
	"Socks": Service{
		Ports:  []int{1080},
		Banner: "SOCKS",
	},
	"Splunk": Service{
		Ports:  []int{8089},
		Banner: "Splunk",
	},
	"SSH": Service{
		Ports:  []int{22},
		Banner: "SSH-2.0-OpenSSH_[int:1].[int:2]",
	},
	"Telnet": Service{
		Ports:  []int{23},
		Banner: "220 [tapeworm,bookworm,[str:12],[mix:12]] Telnet",
	},
	"TFTP": Service{
		Ports:  []int{69},
		Banner: "TFTP",
	},
	"VNC": Service{
		Ports:  []int{5900, 5901},
		Banner: "RFB [int:3].[int:3]",
	},
	"VRRP": Service{
		Ports:  []int{112},
		Banner: "VRRP",
	},
	"ZooKeeper": Service{
		Ports:  []int{2181},
		Banner: "Environment",
	},
}
