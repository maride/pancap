package main

import (
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
)

var (
	filenameFlag *string
)

// Registers the flag --file
func registerFileFlags() {
	filenameFlag = flag.String("file", "", "PCAP file to base analysis on")
}

// Opens the PCAP, returns its packets and the link type or an error
func openPCAP() (*gopacket.PacketSource, layers.LinkType, error) {
	// Open specified file
	handle, openErr := pcap.OpenOffline(*filenameFlag)
	if openErr != nil {
		// There were some problems opening the file
		return nil, 0, openErr
	}

	// Output basic information about this PCAP
	log.Printf("PCAP capture link type is %s (ID %d)", handle.LinkType().String(), handle.LinkType())

	// Open given handle as packet source and return it
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	return packetSource, handle.LinkType(), nil
}