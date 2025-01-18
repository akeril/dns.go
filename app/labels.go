package main

import (
	"encoding/binary"
)

func parseLabels(message []byte, l int) (int, []string) {
	labels := make([]string, 0)
	for l < len(message) {
		if message[l] == '\x00' {
			l++
			break
		}
		if index, ok := pointer(message[l : l+2]); ok {
			_, lbls := parseLabels(message, index)
			labels = append(labels, lbls...)
			l += 2
			break
		}
		length := int(message[l])
		labels = append(labels, string(message[l+1:l+1+length]))
		l += 1 + length
	}
	return l, labels
}

func pointer(prefix []byte) (int, bool) {
	if (prefix[0]&0x80) == 0 || (prefix[0]&0x40) == 0 {
		return 0, false
	}
	index := binary.BigEndian.Uint16(prefix) ^ 0xC000
	return int(index), true
}
