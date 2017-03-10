package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
)

var (
	filename = flag.String("f", "sample.go", "filename")
)

func main() {
	flag.Parse()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var handlers []Handler
	for _, obj := range f.Scope.Objects {
		h, err := parseHandler(obj)
		if err == nil {
			handlers = append(handlers, h)
		}
	}
	if len(handlers) > 0 {
		fmt.Printf("func registerRoutes(%s){\n", handlers[0].FuncParam())
		for _, handler := range handlers {
			fmt.Println(" ", handler)
		}
		fmt.Println(`}`)
	}
}
