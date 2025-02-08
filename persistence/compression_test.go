package persistence

import (
	"testing"

	// "crypto/rand"
	"bytes"
	"encoding/binary"
	"os"
)

func TestVbyteEncodeUInt64(t *testing.T) {
	expected := map[uint64][]byte{
		// some edge cases
		0:          {0},
		0x80:       {1, 0x80},
		0x400:      {8, 0x80},
		0x80000000: {8, 0x80, 0x80, 0x80, 0x80},
		0xFF:       {1, 0xFF},
		0xFFFF:     {3, 0xFF, 0xFF},
		0xFFFFFF:   {7, 0xFF, 0xFF, 0xFF},
		0xFFFFFFFF: {15, 0xFF, 0xFF, 0xFF, 0xFF},
		// some random cases
		17726549421771500413: {1, 246, 128, 214, 216, 250, 144, 226, 214, 253},
		226058532925797298:   {3, 145, 199, 225, 155, 203, 243, 231, 178},
		14236130522442077043: {1, 197, 200, 185, 169, 138, 242, 134, 166, 243},
		16698193177512444535: {1, 231, 221, 249, 199, 150, 131, 236, 132, 247},
		9224194194998898169:  {1, 128, 129, 186, 247, 249, 225, 132, 219, 249},
		11204595683350715874: {1, 155, 191, 174, 141, 141, 175, 250, 187, 226},
		6368678729152153642:  {88, 177, 133, 241, 251, 221, 205, 144, 170},
		10954371072751130213: {1, 152, 130, 239, 206, 187, 137, 149, 172, 229},
		10277747424137571551: {1, 142, 208, 249, 154, 193, 151, 159, 217, 223},
		10158608212689706916: {1, 140, 253, 168, 139, 140, 205, 241, 191, 164},
		10154246904502481127: {1, 140, 245, 200, 184, 169, 195, 235, 225, 231},
		13490447729012836275: {1, 187, 155, 236, 215, 190, 230, 164, 143, 179},
		16383005548227477251: {1, 227, 174, 136, 151, 236, 186, 232, 254, 131},
		12399423697845790562: {1, 172, 137, 230, 197, 136, 141, 162, 206, 226},
		251263606851712567:   {3, 190, 170, 221, 180, 251, 132, 236, 183},
		18370130104345724863: {1, 254, 247, 244, 128, 131, 153, 242, 143, 191},
	}

	for key, value := range expected {
		en := vbyteEncodeUInt64(key)
		t.Logf("%v", en)
		if !bytes.Equal(value, vbyteEncodeUInt64(key)) {
			t.Errorf("Expected %v, got %v", value, en)
		}
	}
}

func TestVbyteDecodeUInt64(t *testing.T) {
	file, err := os.Open("../test_data/random_data.bin")
	if err != nil {
		t.Errorf("Test run is incomplete. Cannot open sample testing data file: %s", err)
	}
	defer file.Close()

	for {
		var num uint64
		err := binary.Read(file, binary.LittleEndian, &num)
		if err != nil {
			break
		}
		t.Logf("Test Case for n = %d \t(%064b)", num, num)
		en := vbyteEncodeUInt64(num)
		de := vbyteDecodeUInt64(en)
		if de != num {
			t.Errorf("Error not encoded or decoded correctly\nOriginal: %064b\nEncoded: %08b\nDecoded: %064b", num, en, de)
		}
	}
}
