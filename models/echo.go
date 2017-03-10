package models

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

type echoHandler struct {
	Name   string
	Method string
	Path   string
}

func (h echoHandler) String() string {
	return fmt.Sprintf(`e.%s("%s", %s)`, h.Method, h.Path, h.Name)
}

func (h echoHandler) FuncParam() string {
	return `e *echo.Echo`
}

func (h echoHandler) Pkg() string {
	return "github.com/labstack/echo"
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
