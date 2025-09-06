package main

import (
	"flag"
	"fmt"
	"os"

	template "github.com/fuwjax/gopase/happy/sample"
	"github.com/fuwjax/gopase/parser"
	json "github.com/fuwjax/gopase/parser/sample"
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
	opts := panicUnless(json.ParseJson(string(config))).(map[string]any)
	rules := panicUnless(parser.Bootstrap(string(grammar)))
	result := panicUnless(template.RenderPeg(rules, opts))
	fmt.Print(result)
}
