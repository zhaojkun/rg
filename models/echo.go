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
	fn     *ast.FuncDecl
}

func (h echoHandler) Method() (string, string, bool) {
	if h.fn.Recv != nil {
		name := h.fn.Recv.List[0].Names[0].Name
		cls := h.fn.Recv.List[0].Names[0].Obj.Decl.(*ast.Field).Type.(*ast.StarExpr).X.(*ast.Ident).Name
		return name, "*" + cls, true
	}
	return "", "", false
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
	if h.fn.Recv != nil {
		name := h.fn.Recv.List[0].Names[0].Name
		return fmt.Sprintf(`e.%s("%s", %s.%s)`, h.method, h.path, name, h.name)
	}
	return fmt.Sprintf(`e.%s("%s", %s)`, h.method, h.path, h.name)
}

func echoHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	if !checkFuncInterface(fnDecl, "echo.Context") {
		return nil, errors.New("type error")
	}
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
		fn:     fnDecl,
	}, nil
}
