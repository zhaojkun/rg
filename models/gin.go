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
	if !checkFuncInterface(fnDecl, "*gin.Context") {
		return nil, errors.New("type error")
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
