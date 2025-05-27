/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

V-byte encoding is a variable-length encoding scheme for unsigned integers. It is used
to compress integers that are likely to be small. Each byte uses the most significant
bit (MSB) as a continuation flag (The remaining 7 bits store the actual value):
	- If MSB = 1, more bytes follow.
	- If MSB = 0, this is the last byte of the integer.

For example, the following numbers are encoded as follows:
  00000000 00000000 00000000 00110101 -> [0]0110101
  00000000 00000000 00000011 10110101 -> [0]0000111 [1]0110101
  00000000 00000000 00100011 10110101 -> [0]1000111 [1]0110101
  00000000 00000000 11100011 10110101 -> [0]0000011 [1]1000111 [1]0110101

This encoding uses fewer bytes for small numbers, which is useful for compressing
the difference between document IDs and positions in an inverted index. Since
document IDs and positions are stored sequentially, the difference between them
is likely to be small, and thus the v-byte encoding will be efficient.
==================================================================================*/

package persistence

import (
	"io"
	"slices"
)

func vbyteEncodeUInt64[Num64Bit ~uint64](uint64Value Num64Bit) []byte {
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

func loadVbyteEncodedUInt64(fileReader io.ByteReader) ([]byte, error) {
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
