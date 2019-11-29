package dhcpv4

import (
	"fmt"
	"git.darknebu.la/maride/pancap/common"
	"github.com/google/gopacket/layers"
	"log"
)

var (
	hostnames []hostname
)

func checkForHostname(dhcppacket layers.DHCPv4) {
	// Search for "Hostname" option (ID 12) in DHCP Packet Options
	for _, o := range dhcppacket.Options {
		if o.Type == layers.DHCPOptHostname {
			// found it. Let's see if it's a request or response
			if dhcppacket.Operation == layers.DHCPOpRequest {
				// request, not granted yet.
				addHostname(hostname{
					hostname:       string(o.Data),
					requestedByMAC: dhcppacket.ClientHWAddr.String(),
					granted:        false,
				})
			} else {
				// Response, DHCP issued this hostname
				addHostname(hostname{
					hostname:       string(o.Data),
					requestedByMAC: "",
					granted:        true,
				})
			}

			return
		}
	}

	// None found, means client or server doesn't support Hostname option field. Ignore.
}

// Prints the list of all hostnames encountered.
func printHostnames() {
	var tmparr []string

	// Construct meaningful text
	for _, h := range hostnames {
		answer := ""

		// check what kind of answer we need to construct
		if h.deniedHostname == "" {
			// Hostname was not denied, let's check if it was officially accepted
			if h.granted {
				// it was. Yay.
				answer = fmt.Sprintf("%s has hostname %s, granted by the DHCP server", h.requestedByMAC, h.hostname)
			} else {
				// it was neither denied nor accepted, either missing the DHCP answer in capture file or misconfigured DHCP server
				answer = fmt.Sprintf("%s has hostname %s, without a response from DHCP server", h.requestedByMAC, h.hostname)
			}
		} else {
			// Hostname was denied, let's check if we captured the request
			if h.hostname == "" {
				// we didn't.
				answer = fmt.Sprintf("%s was forced to have hostname %s by DHCP server,", h.requestedByMAC, h.hostname)
			} else {
				// we did, print desired and de-facto hostname
				answer = fmt.Sprintf("%s asked for hostname %s, but got hostname %s from DHCP server.", h.requestedByMAC, h.deniedHostname, h.hostname)
			}
		}

		tmparr = append(tmparr, answer)
	}

	// and print it as a tree.
	common.PrintTree(tmparr)
}

// Adds the given hostname to the hostname array, or patches an existing entry if found
func addHostname(tmph hostname) {
	// see if we have an existing entry for this hostname
	for i := 0; i < len(hostnames); i++ {
		// get ith hostname in the list
		h := hostnames[i]

		// ... and check if it's the one requested
		if h.hostname == tmph.hostname {
			// Found hostname, check different possible cases
			if tmph.requestedByMAC != "" {
				// Already got that hostname in the list, but received another request for it
				if tmph.requestedByMAC == h.requestedByMAC {
					// Same client asked for the same hostname - that's ok. Ignore.
				} else {
					// Different devices asked for the same hostname - log it.
					log.Printf("Multiple clients (%s, %s) asked for the same hostname (%s)", h.requestedByMAC, tmph.requestedByMAC, h.hostname)
				}
			} else {
				// Received a response for this hostname, check if it was granted
				if h.hostname == tmph.hostname {
					// granted, everything is fine.
					hostnames[i].granted = true
				} else {
					// Received a different hostname than the one requested by the MAC. Report that.
					log.Printf("Client %s asked for hostname '%s' but was given '%s' by DHCP server", h.requestedByMAC, tmph.hostname, h.hostname)
					hostnames[i].deniedHostname = hostnames[i].hostname
					hostnames[i].hostname = tmph.hostname
					hostnames[i].granted = false
				}
				// in either case, it's a response by the DHCP server - hostname is granted in this context

			}

			return
		}
	}

	// We didn't find the desired hostname, append given object to the list
	hostnames = append(hostnames, tmph)
}
