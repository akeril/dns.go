package main

import "fmt"

type DNS struct {
	header    *Header
	questions []Question
	answers   []Answer
}

func New() DNS {
	return DNS{
		header:    &Header{},
		questions: make([]Question, 0),
		answers:   make([]Answer, 0),
	}
}
func (d DNS) Header() *Header {
	return d.header
}

func Parse(message []byte) DNS {
	h := ParseHeader(message[:12])
	q := ParseQuestions(message[12:], h.QDCOUNT)
	return DNS{
		header:    &h,
		questions: q,
	}
}

func (d DNS) Writer() []byte {
	message := make([]byte, 0)
	message = append(message, d.header.Writer()...)
	for _, question := range d.questions {
		message = append(message, question.Writer()...)
	}
	for _, answer := range d.answers {
		message = append(message, answer.Writer()...)
	}
	return message
}

func (d DNS) String() string {
	s := ""
	s += fmt.Sprintf("{{header: %v, questions: %v, answers: %v}}", d.Header(), d.questions, d.answers)
	return s
}
