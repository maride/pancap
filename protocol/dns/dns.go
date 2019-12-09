package dns

import (
	"git.darknebu.la/maride/pancap/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Protocol struct {}

func (p *Protocol) CanAnalyze(packet gopacket.Packet) bool {
	return packet.Layer(layers.LayerTypeDNS) != nil
}

// Analyzes the given DHCP packet
func (p *Protocol) Analyze(packet gopacket.Packet) error {
	var dnspacket layers.DNS

	// Decode raw packet into DNS
	decodeErr := dnspacket.DecodeFromBytes(packet.ApplicationLayer().LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Further process the packet
	p.processDNSQuestion(dnspacket.Questions)
	p.processDNSAnswer(dnspacket.Answers)

	// No error encountered, return clean
	return nil
}

// Print a summary after all DNS packets were processed
func (p *Protocol) PrintSummary() {
	output.PrintBlock("DNS Request Summary", p.generateDNSQuestionSummary())
	output.PrintBlock("DNS Response Summary", p.generateDNSAnswerSummary())
}
