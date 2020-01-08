package http

import (
	"bufio"
	"fmt"
	"git.darknebu.la/maride/pancap/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	responseSummaryLines   []string
)

type httpResponseFactory struct{}

type httpResponseStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// Creates a new HTTPResponseStream for the given packet flow, and analyzes it in a separate thread
func (h *httpResponseFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpResponseStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

// Analyzes the given response
func (h *httpResponseStream) run() {
	iobuf := bufio.NewReader(&h.r)

	for {
		resp, respErr := http.ReadResponse(iobuf, nil)

		if respErr == io.EOF {
			// That's ok, we can ignore EOF errors
			return
		} else if respErr != nil {
			// Ignore, because it may be a request
		} else {
			// Try to process assembled request
			fileBytes, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			// Register file in filemanager
			output.RegisterFile("", fileBytes, "HTTP response")

			// Build summary
			line := fmt.Sprintf("Response %s, Type %s, Size %d bytes", resp.Status, resp.Header.Get("Content-Type"), resp.ContentLength)
			responseSummaryLines = append(responseSummaryLines, line)
		}
	}
}
