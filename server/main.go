package main

import (
	"flag"
	"fmt"
	"log"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/signatures"
)

var (
	cnf config.Config
	server *machinery.Server
	redisClient *RedisClient
)

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
	server, err = machinery.NewServer(&cnf)
	if err != nil {
		log.Fatalln("Failed to initialize server", err)
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
	fmt.Println("Hello World")

	task := signatures.TaskSignature{
		Name: "add",
		Args: []signatures.TaskArg{
			signatures.TaskArg{
				Type: "int64",
				Value: 1,
			},
			signatures.TaskArg{
				Type: "int64",
				Value: 2,
			},
		},
	}

	asyncResult, err := server.SendTask(&task)
	if err != nil {
		log.Fatalln("Failed to send task", err)
	}

	result, err := asyncResult.Get()
	if err != nil {
		log.Fatalln("Getting task state failed with error", err)
	}
	fmt.Println(result.Interface())
	fmt.Println(asyncResult.Signature.UUID)

	uuids, err := redisClient.GetAllTaskUUIDs()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(uuids)
}
