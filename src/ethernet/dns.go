package ethernet

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

var (
	numQuestions int
	numAnswers int
)

// Called on every DNS packet
func processDNSPacket(packet gopacket.Packet) error {
	var dnspacket layers.DNS

	// Decode raw packet into DNS
	decodeErr := dnspacket.DecodeFromBytes(packet.ApplicationLayer().LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Further process the packet
	processDNSQuestion(dnspacket.Questions)
	processDNSAnswers(dnspacket.Answers)

	// No error encountered, return clean
	return nil
}

// Called on every DNS packet to process questions
func processDNSQuestion(questions []layers.DNSQuestion) {
	for _, _ = range questions {
		numQuestions++
	}
}

// Called on every DNS packet to process response(s)
func processDNSAnswers(answers []layers.DNSResourceRecord) {
	for _, _ = range answers {
		numAnswers++
	}
}

// Print a summary after all packets were processed
func printDNSSummary() {
	log.Printf("%d DNS Questions, %d DNS Answers in total", numQuestions, numAnswers)
}