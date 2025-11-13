package http

import (
	"github.com/gin-gonic/gin"
	"github.com/dangdinh2405/web-note-work/internal/data"
	"github.com/dangdinh2405/web-note-work/internal/handler"

)

func TasksRoutes(r *gin.Engine, db *data.Mongo, dbName string) {
	taskColl := db.DB(dbName).Collection("tasks")

	tasks := r.Group("/tasks")
	tasks.GET("/", handler.GetAllTasks(taskColl))
	tasks.POST("/", handler.CreateTask(taskColl))
	tasks.PUT("/:id", handler.UpdateTask(taskColl))
	tasks.DELETE("/:id", handler.DeleteTask(taskColl))
}