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
	typ, ok := fnDecl.Type.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return nil, errors.New("type error")
	}
	pkgIdent, ok := typ.X.(*ast.Ident)
	if !ok {
		return nil, errors.New("")
	}
	isEchoContext := pkgIdent.Name == "echo" && typ.Sel != nil && typ.Sel.Name == "Context"
	if !isEchoContext {
		return nil, errors.New("")
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
	}, nil
}
