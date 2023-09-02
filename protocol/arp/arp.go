package arp

import (
	"fmt"
	"github.com/maride/pancap/common"
	"github.com/maride/pancap/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

var (
	arpStatsList []arpStats
	devices []arpDevice
	linkLocalBlock = net.IPNet{
		IP:   net.IPv4(169, 254, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}
)

type Protocol struct {}

// Checks if the given packet is an ARP packet we can process
func (p *Protocol) CanAnalyze(packet gopacket.Packet) bool {
	return packet.Layer(layers.LayerTypeARP) != nil
}

// Analyzes the given ARP packet
func (p *Protocol) Analyze(packet gopacket.Packet) error {
	var arppacket layers.ARP

	// Decode raw packet into ARP
	decodeErr := arppacket.DecodeFromBytes(packet.Layer(layers.LayerTypeARP).LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Convert MAC address byte array to string
	sourceAddr := net.HardwareAddr(arppacket.SourceHwAddress).String()
	participant := p.getStatOrCreate(sourceAddr)

	// Raise stats
	if arppacket.Operation == layers.ARPRequest {
		// Request packet
		participant.asked++
		participant.askedList = common.AppendIfUnique(net.IP(arppacket.DstProtAddress).String(), participant.askedList)

		// Add device entry
		p.addDeviceEntry(sourceAddr, net.IP(arppacket.SourceProtAddress).String())
	} else {
		// Response packet
		participant.answered++
		participant.answeredList = common.AppendIfUnique(net.IP(arppacket.SourceProtAddress).String(), participant.answeredList)

		// Add device entry
		p.addDeviceEntry(sourceAddr, net.IP(arppacket.SourceProtAddress).String())
	}

	return nil
}

// Print a summary after all packets are processed
func (p *Protocol) PrintSummary() {
	output.PrintBlock("ARP traffic summary", p.generateTrafficStats())
	output.PrintBlock("ARP LAN overview", p.generateLANOverview())
}

// Generates an answer regarding the ARP traffic
func (p *Protocol) generateTrafficStats() string {
	var tmparr []string

	// Iterate over all participants
	for _, p := range arpStatsList {
		// produce a meaningful output
		if p.asked > 0 {
			// device asked at least for one IP
			if p.answered > 0 {
				// and also answered requests
				tmparr = append(tmparr, fmt.Sprintf("%s asked for %d addresses and answered %d requests", p.macaddr, p.asked, p.answered))
			} else {
				// only asked, never answered
				tmparr = append(tmparr, fmt.Sprintf("%s asked for %d addresses", p.macaddr, p.asked))
			}
		} else {
			// Answered, but never asked for any addresses
			tmparr = append(tmparr, fmt.Sprintf("%s answered %d requests", p.macaddr, p.answered))
		}
	}

	// And print it as a tree
	return common.GenerateTree(tmparr)
}

// Generates an overview over all connected devices in the LAN
func (p *Protocol) generateLANOverview() string {
	var tmparr []string

	// iterate over all devices
	for _, d := range devices {
		tmparr = append(tmparr, fmt.Sprintf("%s got address %s", d.macaddr, d.ipaddr))
	}

	// And print it as a tree
	return common.GenerateTree(tmparr)
}

// Returns the arpStats object for the given MAC address, or creates a new one
func (p *Protocol) getStatOrCreate(macaddr string) *arpStats {
	// Try to find the given macaddr
	for i := 0; i < len(arpStatsList); i++ {
		if arpStatsList[i].macaddr == macaddr {
			// Found, return it
			return &arpStatsList[i]
		}
	}

	// None found yet, we need to create a new one
	arpStatsList = append(arpStatsList, arpStats{
		macaddr:  macaddr,
	})

	// And return it
	return &arpStatsList[len(arpStatsList)-1]
}

// Adds a new entry to the devices array, checking if there may be a collision (=ARP Spoofing)
func (p *Protocol) addDeviceEntry(macaddr string, ipaddr string) {
	if ipaddr == "0.0.0.0" {
		// Possible ARP request if sender doesn't have an IP address yet. Ignore.
		return
	}

	for i := 0; i < len(devices); i++ {
		// check if we found a collision (possible ARP spoofing)
		if (devices[i].macaddr == macaddr) != (devices[i].ipaddr == ipaddr) {
			// this operation is practically XOR (which golang doesn't provide e.g. with ^)

			// Check if one address is in the link-local block (169.254.0.0/16), ignore "ARP spoofing" then
			if !linkLocalBlock.Contains(net.ParseIP(devices[i].ipaddr)) && !linkLocalBlock.Contains(net.ParseIP(ipaddr)) {
				// The old and the new IP are both outside of the link-local range - we can warn about ARP spoofing
				log.Printf("Found possible ARP spoofing! Old: (MAC=%s, IP=%s), New: (MAC=%s, IP=%s). Overriding...", devices[i].macaddr, devices[i].ipaddr, macaddr, ipaddr)
			}

			devices[i].macaddr = macaddr
			devices[i].ipaddr = ipaddr
			return
		}

		if devices[i].macaddr == macaddr && devices[i].ipaddr == ipaddr {
			// Found collision, but no ARP spoofing (both values are identical)
			return
		}
	}

	// No device found, add a new entry
	devices = append(devices, arpDevice{
		macaddr: macaddr,
		ipaddr:  ipaddr,
	})
}
