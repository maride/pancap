package protocol

import (
	"github.com/maride/pancap/protocol/arp"
	"github.com/maride/pancap/protocol/dhcpv4"
	"github.com/maride/pancap/protocol/dns"
	"github.com/maride/pancap/protocol/http"
)

var (
	Protocols = []Protocol{
		&arp.Protocol{},
		&dhcpv4.Protocol{},
		&dns.Protocol{},
		&http.Protocol{},
	}
)
