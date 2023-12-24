package drlp

// Buffer is essentially a byte slice. It holds encoded elements.
type Buffer []byte

// Reset resets all content.
func (b *Buffer) Reset() {
	*b = (*b)[:0]
}

// PutUint puts the uint value.
func (b *Buffer) PutUint(i uint64) {
	if i == 0 {
		*b = append(*b, 0x80)
	} else if i < 128 {
		// fits single byte
		*b = append(*b, byte(i))
	} else {
		*b = appendUint(*b, i, 0x80)
	}
}

// PutString puts the string value.
func (b *Buffer) PutString(str []byte) {
	if size := len(str); size == 0 {
		*b = append(*b, 0x80)
	} else if size == 1 && str[0] < 128 {
		// fits single byte, no string header
		*b = append(*b, str[0])
	} else if size < 56 {
		*b = append(*b, 0x80+byte(size))
		*b = append(*b, str...)
	} else {
		*b = appendUint(*b, uint64(size), 0xB7)
		*b = append(*b, str...)
	}
}

// List starts a RLP list.
func (b *Buffer) List() List {
	return List{b, len(*b)}
}

// appendUint appends kind tag and i to b in big endian byte order,
// using the least number of bytes needed to represent i.
func appendUint(b []byte, i uint64, kindTag byte) []byte {
	switch {
	case i < (1 << 8):
		return append(b, kindTag+1, byte(i))
	case i < (1 << 16):
		return append(b, kindTag+2, byte(i>>8), byte(i))
	case i < (1 << 24):
		return append(b, kindTag+3, byte(i>>16), byte(i>>8), byte(i))
	case i < (1 << 32):
		return append(b, kindTag+4, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	case i < (1 << 40):
		return append(b, kindTag+5, byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	case i < (1 << 48):
		return append(b, kindTag+6, byte(i>>40), byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	case i < (1 << 56):
		return append(b, kindTag+7, byte(i>>48), byte(i>>40), byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	default:
		return append(b, kindTag+8, byte(i>>56), byte(i>>48), byte(i>>40), byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	}
}

// uintSize computes the minimum number of bytes required to store i.
func uintSize(i uint64) int {
	switch {
	case i < (1 << 8):
		return 1
	case i < (1 << 16):
		return 2
	case i < (1 << 24):
		return 3
	case i < (1 << 32):
		return 4
	case i < (1 << 40):
		return 5
	case i < (1 << 48):
		return 6
	case i < (1 << 56):
		return 7
	default:
		return 8
	}
}
