package http

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
	"net/http"
)

var (
	requestSummaryLines []string
)

type httpRequestFactory struct{}

type httpRequestStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// Creates a new HTTPRequestStream for the given packet flow, and analyzes it in a separate thread
func (h *httpRequestFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpRequestStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}

	// Start analyzer as thread and return TCP reader stream
	go hstream.run()
	return &hstream.r
}

// Analyzes the given request
func (h *httpRequestStream) run() {
	iobuf := bufio.NewReader(&h.r)

	for {
		req, reqErr := http.ReadRequest(iobuf)

		if reqErr == io.EOF {
			// That's ok, we can ignore EOF errors
			return
		} else if reqErr != nil {
			// Ignore, because it may be a response
		} else {
			// Try to process assembled request
			tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()

			// Build summary
			line := fmt.Sprintf("Request %s http://%s%s", req.Method, req.Host, req.RequestURI)
			requestSummaryLines = append(requestSummaryLines, line)

			// Check for file uploads
			if req.MultipartForm != nil && req.MultipartForm.File != nil {
				for k, v := range req.MultipartForm.File {
					log.Println(k, v)
				}
			}
		}
	}
}
