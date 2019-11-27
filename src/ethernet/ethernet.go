package ethernet

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"./dns"
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

		// Check if we can do some Application Layer statistics with this packet
		if packet.ApplicationLayer() != nil {
			// We can, switch over the type
			switch packet.ApplicationLayer().LayerType() {
			case layers.LayerTypeDNS:
				// Handle DNS packet
				dns.ProcessDNSPacket(packet)
			}
		}
	}

	// After processing all packets, print summary
	printSummary()

	return nil
}

// Prints all the summaries.
func printSummary() {
	dns.PrintDNSSummary()
}
