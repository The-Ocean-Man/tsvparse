package conv

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

type ConvertionResult struct {
	Values     [][]float32
	Timestamps []int64

	// Optional Metadata
	StartDate time.Time
}

func Convert(file *os.File, ext string) (result ConvertionResult, err error) {
	result = ConvertionResult{}
	err = nil

	r := bufio.NewReader(file)

	switch ext {
	case ".tsv":
		err = tsvConvert(r, &result)
	case ".txt":
		err = txtConvert(r, &result)
	default:
		panic("Unknown file type.")
	}

	// Unnecessery, but I keep it for habits sake
	if err != nil {
		return
	}

	return
}

func readln(r *bufio.Reader) (string, bool, error) {
	str, err := r.ReadString('\n')

	if err == io.EOF {
		return str, true, nil
	}

	if err != nil {
		return "", false, err
	}

	return strings.ReplaceAll(str, "\r", ""), false, nil
}
