package output

import (
	"fmt"
	"io/ioutil"

	"crypto/sha256"

	"github.com/google/gopacket"
)

var graphPkgs []GraphPkg

// AddPkgToGraph adds the given packet as communication to the graph
func AddPkgToGraph(pkg gopacket.Packet) {
	// Only proceed if pkg contains a network layer
	if pkg.NetworkLayer() == nil {
		return
	}

	src := pkg.NetworkLayer().NetworkFlow().Src().String()
	dst := pkg.NetworkLayer().NetworkFlow().Dst().String()

	// Search for the given communication pair
	for _, p := range graphPkgs {
		if p.from == src && p.to == dst {
			// Communication pair found, add protocol and finish
			p.AddProtocol("nil")
			return
		}
	}

	// Communcation pair was not in graphPkgs, add to it
	graphPkgs = append(graphPkgs, GraphPkg{
		from:     src,
		to:       dst,
		protocol: []string{""},
	})
}

// CreateGraph writes out a Graphviz digraph
func CreateGraph() {
	if *graphOutput == "" {
		// No graph requested
		return
	}

	// Start with the Graphviz-specific header
	dot := fmt.Sprintf("# Compile with `neato -Tpng %s > %s.png`\n", *graphOutput, *graphOutput)
	dot += "digraph pancap {\n\toverlap = false;\n"

	// First, gather all nodes as-is and write them out
	dot += nodedef(graphPkgs)

	// Iterate over communication
	for _, p := range graphPkgs {
		dot += fmt.Sprintf("\tn%s->n%s\n", hash(p.from), hash(p.to))
	}

	// Close
	dot += "}\n"

	// Write out
	ioutil.WriteFile(*graphOutput, []byte(dot), 0644)
}

// Creates a list of distinct nodes, Graphviz-compatible
func nodedef(pkgs []GraphPkg) string {
	output := ""
	nodes := []string{}
	for _, p := range graphPkgs {
		// Check if src and dst are already present in nodes array
		srcFound := false
		dstFound := false
		for _, n := range nodes {
			if p.from == n {
				srcFound = true
			}
			if p.to == n {
				dstFound = true
			}
		}
		if !srcFound {
			// src not yet present, add to node list
			nodes = append(nodes, p.from)
		}
		if !dstFound {
			// dst not yet present, add to node list
			nodes = append(nodes, p.to)
		}
	}

	// Output Graphviz-compatible node definition
	for _, n := range nodes {
		// As the Graphviz charset for nodes is rather small, rely on hashes
		output += fmt.Sprintf("\tn%s[label=\"%s\"]\n", hash(n), n)
	}

	return output
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))[:6]
}
