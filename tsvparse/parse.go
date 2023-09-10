package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"unsafe"
)

var ProgressCallback func(progress float32, done bool)

type BinaryResult struct {
	RowCount, CollCount int

	Timestamps []int64
	Values     [][]float32
}

var nativeEndian binary.ByteOrder

var size = 0

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

func readInt32(r *bufio.Reader) int32 {
	var i int32 = 0
	binary.Read(r, nativeEndian, &i)
	size += 4
	return i
}

func readInt64(r *bufio.Reader) int64 {
	var i int64 = 0
	binary.Read(r, nativeEndian, &i)
	size += 8
	return i
}

func readFloat32(r *bufio.Reader) float32 {
	var f float32 = 0
	if err := binary.Read(r, nativeEndian, &f); err != nil {
		// fmt.Println(err)
	}
	size += 4
	return f
}

func parseBinary(r *bufio.Reader) (res BinaryResult, err error) {
	// start := now()
	err = nil

	setEndian()

	rowCount := readInt32(r)
	collCount := readInt32(r)
	fmt.Println(rowCount, collCount)
	res = BinaryResult{int(rowCount), int(collCount), make([]int64, rowCount), make([][]float32, rowCount)}

	for i := 0; i < int(rowCount); i++ {
		res.Timestamps[i] = readInt64(r)
		doCallback(float32(i)/float32(rowCount), false)

		row := make([]float32, collCount)

		for j := 0; j < int(collCount); j++ {
			row[j] = readFloat32(r)
		}

		res.Values[i] = row
	}

	// res = BinaryResult{0, 0, make([]int64, 0), make([][]float32, 0)}
	// i := 0

	// for {
	// 	ts := readInt64(r)

	// 	if ts == math.MaxInt64 {
	// 		fmt.Println("Breaking")
	// 		break
	// 	}

	// 	values := make([]float32, 0)

	// 	for {
	// 		v := readFloat32(r)

	// 		if v == math.MaxFloat32 {
	// 			fmt.Println("small breaking")
	// 			break
	// 		}

	// 		values = append(values, v)

	// 		// doCallback(float32(ii), false)
	// 	}
	// 	i++
	// 	fmt.Println(i)

	// 	res.Timestamps = append(res.Timestamps, ts)
	// 	res.Values = append(res.Values, values)
	// }

	// res.RowCount = len(res.Values)
	// res.CollCount = len(res.Values[0])

	doCallback(1, true)

	// fmt.Println("Size is now:", size, "Tru?:", size == 8479880)

	return
}

func doCallback(f float32, done bool) {
	if ProgressCallback != nil {
		ProgressCallback(f, done)
	}
}
