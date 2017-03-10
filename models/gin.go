package models

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

type ginHandler struct {
	Name   string
	Method string
	Path   string
}

func (h ginHandler) String() string {
	return fmt.Sprintf(`r.%s("%s", %s)`, h.Method, h.Path, h.Name)
}

func (h ginHandler) FuncParam() string {
	return `r *gin.Engine`
}

func (h ginHandler) Pkg() string {
	return "github.com/gin-gonic/gin"
}

func ginHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	doc := fnDecl.Doc.Text()
	fields := strings.Fields(doc)
	if len(fields) != 3 {
		return nil, errors.New("comment not ok")
	}
	name := fields[0]
	method := strings.ToUpper(fields[1])
	path := strings.Trim(fields[2], `" '`)
	return ginHandler{
		Name:   name,
		Method: method,
		Path:   path,
	}, nil
}
