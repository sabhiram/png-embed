package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/sabhiram/png-embed"
)

////////////////////////////////////////////////////////////////////////////////

var (
	inputFile  string
	outputFile string
	key        string
	value      string
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	data, err := pngembed.EmbedKeyValue(inputFile, key, value)
	fatalOnError(err)
	fatalOnError(ioutil.WriteFile(outputFile, data, 777))
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(os.Stdout)

	flag.StringVar(&inputFile, "input", "", "input file name for the png")
	flag.StringVar(&outputFile, "output", "out.png", "output file name for the png")
	flag.StringVar(&key, "key", "TEST_KEY", "key name for the data to inject")
	flag.StringVar(&value, "value", "TEST_VALUE", "sample value to inject for key")

	flag.Parse()
	if len(inputFile) == 0 {
		log.Fatalf("No input file specified!\n")
	}
}

////////////////////////////////////////////////////////////////////////////////
