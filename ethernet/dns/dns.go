package dns

import (
	"git.darknebu.la/maride/pancap/output"
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
	output.PrintBlock("DNS Request Summary", generateDNSQuestionSummary())
	output.PrintBlock("DNS Response Summary", generateDNSAnswerSummary())
}
