package arp

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

var (
	arpStatsList []arpStats
	devices []arpDevice
)

// Called on every ARP packet
func ProcessARPPacket(packet gopacket.Packet) error {
	var arppacket layers.ARP

	// Decode raw packet into ARP
	decodeErr := arppacket.DecodeFromBytes(packet.Layer(layers.LayerTypeARP).LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Convert MAC address byte array to string
	sourceAddr := net.HardwareAddr(arppacket.SourceHwAddress).String()
	participant := getStatOrCreate(sourceAddr)

	// Raise stats
	if arppacket.Operation == layers.ARPRequest {
		// Request packet
		participant.asked++
		appendIfUnique(net.IP(arppacket.DstProtAddress).String(), participant.askedList)

		// Add device entry
		addDeviceEntry(sourceAddr, net.IP(arppacket.SourceProtAddress).String())
	} else {
		// Response packet
		participant.answered++
		appendIfUnique(net.IP(arppacket.SourceProtAddress).String(), participant.answeredList)

		// Add device entry
		addDeviceEntry(sourceAddr, net.IP(arppacket.SourceProtAddress).String())
	}

	return nil
}

// Print a summary after all packets are processed
func PrintARPSummary() {
	headline := color.New(color.FgRed, color.Bold)
	headline.Println("ARP traffic summary")
	printTrafficStats()
	headline.Println("ARP LAN overview")
	printLANOverview()
}

// Constructs an answer regarding the ARP traffic
func printTrafficStats() {
	var tmparr []string

	// Iterate over all participants
	for _, p := range arpStatsList {
		tmparr = append(tmparr, fmt.Sprintf("%s asked for %d addresses and answered %d requests", p.macaddr, p.asked, p.answered))
	}

	// And print it as a tree
	printTree(tmparr)
}

// Prints an overview over all connected devices in the LAN
func printLANOverview() {
	var tmparr []string

	// iterate over all devices
	for _, d := range devices {
		tmparr = append(tmparr, fmt.Sprintf("%s got address %s", d.macaddr, d.ipaddr))
	}

	// And print it as a tree
	printTree(tmparr)
}

// Returns the arpStats object for the given MAC address, or creates a new one
func getStatOrCreate(macaddr string) *arpStats {
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
func addDeviceEntry(macaddr string, ipaddr string) {
	if ipaddr == "0.0.0.0" {
		// Possible ARP request if sender doesn't have an IP address yet. Ignore.
		return
	}

	for i := 0; i < len(devices); i++ {
		// check if we found a collision (possible ARP spoofing)
		if (devices[i].macaddr == macaddr) != (devices[i].ipaddr == ipaddr) {
			// this operation is practically XOR (which golang doesn't provide e.g. with ^)
			log.Printf("Found possible ARP spoofing! Old: (MAC=%s, IP=%s), New: (MAC=%s, IP=%s). Overriding...", devices[i].macaddr, devices[i].ipaddr, macaddr, ipaddr)
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
