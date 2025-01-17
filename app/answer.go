package main

import (
	"encoding/binary"
)

type Answer struct {
	NAME     []string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func Check(q Question) (Answer, bool) {
	return Answer{
		NAME:     q.NAME,
		TYPE:     q.TYPE,
		CLASS:    q.CLASS,
		TTL:      uint32(3600),
		RDLENGTH: uint16(4),
		RDATA:    []byte{'\x08', '\x08', '\x08', '\x08'},
	}, true
}

func (a Answer) Writer() []byte {
	message := make([]byte, 0)
	b := binary.BigEndian
	for _, label := range a.NAME {
		message = append(message, byte(len(label)))
		message = append(message, []byte(label)...)
	}
	message = append(message, '\x00')
	message = b.AppendUint16(message, a.TYPE)
	message = b.AppendUint16(message, a.CLASS)
	message = b.AppendUint32(message, a.TTL)
	message = b.AppendUint16(message, a.RDLENGTH)
	message = append(message, a.RDATA...)

	return message
}
