package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kirill-Sirotkin/input_api_go/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	HandlePostTask(c echo.Context) error
	HandleGetTask(c echo.Context) error
}

type RouteHandler struct {
	TaskMap     *models.TaskMap
	TaskLimiter chan bool
}

func NewRouteHandler(taskMap *models.TaskMap, taskLimiter chan bool) *RouteHandler {
	return &RouteHandler{
		TaskMap:     taskMap,
		TaskLimiter: taskLimiter,
	}
}

func (h *RouteHandler) HandlePostTask(c echo.Context) error {
	taskName := models.PostTaskDTO{}

	if err := c.Bind(&taskName); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input. Task name must be a string"})
	}
	if taskName.Name == "" {
		return c.JSON(400, map[string]string{"error": "Invalid input. Task name cannot be empty"})
	}

	task := models.NewTask(taskName.Name, "pending", "")
	h.TaskMap.PostTask(task)

	h.TaskLimiter <- true
	go func() {
		log.Printf("start time: %v", time.Now())
		mockIOTask(task.Id, h.TaskMap)
		log.Printf("end time: %v", time.Now())
		<-h.TaskLimiter
	}()

	return c.JSON(201, task)
}

func (h *RouteHandler) HandleGetTask(c echo.Context) error {
	id := c.Param("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		errorMsg := fmt.Sprintf("Invalid UUID format: %s", err.Error())
		return c.JSON(400, map[string]string{"error": errorMsg})
	}

	task, err := h.TaskMap.GetTaskById(userUUID)
	if err != nil {
		errorMsg := fmt.Sprintf("Task GET request error: %s", err.Error())
		return c.JSON(404, map[string]string{"error": errorMsg})
	}
	return c.JSON(200, task)
}

func mockIOTask(taskId uuid.UUID, taskMap *models.TaskMap) {
	time.Sleep(10 * time.Second)

	filePath := fmt.Sprintf("files/task_%s.json", taskId.String())

	mockDataFirstIOStep := map[string]string{
		"1": fmt.Sprintf("first step for %s", taskId.String()),
	}

	jsonData, err := json.MarshalIndent(mockDataFirstIOStep, "", "    ")
	if err != nil {
		log.Printf("%v", err)
	}
	os.WriteFile(filePath, jsonData, 0644)

	time.Sleep(10 * time.Second)

	mockDataSecondIOStep := map[string]string{
		"1": fmt.Sprintf("first step for %s", taskId.String()),
		"2": fmt.Sprintf("second step for %s", taskId.String()),
	}

	jsonData, err = json.MarshalIndent(mockDataSecondIOStep, "", "    ")
	if err != nil {
		log.Printf("%v", err)
	}
	os.WriteFile(filePath, jsonData, 0644)

	taskMap.UpdateTaskStatus(taskId, filePath)
}
