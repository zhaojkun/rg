package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

var (
	filename = flag.String("f", "sample.go", "filename")
)

type EchoPath struct {
	Name   string
	Method string
	Path   string
}

func (e EchoPath) String() string {
	return fmt.Sprintf(`e.%s("%s", %s)`, e.Method, e.Path, e.Name)
}

func parsePath(doc string) EchoPath {
	fields := strings.Fields(doc)
	name := fields[0]
	method := strings.ToUpper(fields[1])
	path := strings.Trim(fields[2], `" '`)
	return EchoPath{
		Name:   name,
		Method: method,
		Path:   path,
	}
}

func main() {
	flag.Parse()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var handlers []EchoPath
	for _, obj := range f.Scope.Objects {
		if obj.Kind == ast.Fun {
			fnDecl, ok := (obj.Decl).(*ast.FuncDecl)
			if !ok {
				continue
			}
			handlers = append(handlers, parsePath(fnDecl.Doc.Text()))
		}
	}
	fmt.Println(`func registerRoutes(e *echo.Echo){`)
	for _, handler := range handlers {
		fmt.Println("       ", handler)
	}
	fmt.Println(`}`)
}
