package http

import (
	"git.darknebu.la/maride/pancap/common"
	"git.darknebu.la/maride/pancap/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
)

type Protocol struct {
	initialized bool
	requestFactory *httpRequestFactory
	responseFactory *httpResponseFactory
	requestPool *tcpassembly.StreamPool
	responsePool *tcpassembly.StreamPool
	requestAssembler *tcpassembly.Assembler
	responseAssembler *tcpassembly.Assembler
}

// Checks if the given packet is an HTTP packet we can process
func (p *Protocol) CanAnalyze(packet gopacket.Packet) bool {
	return packet.Layer(layers.LayerTypeTCP) != nil && packet.Layer(layers.LayerTypeTLS) == nil
}

// Analyzes the given HTTP packet
func (p *Protocol) Analyze(packet gopacket.Packet) error {
	// Check if we need to init
	if !p.initialized {
		// Initialize
		p.requestFactory = &httpRequestFactory{}
		p.responseFactory = &httpResponseFactory{}
		p.requestPool = tcpassembly.NewStreamPool(p.requestFactory)
		p.responsePool = tcpassembly.NewStreamPool(p.responseFactory)
		p.requestAssembler = tcpassembly.NewAssembler(p.requestPool)
		p.responseAssembler = tcpassembly.NewAssembler(p.responsePool)
		p.initialized = true
	}

	// Try to cast packet and assemble HTTP stream
	tcp := packet.TransportLayer().(*layers.TCP)
	p.requestAssembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
	p.responseAssembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

	return nil
}

// Print a summary after all packets are processed
func (p *Protocol) PrintSummary() {
	output.PrintBlock("HTTP Requests", common.GenerateTree(requestSummaryLines))
	output.PrintBlock("HTTP Responses", common.GenerateTree(responseSummaryLines))
}
