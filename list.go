package drlp

// List presents a list object.
type List struct {
	buf    *Buffer
	offset int
}

// Discard discards the list. It requires the list not ended or just ended.
func (l *List) Discard() {
	*l.buf = (*l.buf)[:l.offset]
}

// End ends the list and returns bytes of encoded list.
func (l *List) End() []byte {
	buf := *l.buf
	offset := l.offset
	contentSize := len(buf) - offset
	if contentSize < 56 {
		// shift the content to make room for list header
		buf = append(buf[:offset+1], buf[offset:]...)
		// write list header
		buf[offset] = 0xC0 + byte(contentSize)
	} else {
		headerSize := uintSize(uint64(contentSize)) + 1
		// shift the content to make room for list header
		buf = append(buf[:offset+headerSize], buf[offset:]...)
		// write list header
		appendUint(buf[:offset], uint64(contentSize), 0xF7)
	}
	*l.buf = buf
	return buf[offset:]
}
