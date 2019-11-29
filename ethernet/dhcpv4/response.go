package dhcpv4

import (
	"fmt"
	"git.darknebu.la/maride/pancap/common"
	"github.com/google/gopacket/layers"
	"log"
)

var (
	responses []dhcpResponse
)

func processResponsePacket(dhcppacket layers.DHCPv4, ethernetpacket layers.Ethernet) {
	addResponseEntry(dhcppacket.ClientIP.String(), dhcppacket.YourClientIP.String(), dhcppacket.ClientHWAddr.String(), ethernetpacket.SrcMAC.String())
}

// Prints the summary of all DHCP offer packets
func printResponseSummary() {
	var tmpaddr []string

	// Iterate over all responses
	for _, r := range responses {
		addition := ""

		if r.askedFor {
			addition = " which the client explicitly asked for."
		}

		tmpaddr = append(tmpaddr, fmt.Sprintf("%s offered %s IP address %s%s", r.serverMACAddr, r.destMACAddr, r.newIPAddr, addition))
	}

	// Draw as tree
	common.PrintTree(tmpaddr)
}

// Adds a new response entry. If an IP address was already issued or a MAC asks multiple times for DNS, the case is examined further
func addResponseEntry(newIP string, yourIP string, destMAC string, serverMAC string) {
	// Check if client asked for a specific address (which was granted by the DHCP server)
	askedFor := false
	if newIP == "0.0.0.0" {
		// Yes, client asked for a specific address. Most likely not the first time in this network.
		newIP = yourIP
		askedFor = true
	}

	for _, r := range responses {
		// Check for interesting cases
		if r.destMACAddr == destMAC {
			// The same client device received multiple IP addresses, let's examine further
			if r.newIPAddr == newIP {
				// the handed IP is the same - this is ok, just badly configured
				if r.serverMACAddr == serverMAC {
					// Same DHCP server answered.
					log.Printf("MAC address %s received the same IP address multiple times via DHCP by the same server.", destMAC)
				} else {
					// Different DHCP servers answered, but with the same address - strange network, but ok...
					log.Printf("MAC address %s received the same IP address multiple times via DHCP by different servers.", destMAC)
				}
			} else {
				// far more interesting - one client received multiple addresses
				if r.serverMACAddr == serverMAC {
					// Same DHCP server answered.
					log.Printf("MAC address %s received different IP addresses (%s, %s) multiple times via DHCP by the same server.", destMAC, r.newIPAddr, newIP)
				} else {
					// Different DHCP servers answered, with different addresses - possibly an attempt to build up MitM
					log.Printf("MAC address %s received different IP addresses (%s, %s) multiple times via DHCP by different servers (%s, %s).", destMAC, r.newIPAddr, newIP, r.serverMACAddr, serverMAC)
				}
			}
		}
	}

	// Add a response entry - even if we found some "strange" behavior before.
	responses = append(responses, dhcpResponse{
		destMACAddr:   destMAC,
		newIPAddr:     newIP,
		serverMACAddr: serverMAC,
		askedFor:      askedFor,
	})
}
