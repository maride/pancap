package analyze

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/maride/pancap/output"
	"github.com/maride/pancap/protocol"
)

var (
	// Store total amount and amount of visited packets
	totalPackets     int
	processedPackets int
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

		// Track if we didn't process a packet
		processed := false

		// Iterate over all possible protocols
		for _, p := range protocol.Protocols {
			// Check if this protocol can handle this packet
			if p.CanAnalyze(packet) {
				handleErr(p.Analyze(packet))
				processed = true
			}
		}

		// Register communication for graph
		output.AddPkgToGraph(packet)

		// Raise statistics
		totalPackets += 1
		if processed {
			processedPackets += 1
		}
	}

	return nil
}

// Prints all the summaries.
func PrintSummary() {
	// First, print base information collected while analyzing
	content := fmt.Sprintf("Processed %d out of %d packets (%d%%)", processedPackets, totalPackets, processedPackets*100/totalPackets)
	output.PrintBlock("Overall statistics", content)

	// Print summary of each protocol
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
