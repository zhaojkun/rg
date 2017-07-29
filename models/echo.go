package models

import (
	"bytes"
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

func (h echoHandler) ToHTTP() HTTPHandler {
	return HTTPHandler{
		Name:   h.name,
		Method: h.method,
		Path:   h.path,
	}
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

func (h echoHandler) Doc() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "@api {%s} %s\n", h.method, h.path)
	fmt.Fprintf(&buf, "@apiName %s", h.name)
	return buf.String()
}

func (h echoHandler) JS() string {
	arr := strings.Split(h.path, "/")
	var res string
	var tarr []string
	var params []string
	for i := 0; i < len(arr); i++ {
		if len(arr[i]) > 0 && arr[i][0] == ':' {
			param := arr[i][1:]
			tarr = append(tarr, `"`+res+`"`, param)
			params = append(params, param)
			res = ""
		} else {
			res += arr[i]
		}
		if i < len(arr)-1 {
			res += "/"
		}
	}
	if len(res) > 0 {
		tarr = append(tarr, `"`+res+`"`)
	}
	u := strings.Join(tarr, "+")
	requestBody := "{}"
	if strings.ToLower(h.method) == "delete" {
		requestBody = "{body:{}}"
	}
	fnHeader := fmt.Sprintf("%s(%s){\n", h.name, strings.Join(params, ","))
	fnBody := fmt.Sprintf("    that.%s(%s,%s)", strings.ToLower(h.method), u, requestBody)
	return fnHeader + fnBody + "\n}"
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
