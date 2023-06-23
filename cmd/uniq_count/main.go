package main

import (
	"flag"
	"fmt"
	"k-way-merge.git/sort"
	"k-way-merge.git/uniq"
	"log"
	"os"
)

var (
	input   = flag.String("input", "input.txt", "")
	maxRows = flag.Int("N", 100, "")
)

func main() {
	flag.Parse()
	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	if *maxRows <= 1 {
		log.Fatal("bad maxRows")
	}
	sortedFile, err := os.Create(fmt.Sprintf("%s_sorted", *input))
	if err != nil {
		log.Fatal(err)
	}

	err = sort.Sort(file, *maxRows, sortedFile)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.Create(fmt.Sprintf("%s_uniq", *input))
	if err != nil {
		log.Fatal(err)
	}
	err = uniq.Uniq(sortedFile, output)
	if err != nil {
		log.Fatal(err)
	}
	err = output.Close()
	if err != nil {
		log.Fatal(err)
	}
}
