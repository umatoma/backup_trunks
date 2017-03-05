package main


import (
	"flag"
	"log"

	"github.com/umatoma/trunks/tasks"

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

	server.RegisterTask("add", tasks.TaskAdd)
	server.RegisterTask("attack", tasks.TaskAttack)

	worker = server.NewWorker("test_worker")
}

func main() {
	err := worker.Launch()
	if err != nil {
		log.Fatalln(err)
	}
}
