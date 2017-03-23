package example

import "github.com/labstack/echo"

type Server struct{}

// GetTodos get "/api/m/todos"
func (s *Server) GetTodos(ctx echo.Context) error {
	return nil
}

// AddTodo post "/api/m/todos"
func (s *Server) AddTodo(ctx echo.Context) error {
	return nil
}
