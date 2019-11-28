package ethernet

import (
	"git.darknebu.la/maride/pancap/ethernet/arp"
	"git.darknebu.la/maride/pancap/ethernet/dns"
	"git.darknebu.la/maride/pancap/ethernet/dhcpv4"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

func Analyze(source *gopacket.PacketSource) error {
	// Loop over all packets now
	for {
		packet, packetErr := source.NextPacket()
		if packet == nil {
			// We iterated over all packets, we're done here
			break
		} else if packetErr != nil {
			// encountered some problem, report it
			log.Printf("Encountered a problem with a packet: %s", packetErr.Error())
			continue
		}

		if packet.Layer(layers.LayerTypeDNS) != nil {
			// Handle DNS packet
			dns.ProcessDNSPacket(packet)
		}

		if packet.Layer(layers.LayerTypeARP) != nil {
			// Handle ARP packet
			arp.ProcessARPPacket(packet)
		}

		if packet.Layer(layers.LayerTypeDHCPv4) != nil {
			// Handle DHCP (v4) packet
			dhcpv4.HandleDHCPv4Packet(packet)
		}
	}

	// After processing all packets, print summary
	printSummary()

	return nil
}

// Prints all the summaries.
func printSummary() {
	arp.PrintARPSummary()
	dns.PrintDNSSummary()
	dhcpv4.PrintDHCPv4Summary()
}
