package protocol

import "github.com/google/gopacket"

type Protocol interface {
	CanAnalyze(gopacket.Packet) bool
	Analyze(gopacket.Packet) error
	PrintSummary()
}
