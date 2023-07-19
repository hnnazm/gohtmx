package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Status uint8

const (
	INCOMPLETE Status = iota
	COMPLETED
)

type Todo struct {
	ID     int
	Title  string
	Status Status
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func RemoveIndex(s []Todo, index int) []Todo {
    ret := make([]Todo, 0)
    ret = append(ret, s[:index]...)

    return append(ret, s[index+1:]...)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}
	e.Renderer = renderer

	todos := []Todo{
		{ID: 0, Title: "Todo 1", Status: 0},
		{ID: 1, Title: "Todo 2", Status: 1},
		{ID: 2, Title: "Todo 3", Status: 1},
	}

	todo := e.Group("todo")

	todo.GET("/list", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", &todos)
	})

	todo.POST("/create", func(c echo.Context) error {
		title := c.FormValue("title")

		if title == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		newTodo := Todo{ID: len(todos), Title: title, Status: 0}
		todos = append(todos, newTodo)

		return c.Render(http.StatusOK, "todo.html", newTodo)
	})

	todo.GET("/update", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	todo.DELETE("/delete/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

    if err != nil {
			return c.NoContent(http.StatusBadRequest)
    }

    todos = RemoveIndex(todos, id)

    if err != nil {
			return c.NoContent(http.StatusBadRequest)
    }

    fmt.Println(todos)

		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.PUT("/messages", func(c echo.Context) error {
		return c.Render(http.StatusOK, "message.html", map[string]interface{}{
			"name": "Dolly!",
		})
	}).Name = "foobar"

	e.Logger.Fatal(e.Start(":1323"))
}
