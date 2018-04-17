package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	pngembed "github.com/sabhiram/png-embed"
)

////////////////////////////////////////////////////////////////////////////////

var (
	inputFile  string
	outputFile string
	key        string
	value      string
)

////////////////////////////////////////////////////////////////////////////////

func main() {
	data, err := pngembed.EmbedFile(inputFile, key, value)
	if err == nil {
		err = ioutil.WriteFile(outputFile, data, 777)
	}

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	flag.StringVar(&inputFile, "input", "", "input file name for the png")
	flag.StringVar(&outputFile, "output", "out.png", "output file name for the png")
	flag.StringVar(&key, "key", "TEST_KEY", "key name for the data to inject")
	flag.StringVar(&value, "value", "TEST_VALUE", "sample value to inject for key")

	flag.Parse()
	if len(inputFile) == 0 {
		fmt.Printf("Fatal error: No input file specified!\n")
		os.Exit(1)
	}
}
