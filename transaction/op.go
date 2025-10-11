package transaction

type BitcoinOpCode struct{}

func NewBitcoinOpCode() *BitcoinOpCode {
	return &BitcoinOpCode{}
}

func (boc *BitcoinOpCode) DecodeNum(element []byte) int64 {
	bigEndian := reverseByteSlice(element)
	negative := false
	result := int64(0)

	if (bigEndian[0] & 0x80) != 0 {
		negative = true
		// reset the msb to 0
		// 0x7f is 0111 1111
		result = int64(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int64(bigEndian[0])
	}

	for i := 1; i < len(bigEndian); i++ {
		result <<= 8 // 8 bit => 1 byte, prepare space for the next byte
		result += int64(bigEndian[i])
	}

	if negative {
		return -result
	}

	return result
}

func (boc *BitcoinOpCode) EncodeNum(num int64) []byte {
	if num == 0 {
		return []byte("")
	}

	result := []byte{}
	absNum := num
	negative := false

	if num < 0 {
		absNum = -num
		negative = true
	}

	if absNum > 0 {
		// Result will be in little endian format so we append the last byte first (actually we're reversing the byte array)
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	// Check the most significant byte, notice the most significant byte is at the end of the result
	// 0x8080 -> 32896 -32896
	// 0x80 -> 1000 0000
	// 0.... & 1000 0000 -> 0000 0000
	// 1.... & 1000 0000 -> 1000 0000

	if (result[len(result)-1] & 0x80) != 0 {
		if negative {
			// insert 0x80 in the head, most significant byte is currently at the end of result
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else {
		// set the most significant bit to one
		// 1000 0000 || 0... ....
		result[len(result)-1] |= 0x80
	}

	return result
}
