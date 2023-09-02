package dhcpv4

type hostname struct {
	hostname       string
	requestedByMAC string
	granted        bool
	deniedHostname string
}
