package main

import (
	"net/http"
	"encoding/base64"
	"bytes"
	"errors"

	"github.com/umatoma/trunks/tasks"

	"github.com/labstack/echo"
)

// BodyCreateTask is body schema for CreateTask
type BodyCreateTask struct {
	Targets []tasks.AttackTarget	`validate:"required,dive,required"`
	Rate uint64 									`validate:"required"`
	Duration uint64 							`validate:"required"`
}

func getTaskResultBytes(taskID string) (*[]byte, error) {
	result, err := jobQueueServer.GetBackend().GetState(taskID)
	if err != nil {
		return nil, err
	}

	if !result.IsCompleted() {
		return nil, err
	}

	encodedResult, ok := result.Result.Value.(string)
	if !ok {
		return nil, errors.New("result data is borken")
	}

	bytesResult, err := base64.StdEncoding.DecodeString(encodedResult)
	if err != nil {
		return nil, err
	}

	return &bytesResult, nil
}

// GetIndex handle GET /
func GetIndex(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!!")
}

// GetPendingTasks handle GET /tasks/pending
func GetPendingTasks(c echo.Context) error {
	queue := jobQueueServer.GetConfig().DefaultQueue
	taskSignatures, err := jobQueueServer.GetBroker().GetPendingTasks(queue)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, taskSignatures)
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

// GetTaskTextRepot handle GET /tasks/:id/report/text
func GetTaskTextRepot(c echo.Context) error {
	taskID := c.Param("id")
	bytesResult, err := getTaskResultBytes(taskID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	reporter, err := GetTextReporter(bytesResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var buf bytes.Buffer
	if err := reporter.Report(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

// GetTaskPlotRepot handle GET /tasks/:id/report/plot
func GetTaskPlotRepot(c echo.Context) error {
	taskID := c.Param("id")
	bytesResult, err := getTaskResultBytes(taskID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	reporter, err := GetPlotReporter(bytesResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var buf bytes.Buffer
	if err := reporter.Report(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}
