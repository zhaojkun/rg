package models

import (
	"fmt"
	"go/ast"
	"strings"
)

type echoHandler struct {
	name   string
	method string
	path   string
}

func (h echoHandler) Pkg() string {
	return "github.com/labstack/echo"
}

func (h echoHandler) FuncParam() string {
	return `e *echo.Echo`
}

func (h echoHandler) Path() string {
	return h.path
}

func (h echoHandler) String() string {
	return fmt.Sprintf(`e.%s("%s", %s)`, h.method, h.path, h.name)
}

func echoHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	doc := fnDecl.Doc.Text()
	fields := strings.Fields(doc)
	p, err := parseURL(fields)
	if err != nil {
		return nil, err
	}
	name := fnDecl.Name.Name
	method := parseHTTPMethod(fields)
	return echoHandler{
		name:   name,
		method: method,
		path:   p,
	}, nil
}
