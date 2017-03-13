package models

import "go/ast"

type Handler interface {
	Pkg() string
	FuncParam() string
	Path() string
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

type HandlerSorter []Handler

func (s HandlerSorter) Len() int      { return len(s) }
func (s HandlerSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s HandlerSorter) Less(i, j int) bool {
	p1, p2 := s[i].Path(), s[j].Path()
	if len(p1) == len(p2) {
		return s[i].Path() < s[j].Path()
	}
	return len(p1) < len(p2)
}

func init() {
	frameworkHandler = make(map[string]HandlerBuilder)
	Register("github.com/labstack/echo", echoHandlerBuilder)
	Register("github.com/gin-gonic/gin", ginHandlerBuilder)
	Register("net/http", originHandlerBuilder)
}
