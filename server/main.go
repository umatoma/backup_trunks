package main

import (
	"flag"
	"fmt"
	"log"
	"encoding/json"

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

type AttackTarget struct {
	Method string `json:"method"`
	URL string `json:"url"`
	Body string `json:"body"`
	Headers map[string][]string `json:"headers"`
}

func main() {
	fmt.Println("Hello World")

	tgt := []AttackTarget{
		AttackTarget{
			Method: "GET",
			URL: "http://localhost:8000/",
			Body: "",
		},
	}
	jsonTgt, err := json.Marshal(&tgt)
	if err != nil {
		log.Fatalln(err)
	}
	rate := uint64(1)
	duration := uint64(5)

	task := signatures.TaskSignature{
		Name: "attack",
		Args: []signatures.TaskArg{
			signatures.TaskArg{
				Type: "string",
				Value: string(jsonTgt),
			},
			signatures.TaskArg{
				Type: "uint64",
				Value: rate,
			},
			signatures.TaskArg{
				Type: "uint64",
				Value: duration,
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
