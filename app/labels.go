package main

import (
	"encoding/binary"
)

// parses a sequence of labels
// assumes l is at the start of a sequence
// if l points to a pointer, it jumps and retrieves the sequence from that index.
func parseLabels(message []byte, l int) (int, []string) {
	labels := make([]string, 0)
	for l < len(message) {
		if message[l] == '\x00' {
			l++ // on terminating the loop, l should always point to the index immediately after the sequence of labels
			break
		}
		if index, ok := pointer(message[l : l+2]); ok {
			_, lbls := parseLabels(message, index)
			labels = append(labels, lbls...)
			l += 2 // on terminating the loop, l should always point to the index immediately after the sequence of labels
			break
		}
		length := int(message[l])
		labels = append(labels, string(message[l+1:l+1+length]))
		l += 1 + length // this could either point to the next label or terminating character
	}
	return l, labels
}

// tests if current label is a pointer
func pointer(prefix []byte) (int, bool) {
	if (prefix[0]&0x80) == 0 || (prefix[0]&0x40) == 0 {
		return 0, false
	}
	index := binary.BigEndian.Uint16(prefix) ^ 0xC000
	return int(index), true
}
