package conv

import (
	"bufio"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func txtConvert(r *bufio.Reader, result *ConvertionResult) error {
	// Read filename
	filename, eof, err := readln(r)
	_ = filename
	_, eof, err = readln(r)
	if eof || err != nil {
		if err != nil {
			log.Fatalln("File has no body:", err)
		}
		log.Fatalln("File has no body")
	}

	if eof {
		panic("Reached unexpected EOF in filename.")
	}
	if err != nil {
		return err
	}

	firstHeader := make([]string, 0)
	// First header
	for {
		str, eof, err := readln(r)

		if eof {
			panic("Reached unexpected EOF in first header.")
		}
		if err != nil {
			return err
		}

		if str == "\n" {
			break
		}

		firstHeader = append(firstHeader, str)
	}

	secondHeader := make([]string, 0)
	// Second header
	for {
		str, eof, err := readln(r)

		if eof {
			panic("Reached unexpected EOF in first header.")
		}
		if err != nil {
			return err
		}

		if str == "\n" {
			break
		}

		secondHeader = append(secondHeader, str)
	}

	var tsStr string

	for _, s := range secondHeader {
		if strings.HasPrefix(s, "T0 =") {
			tsStr = s[:len(s)-1]
			break
		}
	}

	// fmt.Printf("tsStr: %v\n", tsStr)
	timestamps, err := txtParseTimestamps(tsStr)

	if err != nil {
		return err
	}

	// Read bodyStr
	bodyStr := make([][]string, 0)

	for {
		ln, eof, err := readln(r)

		if eof {
			break
		}
		if err != nil {
			return err
		}

		ln = strings.ReplaceAll(ln, "\n", "")
		split := strings.Split(ln, "\t")

		for i, s := range split {
			split[i] = strings.Trim(s, " ")
		}

		bodyStr = append(bodyStr, split)
	}

	// Parse body
	body, ts, err := txtParseBody(bodyStr, timestamps[0])

	if err != nil {
		return err
	}

	result.Timestamps = ts
	result.Values = body

	return nil
}

func txtParseBody(body [][]string, t0 int64) (arr [][]float32, ts []int64, err error) {
	arr = make([][]float32, len(body))
	ts = make([]int64, len(body))
	err = nil

	for i, row := range body {
		tmp := make([]float32, 0)
		for j, value := range row {
			if j == 0 {
				f, err := txtParseTsNumber(value)
				t := int64(f)

				if err != nil {
					return nil, nil, err
				}

				ts[i] = t
				continue
			}

			gage, err := txtParseNumber(value)

			if err != nil {
				return nil, nil, err
			}
			if gage != -1000000 {
				tmp = append(tmp, gage)
			} else {
			}
		}
		arr[i] = tmp
		// fmt.Print(len(tmp), " ")
	}

	// fmt.Println(ts)
	return
}

func txtParseTsNumber(str string) (float64, error) {
	str = strings.Trim(str, " ")
	str = strings.ReplaceAll(str, ",", ".")

	return strconv.ParseFloat(str, 64)
}

func txtParseTimestamps(timestamps string) ([]int64, error) {
	split := strings.Split(timestamps, "\t")

	ts := make([]int64, len(split))

	for i, str := range split {
		str = str[4:]
		t, err := time.Parse("2006-01-02 15:04:05", strings.Trim(str, " "))

		if err != nil {
			t, err = time.Parse("06-01-02 15:04:05", str)

			if err != nil {
				return nil, err
			}
		}

		ts[i] = t.Unix()
	}

	return ts, nil
}

func txtParseNumber(str string) (float32, error) {
	if len(str) == 0 {
		return 0, nil
	}
	str = strings.Trim(str, " \t")
	str = strings.ReplaceAll(str, ",", ".")

	// No need to handle nan, parseFloat will return zero

	if strings.Contains(str, "E") {
		// Exponent fuuuuck
		split := strings.Split(str, "E")

		if len(split) != 2 {
			panic("Malformed number in gage value.")
		}

		return parseExponent(split[0], split[1])
	}
	f, err := strconv.ParseFloat(str, 32)
	return float32(f), err
}

var iii int

func parseExponent(sbase, sexp string) (float32, error) {
	exp, err := strconv.ParseInt(sexp, 10, 32)

	if err != nil {
		return 0, err
	}

	base, err := strconv.ParseFloat(sbase, 32)

	if err != nil {
		return 0, err
	}

	return float32(base * math.Pow10(int(exp))), nil
}
