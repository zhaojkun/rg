package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/zhaojkun/rg/models"
)

var (
	pkg    = flag.String("pkg", "main", "package name")
	fnName = flag.String("name", "registerRoutes", "function name")
)

func main() {
	flag.Parse()
	gofiles := listFiles(flag.Args())
	var handlers []models.Handler
	for _, file := range gofiles {
		hs := parseFile(file)
		handlers = append(handlers, hs...)
	}
	source := generate(*fnName, handlers)
	formatedSource, err := format.Source(source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(formatedSource))
}

func generate(fnName string, handlers []models.Handler) []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("package %s\n", *pkg))
	if len(handlers) > 0 {
		buf.WriteString(fmt.Sprintf("import \"%s\"\n", handlers[0].Pkg()))
		buf.WriteString(fmt.Sprintf("func %s(%s){\n", fnName, handlers[0].FuncParam()))
		for _, handler := range handlers {
			buf.WriteString(fmt.Sprintf(" %s\n", handler))
		}
		buf.WriteString("}")
	}
	return buf.Bytes()
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
