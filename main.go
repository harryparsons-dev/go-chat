package main

import (
	"encoding/binary"
	"fmt"
)

type TcpHandler struct {
	Command byte
	Data    []byte
}

const HEADER_SIZE = 3

func main() {

	var i []byte
	fmt.Println("Chat:")
	fmt.Scan(&i)
	s := &TcpHandler{Command: 1, Data: i}
	packet, _ := s.CreatePacket()
	fmt.Printf("%v\n", packet)
	fmt.Printf("Converting back to string...\n")

	packetString := s.DestructPacket(packet)
	fmt.Println(packetString)

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
