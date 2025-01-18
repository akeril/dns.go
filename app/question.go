package main

import (
	"encoding/binary"
)

type Question struct {
	NAME  []string
	TYPE  uint16
	CLASS uint16
}

func ParseQuestions(message []byte, start int, count uint16) (int, []Question) {
	b := binary.BigEndian
	questions := make([]Question, 0, count)
	p := start
	for i := 0; p < len(message) && i < int(count); i++ {
		delim, labels := parseLabels(message, p)
		questions = append(questions, Question{
			NAME:  labels,
			TYPE:  b.Uint16(message[delim : delim+2]),
			CLASS: b.Uint16(message[delim+2 : delim+4]),
		})
		p = delim + 4
	}
	return p, questions
}

func (q Question) Writer() []byte {
	message := make([]byte, 0)
	b := binary.BigEndian
	for _, label := range q.NAME {
		message = append(message, byte(len(label)))
		message = append(message, []byte(label)...)
	}
	message = append(message, '\x00')
	message = b.AppendUint16(message, q.TYPE)
	message = b.AppendUint16(message, q.CLASS)
	return message
}
