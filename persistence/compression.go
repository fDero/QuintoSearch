package persistence

import (
	"bufio"
	"slices"
)

func vbyteEncodeUInt64(uint64Value uint64) []byte {
	var encoded []byte
	for uint64Value != 0 {
		encodedByte := withMSBtoZero(uint8(uint64Value))
		uint64Value >>= 7
		if uint64Value != 0 {
			encodedByte = withMSBtoOne(encodedByte)
		}
		encoded = prependToSlice(encoded, encodedByte)
	}
	if len(encoded) == 0 {
		encoded = append(encoded, 0)
	}
	return encoded
}

func vbyteDecodeUInt64(encoded []byte) uint64 {
	var decoded uint64
	for _, encodedByte := range encoded {
		decoded <<= 7
		decoded |= uint64(withMSBtoZero(encodedByte))
	}
	return decoded
}

func loadVbyteEncodedUInt64(fileReader *bufio.Reader) ([]byte, error) {
	var encoded []byte
	encodedByte, err := fileReader.ReadByte()
	for ; err == nil; encodedByte, err = fileReader.ReadByte() {
		encoded = append(encoded, encodedByte)
		if encodedByte&0x80 == 0 {
			break
		}
	}
	return encoded, err
}

func withMSBtoZero(x uint8) uint8 {
	return x & 0b01111111
}

func withMSBtoOne(x uint8) uint8 {
	return x | 0b10000000
}

func prependToSlice(slice []byte, value uint8) []byte {
	return slices.Insert(slice, 0, value)
}
