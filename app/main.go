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
		response := Resolve(request)
		_, err = conn.WriteToUDP(response.Writer(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func Resolve(req DNS) DNS {
	resp := New()

	// set Headers
	resp.Header().ID = req.Header().ID
	resp.Header().Set("QR", QR_BIT)
	if req.Header().Get("RD") != 0 {
		resp.Header().Set("RD", 0)
	}
	resp.Header().Set("OPCODE", req.Header().Get("OPCODE"))
	if req.Header().Get("OPCODE") != 0 {
		resp.Header().Set("RCODE", 4)
	}

	// resolve DNS queries
	for _, question := range req.questions {
		resp.questions = append(resp.questions, question)
		resp.Header().QDCOUNT++
		answer, ok := Check(question)
		if ok {
			resp.answers = append(resp.answers, answer)
			resp.Header().ANCOUNT++
		}
	}

	return resp
}
