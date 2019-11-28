package dhcpv4

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

var (
	requestMAC []string
	responses []dhcpResponse
)

// Called on every DHCP (v4) packet
func HandleDHCPv4Packet(packet gopacket.Packet) error {
	var dhcppacket layers.DHCPv4
	var ethernetpacket layers.Ethernet

	// For some reason I can't find an explanation for,
	// 	packet.Layer(layers.LayerTypeDHCPv4).LayerContents(), which effectively is
	// 	packet.Layers()[3].layerContents(), is empty, but
	//  packet.Layers()[2].layerPayload() contains the correct DHCP packet.
	// ... although both calls should return the same bytes.
	// TODO: Open an Issue on github.com/google/gopacket

	// Decode raw packet into DHCPv4
	decodeDHCPErr := dhcppacket.DecodeFromBytes(packet.Layers()[2].LayerPayload(), gopacket.NilDecodeFeedback)
	if decodeDHCPErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeDHCPErr
	}

	// And decode raw packet into Ethernet
	decodeEthernetErr := ethernetpacket.DecodeFromBytes(packet.Layer(layers.LayerTypeEthernet).LayerContents(), gopacket.NilDecodeFeedback)
	if decodeEthernetErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeEthernetErr
	}

	// Examine packet further
	if dhcppacket.Operation == layers.DHCPOpRequest {
		// Request packet
		appendIfUnique(dhcppacket.ClientHWAddr.String(), requestMAC)
	} else {
		// Response/Offer packet
		addResponseEntry(dhcppacket.ClientIP.String(), dhcppacket.YourClientIP.String(), dhcppacket.ClientHWAddr.String(), ethernetpacket.SrcMAC.String())
	}

	return nil
}

// Print summary after all packets are processed
func PrintDHCPv4Summary() {
	headline := color.New(color.FgRed, color.Bold)
	headline.Println("DHCP Requests")
	printRequestSummary()
	headline.Println("DHCP Responses/Offers")
	printResponseSummary()
}

// Prints the summary of all DHCP request packets
func printRequestSummary() {
	fmt.Printf("%d unique DHCP requests\n", len(requestMAC))
	printTree(requestMAC)
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
	printTree(tmpaddr)
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