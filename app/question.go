package main

import (
	"bytes"
	"encoding/binary"
)

type Question struct {
	LABELS []string
	TYPE   uint16
	CLASS  uint16
}

func ParseQuestions(message []byte, count uint16) []Question {
	b := binary.BigEndian
	questions := make([]Question, 0, count)
	for i := 0; i < int(count); i++ {
		p := bytes.IndexByte(message, '\x00')
		questions = append(questions, Question{
			LABELS: parseLabels(message[:p]),
			TYPE:   b.Uint16(message[p+1 : p+3]),
			CLASS:  b.Uint16(message[p+3 : p+5]),
		})
		message = message[p+5:]
	}
	return questions
}

func parseLabels(message []byte) []string {
	labels := make([]string, 0)
	p := 0
	for p < len(message) {
		length := int(message[p])
		labels = append(labels, string(message[p+1:p+1+length]))
		p += 1 + length
	}
	return labels
}

func (q Question) Writer() []byte {
	message := make([]byte, 0)
	b := binary.BigEndian
	for _, label := range q.LABELS {
		message = append(message, byte(len(label)))
		message = append(message, []byte(label)...)
	}
	message = append(message, '\x00')
	message = b.AppendUint16(message, q.TYPE)
	message = b.AppendUint16(message, q.CLASS)
	return message
}
