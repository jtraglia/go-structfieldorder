package main

import (
	"flag"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/jtraglia/go-structfieldorder/analyzer"
)

func main() {
	flag.Bool("unsafeptr", false, "")

	a, err := analyzer.NewAnalyzer(nil, nil)
	if err != nil {
		panic(err)
	}

	singlechecker.Main(a)
}
