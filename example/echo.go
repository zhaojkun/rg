package example

import (
	"strconv"

	"github.com/labstack/echo"
)

// PageInfo
// @Parameter {Number} id
// @Request get /api/page/:id/pageinfo
//
func PageInfo(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(200, id)
}

// GetTodos get "/api/todos"
func GetTodos(ctx echo.Context) error {
	return nil
}

// AddTodo post "/api/todos"
func AddTodo(ctx echo.Context) error {
	return nil
}
