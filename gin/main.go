package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todo/gin/internal/consumer"
	"todo/gin/internal/handlers"
	"todo/gin/internal/repository"
	"todo/gin/internal/services"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/gin-gonic/gin"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "appdb"
)

func main() {
	DSN := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"tasks_exchange",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare(
		"tasks", // имя очереди
		true,    // durable (сохраняется после перезапуска)
		false,   // auto delete
		false,   // exclusive
		false,   // no wait
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		"tasks",          // очередь
		"",               // routing key (не нужен для fanout)
		"tasks_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresTaskRepository(db)
	svc := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(svc)

	ctx := context.Background()

	cons := consumer.NewConsumer(ch)

	go func() {
		if err := cons.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	server := gin.Default()

	server.Any("/ping", handler.Ping)

	tasks := server.Group("/tasks")
	{
		tasks.POST("", handler.CreateTask)
		tasks.GET("", handler.GetTasks)
		tasks.PATCH("/:id", handler.UpdateTask)
		tasks.GET("/:id", handler.GetTaskByID)
		tasks.DELETE("/:id", handler.DeleteTask)
	}

	// ":8080"
	server.Run()
}
