package input

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/graphite-ng/carbon-relay-ng/cfg"
)

var (
	addr   = "localhost"
	config = cfg.Config{}
)

type mockHandler struct {
	sync.Mutex
	data    []byte
	testing *testing.T
}

// just store all the data in m.data
func (m *mockHandler) Handle(r io.Reader) {
	buf := make([]byte, 100)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			return
		}
		m.Lock()
		m.data = append(m.data, buf[:n]...)
		m.Unlock()
	}
}

func (m *mockHandler) String() string {
	m.Lock()
	defer m.Unlock()
	return string(m.data)
}

func TestTcpUdpShutdown(t *testing.T) {
	handler := mockHandler{testing: t}
	addr := "localhost:" // choose random ports
	listener := NewListener("mock", &handler)
	err := listener.listen(addr)
	if err != nil {
		t.Fatalf("Error when trying to listen: %s", err)
	}
	res := listener.Stop()
	if !res {
		t.Fatalf("Failed to shut down cleanly")
	}
}

func TestTcpConnection(t *testing.T) {
	handler := mockHandler{testing: t}
	addr := "localhost:" // choose random ports
	listener := NewListener("mock", &handler)
	err := listener.listen(addr)
	if err != nil {
		t.Fatalf("Error when listening: %s", err)
	}

	r_addr := listener.tcpList.Addr()
	conn, err := net.DialTCP("tcp", nil, r_addr.(*net.TCPAddr))
	if err != nil {
		t.Fatalf("Error when connecting to listening port: %s", err)
	}

	testContent := "test"
	_, err = conn.Write([]byte(testContent))

	// with carbon there's no way of knowing if the server has processed the data
	// we've sent, so we just give it 50ms and then assume it did.
	time.Sleep(time.Millisecond * 50)

	listener.Stop()

	received := handler.String()
	if received != testContent {
		t.Fatalf("Received unexpected content in handler. Expected \"%s\" got \"%s\"", testContent, received)
	}

	// giving the server another 50ms to shut down
	time.Sleep(time.Millisecond * 50)

	// verify that the socket is closed now
	_, err = net.DialTCP("tcp", nil, r_addr.(*net.TCPAddr))
	if err == nil {
		t.Fatalf("Connection to tcp server should have failed, but it did not")
	}
}

func TestUdpConnection(t *testing.T) {
	handler := mockHandler{testing: t}
	addr := "localhost:" // choose random ports
	listener := NewListener("mock", &handler)
	err := listener.listen(addr)
	if err != nil {
		t.Fatalf("Error when listening: %s", err)
	}

	r_addr := listener.udpConn.LocalAddr()
	conn, err := net.DialUDP("udp", nil, r_addr.(*net.UDPAddr))
	if err != nil {
		t.Fatalf("Error when connecting to listening port: %s", err)
	}

	testContent := "test"
	_, err = conn.Write([]byte(testContent))

	// with carbon there's no way of knowing if the server has processed the data
	// we've sent, so we just give it 50ms and then assume it did.
	time.Sleep(time.Millisecond * 50)

	listener.Stop()

	received := handler.String()
	if received != testContent {
		t.Fatalf("Received unexpected content in handler. Expected \"%s\" got \"%s\"", testContent, received)
	}

	// giving the server another 50ms to shut down
	time.Sleep(time.Millisecond * 50)

	buffer := make([]byte, 10)
	listener.udpConn.SetDeadline(time.Now().Add(time.Second))
	_, _, err = listener.udpConn.ReadFrom(buffer)
	if err == nil {
		t.Fatalf("Expected read from udp connection to fail, but it did not")
	}
	if err.(*net.OpError).Timeout() {
		t.Fatalf("Expected i/o error, but got timeout error")
	}
}
