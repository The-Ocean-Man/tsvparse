package conv

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"
)

var tsvHeaderStop = "----------------------------------------\n"
var tsvHeaderSep = ":\t"

func tsvConvert(r *bufio.Reader, result *ConvertionResult) error {
	headers, err := tsvParseHeaders(r)
	_ = headers // unused, maybe metadata or something idk

	if err != nil {
		return err
	}

	body, ts, err := tsvParseBody(r)

	if err != nil {
		return err
	}

	result.Values = body
	result.Timestamps = ts

	// for k, v := range headers {
	// 	fmt.Printf("%s=%s", k, v)
	// }

	return nil
}

const TSV_TIME_FMT = "2006-01-02 15:04:05.999999"
const WS = " \r\n\t"

func tsvParseBody(r *bufio.Reader) (body [][]float32, timestamps []int64, err error) {
	err = nil

	// var arr [][]string

	ln, eof, err := readln(r)

	if err != nil {
		return nil, nil, err
	}
	if eof {
		return nil, nil, errors.New("Reached EOF before parsing finished.")
	}

	topsplit := strings.Split(ln, "\t")
	gageCount := len(topsplit) - 5

	body = make([][]float32, 0)
	timestamps = make([]int64, 0) // append every timestamp, no biggie

	readln(r)
	_, eof, err = readln(r)

	if err != nil || eof {
		panic("Could not parse file.")
	}

	// fmt.Println(gageCount)

	// Start parsing body

	// row := 0
	// for { // Parsing rows
	// 	ln, eof, err = readln(r)
	// 	fmt.Println(row)
	// 	// fmt.Print(".")
	// 	if err != nil {
	// 		return
	// 	}
	// 	if eof {
	// 		break
	// 	}
	// 	split := strings.Split(ln, "\t")
	// 	gages := make([]int32, gageCount)
	// 	for i := 0; i < gageCount; i++ { // Parsing collumns
	// 		if i == 0 {
	// 			var ts time.Time
	// 			ts, err = time.Parse(TSV_TIME_FMT, strings.Trim(split[0], WS))
	// 			if err != nil {
	// 				return
	// 			}
	// 			timestamps = append(timestamps, ts.Unix())
	// 		} else if i > 4 {
	// 			var f float64
	// 			f, err = strconv.ParseFloat(split[i], 32)
	// 			if err != nil {
	// 				return
	// 			}
	// 			gages[i-5] = int32(f)
	// 		}
	// 		// fmt.Println(ln[:20])
	// 		body = append(body, gages)
	// 		gages = make([]int32, gageCount)
	// 	}
	// 	row++
	// }

	var arr [][]string

	for {
		ln, eof, err := readln(r)

		if eof {
			break
		}

		if err != nil {
			return nil, nil, err
		}

		arr = append(arr, strings.Split(ln, "\t"))
	}

	// fmt.Println(len(arr), len(arr[0]))

	for _, row := range arr {
		gages := make([]float32, gageCount)
		for i, str := range row {
			if i == 0 {
				var ts time.Time
				ts, err = time.Parse(TSV_TIME_FMT, strings.Trim(str, WS))
				if err != nil {
					return nil, nil, err
				}
				timestamps = append(timestamps, ts.Unix())
			} else if i >= 5 {
				var f float64
				f, err = strconv.ParseFloat(strings.Trim(str, WS), 32)

				if err != nil {
					return nil, nil, err
				}
				gages[i-5] = float32(f)
			}
		}
		body = append(body, gages)
	}

	return
}

func tsvParseHeaders(r *bufio.Reader) (m map[string]string, err error) {
	m = make(map[string]string)
	err = nil
	var ln string
	var eof bool

	ln, eof, err = readln(r)

	if err != nil {
		return
	}

	for ln != tsvHeaderStop {
		if eof {
			break
		}

		split := strings.Split(ln, tsvHeaderSep)

		if len(split) > 3 {
			panic("Malformed header. Too many `:`s.")
		}

		var value string = ""
		if len(split) == 2 {
			value = split[1]
		}
		m[split[0]] = value

		ln, eof, err = readln(r)

		if err != nil {
			return
		}
	}

	return
}
