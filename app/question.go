package main

import (
	"bytes"
	"encoding/binary"
)

type Question struct {
	NAME  []string
	TYPE  uint16
	CLASS uint16
}

func ParseQuestions(message []byte, count uint16) []Question {
	b := binary.BigEndian
	cache := make(map[int]string)
	questions := make([]Question, 0, count)
	p := 12 // offset header bytes
	for i := 0; p < len(message) && i < int(count); i++ {
		delim := p + bytes.IndexByte(message[p:], '\x00')
		questions = append(questions, Question{
			NAME:  parseLabels(message, p, delim, cache),
			TYPE:  b.Uint16(message[delim+1 : delim+3]),
			CLASS: b.Uint16(message[delim+3 : delim+5]),
		})
		p = delim + 5
	}
	return questions
}

func parseLabels(message []byte, l, r int, cache map[int]string) []string {
	labels := make([]string, 0)
	p := l
	for p < r {
		if index, ok := pointer(message[p : p+2]); ok {
			labels = append(labels, cache[index])
			p += 2
			continue
		}
		length := int(message[p])
		cache[p] = string(message[p+1 : p+1+length])
		labels = append(labels, cache[p])
		p += 1 + length
	}
	return labels
}

func pointer(prefix []byte) (int, bool) {
	if (prefix[0]&0x80) == 0 || (prefix[0]&0x40) == 0 {
		return 0, false
	}
	index := binary.BigEndian.Uint16(prefix) ^ 0xC000
	return int(index), true
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
