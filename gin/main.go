package main

import (
	"database/sql"
	"fmt"
	"log"
	"todo/gin/internal/handlers"
	"todo/gin/internal/repository"
	"todo/gin/internal/services"

	_ "github.com/lib/pq"

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

	repo := repository.NewPostgresTaskRepository(db)
	svc := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(svc)

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
