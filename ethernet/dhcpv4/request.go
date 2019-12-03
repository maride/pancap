package dhcpv4

import (
	"fmt"
	"git.darknebu.la/maride/pancap/common"
	"github.com/google/gopacket/layers"
)

var (
	requestMAC []string
)

// Processes the DHCP request packet handed over
func processRequestPacket(dhcppacket layers.DHCPv4) {
	requestMAC = common.AppendIfUnique(dhcppacket.ClientHWAddr.String(), requestMAC)
}

// Generates the summary of all DHCP request packets
func generateRequestSummary() string {
	return fmt.Sprintf("%d unique DHCP requests\n%s", len(requestMAC), common.GenerateTree(requestMAC))
}
