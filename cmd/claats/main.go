package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ShawnLabo/claats"
)

func main() {
	var in, out, pattern, ga, title string

	flag.StringVar(&in, "in", "./", "Input directory")
	flag.StringVar(&out, "out", "docs", "Output directory")
	flag.StringVar(&pattern, "pattern", "**/*.md", "Source file pattern")
	flag.StringVar(&ga, "ga", "", "Google Analytics Tracking ID (UA-XXXX-X")
	flag.StringVar(&title, "title", "claats site", "Site name")

	flag.Parse()

	opt := claats.Option{
		In:      in,
		Out:     out,
		Pattern: pattern,
		GA:      ga,
		Title:   title,
	}

	if err := claats.Generate(opt); err != nil {
		fmt.Fprintf(os.Stderr, "error\n%v\n", err)
		os.Exit(1)
	}
}
