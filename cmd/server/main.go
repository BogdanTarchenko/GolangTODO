package main

import (
	"database/sql"
	"github.com/robfig/cron/v3"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "todo/docs"
	"todo/internal/delivery/http"
	"todo/internal/delivery/http/middleware"
	"todo/internal/repository"
	"todo/internal/usecase"
)

func main() {
	dsn := "host=localhost user=bogdantarchenko dbname=todo sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	taskRepo := repository.NewTaskPgRepository(db)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	taskHandler := http.NewTaskHandler(taskUsecase)

	c := cron.New()
	_, err = c.AddFunc("@every 1m", func() {
		log.Println("[CRON] Running UpdateOverdueTasks")
		if err := taskUsecase.UpdateOverdueTasks(); err != nil {
			log.Printf("[CRON] Failed to update overdue tasks: %v", err)
		} else {
			log.Println("[CRON] UpdateOverdueTasks completed successfully")
		}
	})
	if err != nil {
		log.Fatalf("failed to schedule cron job: %v", err)
	}
	c.Start()

	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	taskHandler.RegisterRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
