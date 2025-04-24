package main

import (
	"github.com/Kirill-Sirotkin/input_api_go/handlers"
	"github.com/Kirill-Sirotkin/input_api_go/models"
	"github.com/labstack/echo/v4"
)

const MAX_TASKS int = 10

func main() {
	e := echo.New()
	e.Static("/files", "files")

	taskLimiter := make(chan bool, MAX_TASKS)
	taskMap := models.NewTaskMap()
	var h handlers.Handler = handlers.NewRouteHandler(taskMap, taskLimiter)

	e.POST("/create-task", h.HandlePostTask)
	e.GET("/poll-task/:id", h.HandleGetTask)

	e.Logger.Fatal(e.Start(":1323"))
}
