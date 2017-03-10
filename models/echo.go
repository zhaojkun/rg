package models

import (
	"errors"
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
	if len(fields) != 3 {
		return nil, errors.New("comment not ok")
	}
	name := fields[0]
	method := strings.ToUpper(fields[1])
	path := strings.Trim(fields[2], `" '`)
	return echoHandler{
		name:   name,
		method: method,
		path:   path,
	}, nil
}
