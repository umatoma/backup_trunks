package main


import (
	"flag"
	"log"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

var (
	broker = flag.String("b", "redis://127.0.0.1:6379", "Broker URL")
	resultBackend = flag.String("r", "redis://127.0.0.1:6379", "Result backend")

	cnf config.Config
	server *machinery.Server
	worker *machinery.Worker
)

func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

func init() {
	flag.Parse()

	cnf = config.Config{
		Broker: *broker,
		ResultBackend: *resultBackend,
	}

	var err error
	server, err = machinery.NewServer(&cnf)
	if err != nil {
		log.Fatalln("Failed to initialize server", err)
	}

	server.RegisterTask("add", Add)

	worker = server.NewWorker("test_worker")
}

func main() {
	err := worker.Launch()
	if err != nil {
		log.Fatalln(err)
	}
}
