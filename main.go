package main

import (
	"bytes"
	"encoding/json"
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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/zhaojkun/rg/models"
)

var (
	pkg          = flag.String("pkg", "main", "package name")
	fnName       = flag.String("name", "registerRoutes", "function name")
	Js           = flag.Bool("js", false, "for js")
	postman      = flag.Bool("postman", false, "for postman")
	fileIgnoreRe = regexp.MustCompile(`(?P<build>\+build)\s+(?P<tool>!rg)`)
)

func main() {
	flag.Parse()
	gofiles := listFiles(flag.Args())
	var handlers []models.Handler
	for _, file := range gofiles {
		hs := parseFile(file)
		handlers = append(handlers, hs...)
	}
	sort.Sort(models.HandlerSorter(handlers))
	if *Js {
		source := generateJS(handlers)
		fmt.Println(string(source))
		return
	}
	if *postman {
		source := generatePostman(flag.Args()[0], handlers)
		fmt.Println(string(source))
		return
	}
	source := generate(*fnName, handlers)
	formatedSource, err := format.Source(source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(formatedSource))
}

func generateJS(handlers []models.Handler) []byte {
	var buf bytes.Buffer
	for _, handler := range handlers {
		buf.WriteString(handler.Doc())
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

type postmanRequest struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type postmanItem struct {
	Name    string         `json:"name"`
	Request postmanRequest `json:"request"`
}

type postMan struct {
	Variables []string          `json:"variables"`
	Info      map[string]string `json:"info"`
	Item      []postmanItem     `json:"item"`
}

func generatePostman(name string, handlers []models.Handler) []byte {
	var items []postmanItem
	for _, handler := range handlers {
		h := handler.ToHTTP()
		item := postmanItem{
			Name: h.Name,
			Request: postmanRequest{
				URL:    "{{baseurl}}" + h.Path,
				Method: h.Method,
			},
		}
		items = append(items, item)
	}
	name = fmt.Sprintf("%s_%v", name, time.Now().Format("2006-01-02_15:04:05"))
	pm := postMan{
		Variables: make([]string, 0, 0),
		Info: map[string]string{
			"name":        name,
			"description": "",
			"schema":      "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		},
		Item: items,
	}
	buf, _ := json.Marshal(pm)
	return buf
}

func generate(fnName string, handlers []models.Handler) []byte {
	var buf bytes.Buffer
	buf.WriteString("// AUTOGENERATED FILE BY github.com/zhaojkun/rg\n\n")
	buf.WriteString(fmt.Sprintf("package %s\n", *pkg))
	if len(handlers) > 0 {
		pkg := handlers[0].Pkg()
		buf.WriteString(fmt.Sprintf("import \"%s\"\n", pkg))
		name, cls, ok := handlers[0].Method()
		if ok {
			buf.WriteString(fmt.Sprintf("func (%s %s)%s(%s){\n", name, cls, fnName, handlers[0].FuncParam()))
		} else {
			buf.WriteString(fmt.Sprintf("func %s(%s){\n", fnName, handlers[0].FuncParam()))
		}
		for _, handler := range handlers {
			if pkg != handler.Pkg() {
				continue
			}
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
	for i := 0; i < len(f.Comments); i++ {
		if fileIgnoreRe.MatchString(f.Comments[i].Text()) {
			return nil
		}
	}
	var imports []string
	for _, s := range f.Imports {
		imports = append(imports, strings.Trim(s.Path.Value, "\""))
	}
	builder, ok := models.GetBuilder(imports)
	if !ok {
		return nil
	}
	//ast.Print(fset, f)
	var handlers []models.Handler
	for _, obj := range f.Decls {
		fnDecl, ok := (obj).(*ast.FuncDecl)
		if !ok {
			continue
		}
		if h, err := builder(fnDecl); err == nil {
			handlers = append(handlers, h)
		}
	}
	return handlers
}
