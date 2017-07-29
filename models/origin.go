package models

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

type originHandler struct {
	name string
	path string
}

func (h originHandler) ToHTTP() HTTPHandler {
	return HTTPHandler{
		Name: h.name,
		Path: h.path,
	}
}
func (h originHandler) Method() (string, string, bool) {
	return "", "", false
}

func (h originHandler) Pkg() string {
	return "net/http"
}

func (h originHandler) FuncParam() string {
	return ""
}

func (h originHandler) Path() string {
	return h.path
}

func (h originHandler) String() string {
	return fmt.Sprintf(`http.HandleFunc("%s", %s)`, h.path, h.name)
}
func (h originHandler) JS() string {
	return ""
}

func (h originHandler) Doc() string {
	return ""
}

func originHandlerBuilder(fnDecl *ast.FuncDecl) (Handler, error) {
	if !checkFuncInterface(fnDecl, "http.ResponseWriter", "*http.Request") {
		return nil, errors.New("type error")
	}
	doc := fnDecl.Doc.Text()
	fields := strings.Fields(doc)
	p, err := parseURL(fields)
	if err != nil {
		return nil, err
	}
	name := fnDecl.Name.Name
	return originHandler{
		name: name,
		path: p,
	}, nil
}
