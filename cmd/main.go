package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mayank18049/GoChessMoveShare/internal/adapters/gamehandler/http"
	"github.com/mayank18049/GoChessMoveShare/internal/adapters/gamerepo/memory"
	"github.com/mayank18049/GoChessMoveShare/internal/adapters/messagebroker/rabbitmq"
	"github.com/mayank18049/GoChessMoveShare/internal/domain/ports"
	"github.com/mayank18049/GoChessMoveShare/internal/service"
)

func main() {
	l := log.New(os.Stdout, "MovePublisher-API", log.LstdFlags)

	env, err := godotenv.Read()

	if err != nil {
		l.Fatalf("[main]: Failed to read .env properly %s", err.Error())
	}
	l.Printf("Data %s\n", env)

	var broker ports.MessageBroker = rabbitmq.NewRabbitMQMessageBroker(context.TODO(), "amqp://"+env["RABBITMQ_USER"]+":"+env["RABBITMQ_PASSWORD"]+"@"+env["RABBITMQ_URI"]+"/", l)
	var gamerepo ports.GameRepo = memory.NewInMemoryGameRepo(l)
	registerService, _ := service.NewGameRegistration(broker, gamerepo, l)
	httpHandler, _ := http.NewHttpHandler(*registerService, l)
	httpHandler.SetAndServe(env["SERVER_ADDR"])
}
