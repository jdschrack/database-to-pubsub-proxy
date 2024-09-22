package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	postgresHost = "localhost:5432" // Replace with your actual PostgreSQL host and port
)

// Start the proxy
func main() {
	// Create the listener for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:15432")
	if err != nil {
		log.Fatalf("Error starting proxy: %v", err)
	}
	log.Printf("Proxy listening on %s", listener.Addr())

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from: %s", clientConn.RemoteAddr())

		go handleConnection(clientConn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	log.Println("Reading data sent by client...")

	clientReader := bufio.NewReader(clientConn)

	// Read the first 8 bytes to detect SSL request
	buf := make([]byte, 8)
	_, err := clientReader.Read(buf)
	if err != nil {
		log.Printf("Error reading from client: %v", err)
		return
	}

	// Check if it's an SSL request (PostgreSQL sends SSL request with a magic number)
	if buf[4] == 0x04 && buf[5] == 0xd2 && buf[6] == 0x16 && buf[7] == 0x2f {
		log.Println("Client requested SSL, responding with 'S'")

		// Respond with 'S' to initiate SSL handshake
		_, err = clientConn.Write([]byte("S"))
		if err != nil {
			log.Printf("Error sending SSL response: %v", err)
			return
		}

		// Upgrade connection to SSL
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{loadTLSCertificate()},
		}
		clientConn = tls.Server(clientConn, tlsConfig)
		clientReader = bufio.NewReader(clientConn)
		log.Println("SSL handshake completed with client.")
	} else {
		log.Println("No SSL requested, proceeding with plaintext connection.")
	}

	// Connect to PostgreSQL server
	serverConn, err := net.Dial("tcp", postgresHost)
	if err != nil {
		log.Printf("Error connecting to PostgreSQL server: %v", err)
		return
	}
	defer serverConn.Close()

	serverReader := bufio.NewReader(serverConn)

	// Relay traffic between client and PostgreSQL server, intercept SQL commands
	go forwardTraffic(
		clientConn,
		serverConn,
		clientReader,
		"client",
		true,
	) // From client to server (intercept and log SQL)
	forwardTraffic(
		serverConn,
		clientConn,
		serverReader,
		"server",
		false,
	) // From server to client (forward responses)
}

// Load TLS certificate for SSL connection
func loadTLSCertificate() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Error loading TLS certificate: %v", err)
	}
	return cert
}

// Relay traffic between client and server, with option to intercept SQL commands
func forwardTraffic(
	src net.Conn,
	dst net.Conn,
	reader *bufio.Reader,
	direction string,
	intercept bool,
) {
	for {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from %s: %v", direction, err)
			}
			return
		}

		if intercept && direction == "client" {
			logSQLCommand(buf[:n]) // Intercept and log SQL commands sent by the client
		}

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("Error writing to %s: %v", direction, err)
			return
		}
	}
}

// Log intercepted SQL commands
func logSQLCommand(data []byte) {
	logFile, err := os.OpenFile("sql_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	logData := fmt.Sprintf("Intercepted SQL Command: %s\n", string(data))
	_, err = logFile.WriteString(logData)
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}

	log.Printf("Intercepted SQL Command: %s", string(data))
}
