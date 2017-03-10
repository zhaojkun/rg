package main

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

type Handler interface {
	FuncParam() string
	String() string
}
type HandlerBuilder func(*ast.FuncDecl) (Handler, error)

var (
	frameworkHandler map[string]HandlerBuilder
)

func init() {
	frameworkHandler = make(map[string]HandlerBuilder)
	Register("echo", echoHandlerBuilder)
}

func Register(path string, h HandlerBuilder) {
	frameworkHandler[path] = h
}

func parseHandler(obj *ast.Object) (Handler, error) {
	if obj.Kind == ast.Fun {
		fnDecl, ok := (obj.Decl).(*ast.FuncDecl)
		if !ok {
			return nil, errors.New("not func")
		}
		for _, builder := range frameworkHandler {
			if h, err := builder(fnDecl); err == nil {
				return h, nil
			}
		}
	}
	return nil, errors.New("not func")
}

type echoHandler struct {
	Name   string
	Method string
	Path   string
}

func echoHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	doc := fnDecl.Doc.Text()
	fields := strings.Fields(doc)
	if len(fields) != 3 {
		return nil, errors.New("comment not ok")
	}
	name := fields[0]
	method := strings.ToUpper(fields[1])
	path := strings.Trim(fields[2], `" '`)
	return echoHandler{
		Name:   name,
		Method: method,
		Path:   path,
	}, nil

}

func (h echoHandler) String() string {
	return fmt.Sprintf(`e.%s("%s", %s)`, h.Method, h.Path, h.Name)
}

func (h echoHandler) FuncParam() string {
	return `e *echo.Echo`
}
