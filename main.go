package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/sample"
)

func panicUnless[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.json", "path to the json config file")
	flag.Parse()
	grammarPath := flag.Arg(0)
	config := panicUnless(os.ReadFile(configPath))
	grammar := panicUnless(os.ReadFile(grammarPath))
	opts := panicUnless(sample.ParseJson(string(config))).(map[string]any)
	rules := panicUnless(parser.Bootstrap(string(grammar)))
	result := panicUnless(sample.RenderPeg(rules, opts))
	fmt.Print(result)
}
