package dhcpv4

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

var (
	watchedOpts = []layers.DHCPOpt{
		layers.DHCPOptSubnetMask, // Option 1
		layers.DHCPOptRouter, // Option 3
		layers.DHCPOptDNS, // Option 6
		layers.DHCPOptBroadcastAddr, // Option 28
		layers.DHCPOptNTPServers, // Option 42
		layers.DHCPOptLeaseTime, // Option 51
		layers.DHCPOptT1, // Option 58
	}
)

// Generates the summary of relevant DHCP options
func (p *Protocol) generateNetworkSummary() string {
	// It's also possible to use strings.Builder here, but it produces code which is longer than this solution :shrug:
	summary := fmt.Sprintf("Subnet Mask: %s\n", formatIP(p.networkSetup[layers.DHCPOptSubnetMask]))
	summary = fmt.Sprintf("%sBroadcast: %s\n", summary, formatIP(p.networkSetup[layers.DHCPOptBroadcastAddr]))
	summary = fmt.Sprintf("%sRouter: %s\n", summary, formatIP(p.networkSetup[layers.DHCPOptRouter]))
	summary = fmt.Sprintf("%sDNS Server: %s\n", summary, formatIP(p.networkSetup[layers.DHCPOptDNS]))
	summary = fmt.Sprintf("%sNTP Server: %s\n", summary, formatIP(p.networkSetup[layers.DHCPOptNTPServers]))
	summary = fmt.Sprintf( "%sLease Time: %s\n", summary, formatDate(p.networkSetup[layers.DHCPOptLeaseTime]))
	summary = fmt.Sprintf("%sRenewal Time: %s\n", summary, formatDate(p.networkSetup[layers.DHCPOptT1]))
	return summary
}

// Looks for information specifying the setup of the network. This includes
//  - Option  1: Subnet Mask
//  - Option  3: Router address
//  - Option  6: Domain Name Server address
//  - Option 28: Broadcast address
//  - Option 42: NTP Server address
//  - Option 51: IP Address Lease time
//  - Option 58: IP Renewal time
func (p *Protocol) checkForNetworkInfos(dhcppacket layers.DHCPv4) {
	// Check if it is a DHCP request
	if dhcppacket.Operation == layers.DHCPOpRequest {
		// We can ignore requests, they won't help us here
		return
	}

	// Search for different options (1, 3, 6, 28, 42, 51, 58) in DHCP Packet Options
	for _, o := range dhcppacket.Options {
		if isRelevantOption(o) {
			// Found DHCP option to be watched, let's watch it
			p.saveOption(o)
		}
	}

}

// Saves the given option in the networkSetup map, and informs the user if the value changes
func (p *Protocol) saveOption(opt layers.DHCPOption) {
	// check if we already stored this value
	if p.networkSetup[opt.Type] != nil {
		// We already stored a value, let's check if it's the same as the new one
		if !bytes.Equal(p.networkSetup[opt.Type], opt.Data) {
			// Already stored a value and it's different from our new value - inform user and overwrite value later
			log.Printf("Received different values for DHCP Option %s (ID %d). (Old: %s, New. %s)", opt.Type.String(), opt.Type, p.networkSetup[opt.Type], opt.Data)
		} else {
			// Exactly this value was already stored, no need to overwrite it
			return
		}
	}

	p.networkSetup[opt.Type] = opt.Data
}

// Checks if the given DHCPOption is part of the watchlist
func isRelevantOption(opt layers.DHCPOption) bool {
	// Iterate over all DHCP options in our watchlist
	for _, o := range watchedOpts {
		if o == opt.Type {
			// Found.
			return true
		}
	}

	// This option is not on our watchlist.
	return false
}

// Formats the given byte array as string representing the IP address, or returns an error (as string)
func formatIP(rawIP []byte) string {
	// Check if we even have an IP
	if rawIP == nil {
		// We don't have an IP, construct an error message (as string)
		error := color.New(color.FgRed)
		return error.Sprint("(not found)")
	}

	// Return formatted IP
	return net.IP(rawIP).String()
}

func formatDate(rawDate []byte) string {
	// Check if we even have a date
	if rawDate == nil {
		// We don't have a date, construct an error message (as string)
		error := color.New(color.FgRed)
		return error.Sprint("(not found)")
	}

	// Actually format date
	intDate := binary.LittleEndian.Uint32(rawDate)
	seconds := intDate % 60
	minutes := intDate / 60 % 60
	hours   := intDate / 60 / 60 % 60
	formattedDate := ""

	// Check which words we need to pick
	// ... regarding hours
	if hours > 0 {
		formattedDate = fmt.Sprintf("%d hours", hours)
	}

	// ... regarding minutes
	if minutes > 0 {
		// check if we got a previous string we need to take care of
		if len(formattedDate) > 0 {
			// yes, append our information to existing string
			formattedDate = fmt.Sprintf("%s, %d minutes", formattedDate, minutes)
		} else {
			// no, use our string
			formattedDate = fmt.Sprintf("%d minutes", minutes)
		}
	}

	// ... regarding seconds
	if seconds > 0 {
		// check if we got a previous string we need to take care of
		if len(formattedDate) > 0 {
			// yes, append our information to existing string
			formattedDate = fmt.Sprintf("%s, %d seconds", formattedDate, seconds)
		} else {
			// no, use our string
			formattedDate = fmt.Sprintf("%d seconds", seconds)
		}
	}

	return formattedDate
}
