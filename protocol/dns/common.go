package dns

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
)

var (
	privateBlocks = []net.IPNet{
		{net.IPv4(10, 0, 0, 0), net.IPv4Mask(255, 0, 0, 0)},      // 10.0.0.0/8
		{net.IPv4(172, 16, 0, 0), net.IPv4Mask(255, 240, 0, 0)},  // 172.16.0.0/12
		{net.IPv4(192, 168, 0, 0), net.IPv4Mask(255, 255, 0, 0)}, // 192.168.0.0/24
		{net.IPv4(100, 64, 0, 0), net.IPv4Mask(255, 192, 0, 0)},  // 100.64.0.0/10
		{net.IPv4(169, 254, 0, 0), net.IPv4Mask(255, 255, 0, 0)}, // 169.254.0.0/16
	}
)

// Processes the given dnstype and raises its stats in the given array
func (p *Protocol) processType(typearr map[layers.DNSType]int, dnstype layers.DNSType) {
	typearr[dnstype]++
}

// Checks if the given IP is in a private range or not
func ipIsPrivate(ip net.IP) bool {
	// check every private IP block for our IP
	for _, block := range privateBlocks {
		if block.Contains(ip) {
			// found, is a private IP
			return true
		}
	}

	// Not in any of the private blocks, not private
	return false
}

// Generates a summary string for DNS types in the given array
func (p *Protocol) generateDNSTypeSummary(typearr map[layers.DNSType]int) string {
	var answerarr []string

	// Iterate over all possible DNS types
	for iter, typeelem := range typearr {
		// Read amount of type hits for this type
		answerarr = append(answerarr, fmt.Sprintf("%d %s", typeelem, iter.String()))
	}

	// Check if we even processed a single type
	if len(answerarr) == 0 {
		// we didn't, strange.
		return "(no types encountered)"
	}

	// now, glue all array elements together
	answerstr := ""
	for iter, elem := range answerarr {
		// Check if we need to apply to proper sentence rules
		if iter == 0 {
			// We don't need to append yet
			answerstr = elem
		} else if iter == len(answerarr)-1 {
			// Last element, use "and" instead of a comma
			answerstr = fmt.Sprintf("%s and %s", answerstr, elem)
		} else {
			// Some entry, just add it with a comma
			answerstr = fmt.Sprintf("%s, %s", answerstr, elem)
		}
	}

	return answerstr
}
