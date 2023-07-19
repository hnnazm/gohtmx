package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Status uint8

const (
	INCOMPLETE Status = iota
	COMPLETED
)

type Todo struct {
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

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}
	e.Renderer = renderer

	todos := []Todo{
		{Title: "Todo 1", Status: 0},
		{Title: "Todo 2", Status: 1},
		{Title: "Todo 3", Status: 1},
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

		newTodo := Todo{Title: title, Status: 0}
		todos = append(todos, newTodo)

		return c.Render(http.StatusOK, "todo.html", newTodo)
	})
	todo.GET("/update", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	todo.GET("/delete", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
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
