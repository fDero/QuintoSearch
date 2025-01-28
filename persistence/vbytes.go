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
	for _, encodedByte := range encoded {
		decoded |= uint64(encodedByte & 0x7F)
		decoded <<= 7
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
