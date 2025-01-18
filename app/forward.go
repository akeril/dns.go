package main

import (
	"fmt"
	"net"
)

func Check(IP string, req DNS, q Question) (DNS, error) {
	req.questions = []Question{q}
	req.Header().QDCOUNT = 1
	buf, err := Forward(IP, req.Writer())
	if err != nil {
		fmt.Println("Forwarding error: ", err)
		return DNS{}, err
	}
	resp := Parse(buf)
	return resp, nil
}

// forwards DNS request to another server
func Forward(ip string, message []byte) ([]byte, error) {
	addr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		return []byte{}, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return []byte{}, err
	}

	_, err = conn.Write(message)
	if err != nil {
		return []byte{}, err
	}

	buf := make([]byte, 512)
	size, err := conn.Read(buf)
	if err != nil {
		return []byte{}, err
	}

	return buf[:size], nil
}
