package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tsvconv/conv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Expected an input file as an only argument.")
	}
	filename := os.Args[1]
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalln("An error occured", err)
	}

	defer file.Close()

	// Handle file
	result, err := conv.Convert(file, filepath.Ext(filename))

	if err != nil {
		panic(err)
	}

	filename = strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + ".plotbin"

	f, _ := os.Create(filename)
	defer f.Close()

	err = conv.WriteToFile(f, result.Values, result.Timestamps)

	if err != nil {
		fmt.Println(err)
	}
}
