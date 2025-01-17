package main

type DNS struct {
	header *Header
}

func New(id uint16) DNS {
	return DNS{
		header: &Header{ID: id},
	}
}

func (d DNS) Header() *Header {
	return d.header
}

func Parse(message []byte) DNS {
	h := ParseHeader(message[:12])
	return DNS{
		header: &h,
	}
}

func (d DNS) Writer() []byte {
	message := make([]byte, 0)
	message = append(message, d.header.Writer()...)
	return message
}
