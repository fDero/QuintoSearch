package persistence

import (
	"bufio"
)

func vbyteEncodeUInt64(uint64Value uint64) []byte {
	var encoded []byte
	for uint64Value != 0 {
		encodedByte := uint8(uint64Value & 0x7F)
		uint64Value >>= 7
		if uint64Value != 0 {
			encodedByte |= 0x80
		}
		encoded = append(encoded, encodedByte)
	}
	if len(encoded) == 0 {
		encoded = append(encoded, 0)
	}
	return encoded
}

func vbyteDecodeUInt64(encoded []byte) uint64 {
	var decoded uint64
	for i := len(encoded) - 1; i >= 0; i-- {
		decoded <<= 7
		encodedByte := encoded[i]
		decoded |= uint64(encodedByte & 0x7F)
	}
	return decoded
}

func loadVbyteEncodedUInt64(fileReader bufio.Reader) ([]byte, error) {
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
