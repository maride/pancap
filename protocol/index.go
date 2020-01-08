package protocol

import (
	"git.darknebu.la/maride/pancap/protocol/arp"
	"git.darknebu.la/maride/pancap/protocol/dhcpv4"
	"git.darknebu.la/maride/pancap/protocol/dns"
	"git.darknebu.la/maride/pancap/protocol/http"
)

var (
	Protocols = []Protocol{
		&arp.Protocol{},
		&dhcpv4.Protocol{},
		&dns.Protocol{},
		&http.Protocol{},
	}
)
