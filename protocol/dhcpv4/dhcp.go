package dhcpv4

import (
	"git.darknebu.la/maride/pancap/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Protocol struct {
	hostnames []hostname
	networkSetup map[layers.DHCPOpt][]byte
	requestMAC []string
	responses []dhcpResponse
}

// Checks if the given packet is a DHCP packet we can process
func (p *Protocol) CanAnalyze(packet gopacket.Packet) bool {
	return packet.Layer(layers.LayerTypeDHCPv4) != nil && packet.Layer(layers.LayerTypeEthernet) != nil && packet.Layers()[2].LayerPayload() != nil
}

// Analyzes the given DHCP packet
func (p *Protocol) Analyze(packet gopacket.Packet) error {
	var dhcppacket layers.DHCPv4
	var ethernetpacket layers.Ethernet

	// Check if it's the first run - init networkSetup map then
	if p.networkSetup == nil {
		p.networkSetup = make(map[layers.DHCPOpt][]byte)
	}

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
		p.processRequestPacket(dhcppacket)
	} else {
		// Response/Offer packet
		p.processResponsePacket(dhcppacket, ethernetpacket)
	}

	// Check for Hostname DHCP option (12)
	p.checkForHostname(dhcppacket)
	p.checkForNetworkInfos(dhcppacket)

	return nil
}

// Print summary after all packets are processed
func (p *Protocol) PrintSummary() {
	output.PrintBlock("DHCP Network Overview", p.generateNetworkSummary())
	output.PrintBlock("DHCP Requests", p.generateRequestSummary())
	output.PrintBlock("DHCP Responses/Offers", p.generateResponseSummary())
	output.PrintBlock("DHCP Hostnames", p.generateHostnamesSummary())
}
