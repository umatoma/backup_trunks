package main

import (
	"flag"
	"log"
	"os"
	"encoding/base64"

	"github.com/umatoma/trunks/tasks"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
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
	tgt := []tasks.AttackTarget{
		tasks.AttackTarget{
			Method: "GET",
			URL: "http://localhost:8000/",
			Body: "",
		},
	}
	rate := uint64(1)
	duration := uint64(5)
	task, err := tasks.AttackTaskSignature(&tgt, rate, duration)
	if err != nil {
		log.Fatalln(err)
	}

	asyncResult, err := server.SendTask(task)
	if err != nil {
		log.Fatalln("Failed to send task", err)
	}

	result, err := asyncResult.Get()
	if err != nil {
		log.Fatalln("Getting task state failed with error", err)
	}

	bytesResult, err := base64.StdEncoding.DecodeString(result.String())
	if err != nil {
		log.Fatalln(err)
	}

	reporter, err := GetPlotReporter(&bytesResult)
	if err != nil {
		log.Fatalln(err)
	}
	reporter.Report(os.Stdout)

	// uuids, err := redisClient.GetAllTaskUUIDs()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(uuids)
}
