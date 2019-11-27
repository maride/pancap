package dns

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Called on every DNS packet
func ProcessDNSPacket(packet gopacket.Packet) error {
	var dnspacket layers.DNS

	// Decode raw packet into DNS
	decodeErr := dnspacket.DecodeFromBytes(packet.ApplicationLayer().LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Further process the packet
	processDNSQuestion(dnspacket.Questions)
	processDNSAnswer(dnspacket.Answers)

	// No error encountered, return clean
	return nil
}

// Print a summary after all DNS packets were processed
func PrintDNSSummary() {
	headline := color.New(color.FgRed, color.Bold)
	headline.Println("DNS Request Summary")
	printDNSQuestionSummary()
	headline.Println("DNS Response Summary")
	printDNSAnswerSummary()
}
