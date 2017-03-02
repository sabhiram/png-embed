package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/sabhiram/png-embed"
)

var (
	inputFile  string
	outputFile string
	key        string
)

func handleErr(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}

func main() {
	flag.Parse()

	if len(key) == 0 {
		log.Fatalf("No key specified!\n")
	}
	if len(inputFile) == 0 {
		log.Fatalf("No input file specified!\n")
	}
	if len(outputFile) == 0 {
		log.Fatalf("No output file specified!\n")
	}

	s := struct {
		Foo string `json:"Foo"`
		Bar string `json:"Bar"`
	}{
		Foo: "FooValue",
		Bar: "BarValue",
	}

	data, err := pngembed.EmbedMap(inputFile, key, s)
	handleErr(err)

	handleErr(ioutil.WriteFile(outputFile, data, 777))
}

func init() {
	flag.StringVar(&inputFile, "input", "", "input file name for the png")
	flag.StringVar(&inputFile, "i", "", "input file name for the png")

	flag.StringVar(&outputFile, "output", "", "output file name for the png")
	flag.StringVar(&outputFile, "o", "", "output file name for the png")

	flag.StringVar(&key, "key", "", "key name for the data to inject")
	flag.StringVar(&key, "k", "", "key name for the data to inject")
}
