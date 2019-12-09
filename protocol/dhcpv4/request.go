package dhcpv4

import (
	"fmt"
	"git.darknebu.la/maride/pancap/common"
	"github.com/google/gopacket/layers"
)

// Processes the DHCP request packet handed over
func (p *Protocol) processRequestPacket(dhcppacket layers.DHCPv4) {
	p.requestMAC = common.AppendIfUnique(dhcppacket.ClientHWAddr.String(), p.requestMAC)
}

// Generates the summary of all DHCP request packets
func (p *Protocol) generateRequestSummary() string {
	return fmt.Sprintf("%d unique DHCP requests\n%s", len(p.requestMAC), common.GenerateTree(p.requestMAC))
}
