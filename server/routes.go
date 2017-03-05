package main

import (
	"net/http"
	"encoding/base64"
	"bytes"

	"github.com/umatoma/trunks/tasks"

	"github.com/labstack/echo"
)

// BodyCreateTask is body schema for CreateTask
type BodyCreateTask struct {
	Targets []tasks.AttackTarget	`json:"targets" validate:"required,dive,required"`
	Rate uint64 									`json:"rate" validate:"required"`
	Duration uint64 							`json:"duration" validate:"required"`
}

// GetIndex handle GET /
func GetIndex(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!!")
}

// GetTasks handle GET /tasks
func GetTasks(c echo.Context) error {
	uuids, err := redisClient.GetAllTaskUUIDs()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, uuids)
}

// CreateTask handle POST /tasks
func CreateTask(c echo.Context) error {
	var body BodyCreateTask

	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	task, err := tasks.AttackTaskSignature(&body.Targets, body.Rate, body.Duration)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	asyncResult, err := jobQueueServer.SendTask(task)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, asyncResult.Signature)
}

// GetTask handle GET /tasks/:id
func GetTask(c echo.Context) error {
	taskID := c.Param("id")
	result, err := jobQueueServer.GetBackend().GetState(taskID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !result.IsCompleted() {
		return c.String(http.StatusInternalServerError, "Status is not completed")
	}

	encodedResult, ok := result.Result.Value.(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "Result data is borken.")
	}

	bytesResult, err := base64.StdEncoding.DecodeString(encodedResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	reporter, err := GetPlotReporter(&bytesResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var buf bytes.Buffer
	if err := reporter.Report(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}
