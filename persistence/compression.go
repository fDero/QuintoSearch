package persistence

import (
	"bufio"
	"slices"
)

/*
In v-byte encoding each byte uses the most significant bit (MSB) as a continuation flag:
- If MSB = 1, more bytes follow.
- If MSB = 0, this is the last byte of the integer.
The remaining 7 bits in each byte store the actual value.
*/
func vbyteEncodeUInt64(uint64Value uint64) []byte {
	var encoded []byte
	for uint64Value != 0 {
		// extract and remove 7 bits from the integer
		encodedByte := uint8(uint64Value & 0x7F)
		uint64Value >>= 7
		// if the remaining bits aren't 0 set MSB=1
		if uint64Value != 0 {
			encodedByte |= 0x80
		}

		// prepend the byte to the array of encoded bytes
		encoded = slices.Insert(encoded, 0, encodedByte)
	}

	if len(encoded) == 0 {
		encoded = append(encoded, 0)
	}
	return encoded
}

func vbyteDecodeUInt64(encoded []byte) uint64 {
	var decoded uint64

	// Decoding of v-byte encoded data simply concatenates first 7 bits of each byte
	for _, encodedByte := range encoded {
		decoded <<= 7
		decoded |= uint64(encodedByte & 0x7F)
	}

	return decoded
}

func loadVbyteEncodedUInt64(fileReader *bufio.Reader) ([]byte, error) {
	var encoded []byte
	encodedByte, err := fileReader.ReadByte()
	for ; err != nil; encodedByte, err = fileReader.ReadByte() {
		encoded = append(encoded, encodedByte)
		if encodedByte&0x80 == 0 {
			break
		}
	}
	return encoded, err
}
