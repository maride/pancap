package ethernet

import (
	"git.darknebu.la/maride/pancap/ethernet/arp"
	"git.darknebu.la/maride/pancap/ethernet/dhcpv4"
	"git.darknebu.la/maride/pancap/ethernet/dns"
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
			handleErr(dns.ProcessDNSPacket(packet))
		}

		if packet.Layer(layers.LayerTypeARP) != nil {
			// Handle ARP packet
			handleErr(arp.ProcessARPPacket(packet))
		}

		if packet.Layer(layers.LayerTypeDHCPv4) != nil {
			// Handle DHCP (v4) packet
			handleErr(dhcpv4.HandleDHCPv4Packet(packet))
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

// Handles an error, if err is not nil.
func handleErr(err error) {
	// (hopefully) most calls to this function will contain a nil error, so we need to check if we really got an error
	if err != nil {
		log.Printf("Encountered error while examining packets, continuing anyway. Error: %s", err.Error())
	}
}
