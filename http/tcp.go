package HTTP

import (
	"net"
)

// HTTPTransport represents the HTTP transport layer for handling peer-to-peer communication.
type HTTPTransport struct {
	HTTPTransportOptions
	listener     net.Listener  // Listener for accepting incoming HTTP connections.
	responseChan chan Response // Channel for receiving responses from remote peers.
}

// NewHTTPTransport creates a new HTTPTransport instance with the given options.
//
// opts: The HTTP transport options that configure the transport behavior.
//
// Returns a pointer to a new HTTPTransport instance.
func NewHTTPTransport(opts HTTPTransportOptions) *HTTPTransport {
	return &HTTPTransport{
		HTTPTransportOptions: opts,
		responseChan:         make(chan Response),
	}
}

// Dial establishes an outbound HTTP connection to the given address and starts handling the connection.
//
// address: The address of the remote peer to connect to.
//
// Returns an error if the dial operation fails.
func (t *HTTPTransport) Dial(address string) error {
	con, err := net.Dial("HTTP", address)
	if err != nil {
		t.Logger.Errorf("HTTP Dial failed for address : %s", address)
		return err
	}

	// Handle the established connection in a goroutine.
	go t.handleConn(con)

	return nil
}

// ListenAndAccept starts listening for incoming HTTP connections on the specified address.
//
// Returns an error if there is an issue initializing the listener or starting the accept loop.
func (t *HTTPTransport) ListenAndAccept() error {
	t.Logger.Infof("Initiating HTTP to listen on %s ", t.ListenAddress)

	// Initialize the listener.
	var err error
	t.listener, err = net.Listen("HTTP", t.ListenAddress)
	if err != nil {
		return err
	}

	t.Logger.Infof("HTTP listen to %s successful", t.ListenAddress)

	// Start the loop for accepting incoming connections asynchronously.
	go t.startAcceptLoop()

	return nil
}

// startAcceptLoop is an infinite loop that accepts incoming HTTP connections from the listener.
//
// It handles each new connection in a separate goroutine.
func (t *HTTPTransport) startAcceptLoop() {
	t.Logger.Info("Starting HTTP accept loop")

	for {
		// Accept a connection from the listener.
		conn, err := t.listener.Accept()
		if err == net.ErrClosed {
			return
		}
		if err != nil {
			t.Logger.Errorf("HTTP accept error: %s", err)
		}

		t.Logger.Infof("Accepted HTTP connection from %s", conn.RemoteAddr())

		// Handle the accepted connection in a separate goroutine.
		go t.handleConn(conn)
	}
}

// handleConn handles the established HTTP connection, performs the handshake, and processes incoming data.
//
// conn: The HTTP connection to the remote peer.
// outbound: A boolean flag indicating whether this connection is outbound (true) or inbound (false).
func (t *HTTPTransport) handleConn(conn net.Conn) {
	var err error

	// Read loop: Continuously receive data from the connection.
	rpc := Response{}
	for {
		// Decode the incoming data into a Response object.
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			t.Logger.Errorf("HTTP failed to decode payload from %+v : %s", conn, err)
			return
		}

		// Set the remote address for the response and send it to the response channel.
		rpc.From = conn.RemoteAddr()
		t.responseChan <- rpc
		t.Logger.Infof("Response: %+v", rpc)
	}
}

// Consume implements the Transport interface, returning a read-only channel to consume incoming responses.
//
// Returns a read-only channel (chan Response) that allows the caller to receive incoming messages from remote peers.
func (t *HTTPTransport) Consume() <-chan Response {
	// The channel is read-only, indicated by <-chan.
	return t.responseChan
}

// Close closes the HTTP listener and performs any necessary cleanup operations.
//
// Returns an error if there is an issue closing the listener.
func (t *HTTPTransport) Close() error {
	t.Logger.Infof("Dropping HTTP connection with %s", t.ListenAddress)
	return t.listener.Close() // Close the HTTP listener
}
