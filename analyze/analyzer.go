package analyze

import (
	"git.darknebu.la/maride/pancap/protocol"
	"github.com/google/gopacket"
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

		// Iterate over all possible protocols
		for _, p := range protocol.Protocols {
			// Check if this protocol can handle this packet
			if p.CanAnalyze(packet) {
				handleErr(p.Analyze(packet))
			}
		}
	}

	return nil
}

// Prints all the summaries.
func PrintSummary() {
	for _, p := range protocol.Protocols {
		p.PrintSummary()
	}
}

// Handles an error, if err is not nil.
func handleErr(err error) {
	// (hopefully) most calls to this function will contain a nil error, so we need to check if we really got an error
	if err != nil {
		log.Printf("Encountered error while examining packets, continuing anyway. Error: %s", err.Error())
	}
}

