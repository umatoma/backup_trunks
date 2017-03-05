package main

import (
	"flag"
	"log"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	cnf config.Config
	jobQueueServer *machinery.Server
	redisClient *RedisClient
)

// BodyValidator is custom body validator
type BodyValidator struct {
	validator *validator.Validate
}

// NewBodyValidator is constructor of BodyValidator
func NewBodyValidator() *BodyValidator {
	return &BodyValidator{validator: validator.New()}
}

// Validate is an interface of echo.Validator
func (v *BodyValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func init() {
	var (
		url = flag.String("redis", "redis://127.0.0.1:6379", "Redis URL")
	)
	flag.Parse()

	cnf = config.Config{
		Broker: *url,
		ResultBackend: *url,
	}

	var err error
	jobQueueServer, err = machinery.NewServer(&cnf)
	if err != nil {
		log.Fatalln("Failed to initialize jobQueueServer", err)
	}

	var host, password, socketPath string
	var db int
	host, password, db, err = machinery.ParseRedisURL(*url)
	if err != nil {
		log.Fatalln(err)
	}
	redisClient = NewRedisClient(host, password, socketPath, db)
}

func main() {
	e := echo.New()
	e.Validator = NewBodyValidator()
	e.Use(middleware.Logger())
	e.GET("/", GetIndex)
	e.GET("/tasks", GetTasks)
	e.POST("/tasks", CreateTask)
	e.GET("/tasks/:id", GetTask)

	if err := e.Start(":3000"); err != nil {
		log.Fatalln(err)
	}
}
