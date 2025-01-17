package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 512)
	for {
		size, source, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		request := Parse(buf[:size])
		response := New(request.Header().ID)
		response.Header().Set("QR")
		response.questions = request.questions
		response.Header().QDCOUNT = request.Header().QDCOUNT

		for _, question := range request.questions {
			answer, ok := Check(question)
			if ok {
				response.answers = append(response.answers, answer)
				response.Header().ANCOUNT++
			}
		}

		_, err = conn.WriteToUDP(response.Writer(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
