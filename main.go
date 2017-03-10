package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/zhaojkun/rg/models"
)

func main() {
	flag.Parse()
	gofiles := listFiles(flag.Args())
	var handlers []models.Handler
	for _, file := range gofiles {
		hs := parseFile(file)
		handlers = append(handlers, hs...)
	}
	if len(handlers) > 0 {
		fmt.Printf("func registerRoutes(%s){\n", handlers[0].FuncParam())
		for _, handler := range handlers {
			fmt.Println(" ", handler)
		}
		fmt.Println(`}`)
	}
}

func listFiles(args []string) []string {
	var gofiles []string
	for _, filename := range args {
		info, err := os.Lstat(filename)
		if err != nil {
			continue
		}
		if info.IsDir() {
			files, err := ioutil.ReadDir(filename)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				gofiles = append(gofiles, path.Join(filename, file.Name()))
			}
		} else {
			gofiles = append(gofiles, filename)
		}
	}
	return gofiles
}

func parseFile(filename string) []models.Handler {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}
	var imports []string
	for _, s := range f.Imports {
		imports = append(imports, strings.Trim(s.Path.Value, "\""))
	}
	builder, ok := models.GetBuilder(imports)
	if !ok {
		return nil
	}
	var handlers []models.Handler
	for _, obj := range f.Scope.Objects {
		if obj.Kind == ast.Fun {
			fnDecl, ok := (obj.Decl).(*ast.FuncDecl)
			if !ok {
				continue
			}
			if h, err := builder(fnDecl); err == nil {
				handlers = append(handlers, h)
			}
		}
	}
	return handlers
}
