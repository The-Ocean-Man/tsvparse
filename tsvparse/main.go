package main

// TODO:
// Some .txt files (primarily haparandas) have 0 coll length

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH = 1, 1, 0

var VERSION = fmt.Sprint(VERSION_MAJOR) + "." + fmt.Sprint(VERSION_MINOR) + "." + fmt.Sprint(VERSION_PATCH)

func main() {
	if len(os.Args) == 1 {
		StartGUI()
	} else {
		useCMD()
	}
}

func useCMD() {
	f, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatalln(err)
	}

	result, err := parseBinary(bufio.NewReader(f))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result.Timestamps[:10])
	fmt.Println(len(result.Values))
	fmt.Println(len(result.Values[0]))
}
