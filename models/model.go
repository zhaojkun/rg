package models

import "go/ast"

type Handler interface {
	Pkg() string
	FuncParam() string
	String() string
}

type HandlerBuilder func(*ast.FuncDecl) (Handler, error)

var frameworkHandler map[string]HandlerBuilder

func Register(path string, h HandlerBuilder) {
	frameworkHandler[path] = h
}

func GetBuilder(paths []string) (HandlerBuilder, bool) {
	for _, p := range paths {
		builder, ok := frameworkHandler[p]
		if ok {
			return builder, true
		}
	}
	return nil, false
}

func init() {
	frameworkHandler = make(map[string]HandlerBuilder)
	Register("github.com/labstack/echo", echoHandlerBuilder)
}
