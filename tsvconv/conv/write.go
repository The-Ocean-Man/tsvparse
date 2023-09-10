package conv

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

var nativeEndian binary.ByteOrder

func setEndian() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

func WriteToFile(file *os.File, data [][]float32, unixTimestamps []int64) error {
	setEndian()

	w := bufio.NewWriterSize(file, 4096)

	binary.Write(w, nativeEndian, int64(0)) // Placeholder for sizes

	lastUnix := int64(0)

	rowCT, collCT := 0, 0

	fmt.Println("lengths:", len(data), len(unixTimestamps), len(data[0100]))

	for i, row := range data {

		// Filter out duplicate seconds
		if unixTimestamps[i] == lastUnix {
			continue
		}
		rowCT++
		lastUnix = unixTimestamps[i]

		if lastUnix == 0 {
			fmt.Println("Ts is 0")
		}

		binary.Write(w, nativeEndian, unixTimestamps[i]) // Unix timestamp for row, 8
		for _, f := range row {
			binary.Write(w, nativeEndian, f) // Float gage value, 4
			if i == 0 {
				collCT++
			}
		}
	}

	w.Flush()

	file.Seek(0, io.SeekStart)

	binary.Write(file, nativeEndian, int32(rowCT))
	binary.Write(file, nativeEndian, int32(collCT))

	w.Flush()

	// fmt.Println("Ugh", rowCT, collCT)

	return nil
}
