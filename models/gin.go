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

func (h ginHandler) ToHTTP() HTTPHandler {
	return HTTPHandler{
		Name:   h.name,
		Method: h.method,
		Path:   h.path,
	}
}
func (h ginHandler) Method() (string, string, bool) {
	return "", "", false
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

func (h ginHandler) JS() string {
	return ""
}
func (h ginHandler) Doc() string {
	return ""
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
