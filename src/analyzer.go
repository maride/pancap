package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

// Analyzes the given packet source
func analyzePCAP(source *gopacket.PacketSource, linkType layers.LinkType) error {
	log.Printf("PCAP capture link type is %s (ID %d)", getNameOfLinkType(linkType), linkType)
	// TODO: maybe, just maybe, we wanna print more here than just the link type :)
	_, _ = source, linkType
	return nil
}

// Returns the name of the LinkType constant handed over
func getNameOfLinkType(lt layers.LinkType) string {
	return lt.String()
}