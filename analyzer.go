package main

import (
	"errors"
	"fmt"
	"git.darknebu.la/maride/pancap/ethernet"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Analyzes the given packet source
func analyzePCAP(source *gopacket.PacketSource, linkType layers.LinkType) error {
	// Switch over link type to determine correct module to ask for analysis
	switch linkType {
	case layers.LinkTypeEthernet:
		// Ethernet
		return ethernet.Analyze(source)
	}

	// if we reach this point, the given PCAP contains a link type we can't handle (yet).
	errorMsg := fmt.Sprintf("Asked for link type %s (ID %d), but not supported by pancap. :( sorry!", linkType.String(), linkType)
	return errors.New(errorMsg)
}
