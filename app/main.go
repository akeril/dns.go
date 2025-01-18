package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	resolveIP := flag.String("resolver", "", "DNS Resolution Server")
	flag.Parse()

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
		response := Resolve(*resolveIP, request)

		_, err = conn.WriteToUDP(response.Writer(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func Resolve(IP string, req DNS) DNS {
	resp := New()

	resp.Header().ID = req.Header().ID
	resp.Header().Set("QR", QR_BIT)
	if req.Header().Get("RD") != 0 {
		resp.Header().Set("RD", req.Header().Get("RD"))
	}
	resp.Header().Set("OPCODE", req.Header().Get("OPCODE"))
	if req.Header().Get("OPCODE") != 0 {
		resp.Header().Set("RCODE", 4)
	}
	resp.Header().QDCOUNT = req.Header().QDCOUNT
	resp.questions = req.questions

	for _, question := range req.questions {
		// fwdResp, err := Check(IP, req, question)
		answer, err := Check(IP, question)
		if err == nil {
			resp.answers = append(resp.answers, answer)
			resp.Header().ANCOUNT++
		}
	}

	return resp
}
