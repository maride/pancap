package output

// GraphPkg resembles a directed communication from one address to another
// It wraps up required information to draw a graph of the communication, including spoken protocols.
type GraphPkg struct {
	from     string
	to       string
	protocol []string
}

// AddProtocol adds the given protocol to the list of protocols if not already present
func (p *GraphPkg) AddProtocol(protocol string) {
	for _, p := range p.protocol {
		if p == protocol {
			return
		}
	}
	p.protocol = append(p.protocol, protocol)
}
