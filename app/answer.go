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

func ParseAnswers(message []byte, start int, count uint16) (int, []Answer) {
	b := binary.BigEndian
	answers := make([]Answer, 0)
	p := start
	for i := 0; p < len(message) && i < int(count); i++ {
		delim, labels := parseLabels(message, p)
		length := b.Uint16(message[delim+8 : delim+10])
		answers = append(answers, Answer{
			NAME:     labels,
			TYPE:     b.Uint16(message[delim : delim+2]),
			CLASS:    b.Uint16(message[delim+2 : delim+4]),
			TTL:      b.Uint32(message[delim+4 : delim+8]),
			RDLENGTH: length,
			RDATA:    message[delim+10 : delim+10+int(length)],
		})
		p = delim + 10 + int(length)
	}
	return p, answers
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
