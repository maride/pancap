package dhcpv4

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
		processRequestPacket(dhcppacket)
	} else {
		// Response/Offer packet
		processResponsePacket(dhcppacket, ethernetpacket)
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
