package main

import (
	"errors"
	"net"
)

func Check(IP string, q Question) (Answer, error) {
	if IP == "" {
		return Answer{
			NAME:     q.NAME,
			TYPE:     q.TYPE,
			CLASS:    q.CLASS,
			TTL:      uint32(3600),
			RDLENGTH: uint16(4),
			RDATA:    []byte{'\x08', '\x08', '\x08', '\x08'},
		}, nil
	}
	req := New()
	req.questions = []Question{q}
	req.Header().QDCOUNT = 1
	buf, err := Forward(IP, req.Writer())
	if err != nil {
		return Answer{}, err
	}
	resp := Parse(buf)
	if len(resp.answers) != 1 {
		return Answer{}, errors.New("No response from forwarding server")
	}
	return resp.answers[0], nil
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
