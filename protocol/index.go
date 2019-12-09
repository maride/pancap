package protocol

import (
	"git.darknebu.la/maride/pancap/protocol/arp"
	"git.darknebu.la/maride/pancap/protocol/dhcpv4"
	"git.darknebu.la/maride/pancap/protocol/dns"
)

var (
	Protocols = []Protocol{
		&arp.Protocol{},
		&dhcpv4.Protocol{},
		&dns.Protocol{},
	}
)
