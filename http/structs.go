package HTTP

import (
	"net"

	"github.com/sumit-behera-in/goLogger"
)

type HTTPTransportOptions struct {
	Logger        *goLogger.Logger // Logger is the logger instance
	ListenAddress string           // ListenAddress is the address on which the transport listens
	Decoder       Decoder          // Decoder is the decoder that is used to decode the incoming data
}

// Response hold arbitrary data that is being sent over each transport between two nodes in a network.
type Response struct {
	From    net.Addr // From is the address of the node that sent the data
	Payload []byte   // Payload is the data that is being sent over the network
}
