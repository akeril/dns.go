package main

import "encoding/binary"

const (
	QR_BIT = 0b1000000000000000
	OP_BIT = 0b0111100000000000
	AA_BIT = 0b0000010000000000
	TC_BIT = 0b0000001000000000
	RD_BIT = 0b0000000100000000
	RA_BIT = 0b0000000010000000
	RZ_BIT = 0b0000000001110000
	RC_BIT = 0b0000000000001111
)

type Header struct {
	ID      uint16
	FLAGS   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h Header) Get(flag string) uint16 {
	switch flag {
	case "QR":
		return h.FLAGS & QR_BIT
	case "OPCODE":
		return h.FLAGS & OP_BIT
	case "AA":
		return h.FLAGS & AA_BIT
	case "TC":
		return h.FLAGS & TC_BIT
	case "RD":
		return h.FLAGS & RD_BIT
	case "RA":
		return h.FLAGS & RA_BIT
	case "RESERVED":
		return h.FLAGS & RZ_BIT
	case "RCODE":
		return h.FLAGS & RC_BIT
	default:
		panic("Unknown Flag")
	}
}

// value bits are only considered for "OPCODE", "RESERVED", "RCODE"
func (h *Header) Set(flag string, value uint16) {
	switch flag {
	case "QR":
		h.FLAGS |= QR_BIT
	case "OPCODE":
		h.FLAGS |= (OP_BIT & value)
	case "AA":
		h.FLAGS |= AA_BIT
	case "TC":
		h.FLAGS |= TC_BIT
	case "RD":
		h.FLAGS |= RD_BIT
	case "RA":
		h.FLAGS |= RA_BIT
	case "RESERVED":
		h.FLAGS |= (RZ_BIT & value)
	case "RCODE":
		h.FLAGS |= (RC_BIT & value)
	default:
		panic("Unknown Flag")
	}
}

func ParseHeader(message []byte) Header {
	b := binary.BigEndian
	return Header{
		ID:      b.Uint16(message[:2]),
		FLAGS:   b.Uint16(message[2:4]),
		QDCOUNT: b.Uint16(message[4:6]),
		ANCOUNT: b.Uint16(message[6:8]),
		NSCOUNT: b.Uint16(message[8:10]),
		ARCOUNT: b.Uint16(message[10:12]),
	}
}

func (h Header) Writer() []byte {
	header := make([]byte, 0, 12)

	header = binary.BigEndian.AppendUint16(header, h.ID)
	header = binary.BigEndian.AppendUint16(header, h.FLAGS)
	header = binary.BigEndian.AppendUint16(header, h.QDCOUNT)
	header = binary.BigEndian.AppendUint16(header, h.ANCOUNT)
	header = binary.BigEndian.AppendUint16(header, h.NSCOUNT)
	header = binary.BigEndian.AppendUint16(header, h.ARCOUNT)

	return header
}
