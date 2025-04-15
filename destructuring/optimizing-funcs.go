package destructuring

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"pfanalyzer2/utils"
	"unicode/utf16"
)

//here i'll put function that are used various time to re-utilize code

// FormatTime formatta una data nel formato "2025-03-13 20:45:08"
func PrintTitle(title string) {
	fmt.Printf("\n\033[91m%s:\033[0m\n", title)
}

func ReadData(reader *bytes.Reader, field any) error {
	return binary.Read(reader, binary.LittleEndian, field)
}

func uint64ToStringDate(decimalFT uint64) string {
	return utils.IntFTToStringDate(int64(decimalFT))
}

func decodeUTF16(utf16Bytes []byte) string {
	utf16Decoded := []uint16{}
	for i := 0; i < len(utf16Bytes); i += 2 {
		codeUnit := uint16(utf16Bytes[i]) | uint16(utf16Bytes[i+1])<<8
		if codeUnit == 0 {
			break
		} //if eos is found break
		utf16Decoded = append(utf16Decoded, codeUnit)
	}
	return string(utf16.Decode(utf16Decoded))
}
