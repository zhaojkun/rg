package models

import (
	"errors"
	"go/ast"
	"strings"
)

var (
	HTTPMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
)

func parseHTTPMethod(fields []string) string {
	for _, field := range fields {
		name := strings.ToUpper(strings.TrimSpace(field))
		for _, hname := range HTTPMethods {
			if name == hname {
				return name
			}
		}
	}
	return "GET"
}

func parseURL(fields []string) (string, error) {
	for _, field := range fields {
		if strings.Contains(field, "/") {
			return strings.Trim(field, `" '`), nil
		}
	}
	return "", errors.New("url not found")
}

func checkFuncInterface(fnDecl *ast.FuncDecl, args ...string) bool {
	if len(args) != len(fnDecl.Type.Params.List) {
		return false
	}
	for i := 0; i < len(args); i++ {
		if !checkParameter(fnDecl.Type.Params.List[i].Type, args[i]) {
			return false
		}
	}
	return true
}

func checkParameter(pt ast.Expr, target string) bool {
	var src string
	var typ *ast.SelectorExpr
	stype, ok := pt.(*ast.StarExpr)
	if ok {
		src = "*"
		typ, ok = stype.X.(*ast.SelectorExpr)
		if !ok {
			return false
		}
	} else {
		typ, ok = pt.(*ast.SelectorExpr)
		if !ok {
			return false
		}
	}
	pkgIdent, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}
	if typ.Sel == nil {
		return false
	}
	src = src + pkgIdent.Name + "." + typ.Sel.Name
	return src == target
}
