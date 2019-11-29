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

// Prints the summary of all DHCP request packets
func printRequestSummary() {
	fmt.Printf("%d unique DHCP requests\n", len(requestMAC))
	common.PrintTree(requestMAC)
}
