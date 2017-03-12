package models

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

type ginHandler struct {
	name   string
	method string
	path   string
}

func (h ginHandler) Pkg() string {
	return "github.com/gin-gonic/gin"
}

func (h ginHandler) FuncParam() string {
	return `r *gin.Engine`
}

func (h ginHandler) Path() string {
	return h.path
}

func (h ginHandler) String() string {
	return fmt.Sprintf(`r.%s("%s", %s)`, h.method, h.path, h.name)
}

func ginHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	styp, ok := fnDecl.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil, errors.New("type error")
	}
	typ, ok := styp.X.(*ast.SelectorExpr)
	if !ok {
		return nil, errors.New("type error")
	}
	pkgIdent, ok := typ.X.(*ast.Ident)
	if !ok {
		return nil, errors.New("")
	}
	isGinContext := pkgIdent.Name == "gin" && typ.Sel != nil && typ.Sel.Name == "Context"
	if !isGinContext {
		return nil, errors.New("")
	}
	doc := fnDecl.Doc.Text()
	fields := strings.Fields(doc)
	p, err := parseURL(fields)
	if len(fields) != 3 {
		return nil, err
	}
	name := fnDecl.Name.Name
	method := parseHTTPMethod(fields)
	return ginHandler{
		name:   name,
		method: method,
		path:   p,
	}, nil
}
