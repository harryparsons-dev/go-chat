package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type TcpHandler struct {
	Command byte
	Data    []byte
}

type Server struct {
	listenAddr *net.UDPAddr
	ch         chan struct{}
}

func NewServer() *Server {
	return &Server{
		listenAddr: &net.UDPAddr{
			Port: 3000,
			IP:   net.ParseIP("127.0.0.1"),
		},
		ch: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	conn, err := net.ListenUDP("udp", s.listenAddr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return err
	}
	defer conn.Close()
	go s.Read(conn)

	<-s.ch
	return nil
}

func (S *Server) Read(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}

		go S.handleMessages(buf[:n])

	}
}

func (S *Server) handleMessages(msg []byte) {
	time.Sleep(time.Second * 5)
	fmt.Printf("Message: %v", string(msg))
}

const HEADER_SIZE = 3

func main() {

	server := NewServer()
	log.Fatal(server.Start())

	// var i []byte
	// fmt.Println("Chat:")
	// fmt.Scan(&i)
	// s := &TcpHandler{Command: 1, Data: i}
	// packet, _ := s.CreatePacket()
	// fmt.Printf("%v\n", packet)
	// fmt.Printf("Converting back to string...\n")

	// packetString := s.DestructPacket(packet)
	// fmt.Println(packetString)

}

func (tcp *TcpHandler) CreatePacket() (packet []byte, err error) {
	length := uint16(len(tcp.Data))
	lengthData := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthData, length)

	fmt.Print(length)

	p := make([]byte, 0, 2+2+length)
	p = append(p, tcp.Command)
	p = append(p, lengthData...)
	p = append(p, tcp.Data...)

	return p, nil

}

func (tcp *TcpHandler) DestructPacket(packet []byte) string {
	command := packet[:1]
	length := int(binary.BigEndian.Uint16(packet[1:4]))

	end := HEADER_SIZE + length

	data := string(packet[HEADER_SIZE:end])

	packetString := fmt.Sprintf("COMMAND: %v\nLENGTH: [%v]\nDATA: %v", command, length, data)

	return packetString

}

func client() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println(err)
	}

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Message: %v", p)
	}
	conn.Close()
}

func server() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("127.0.0.1"),
	}

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your message "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
