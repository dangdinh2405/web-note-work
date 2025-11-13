package main

import (
	"os"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/dangdinh2405/web-note-work/internal/data"
	"github.com/dangdinh2405/web-note-work/internal/http"
)

func main() {
	godotenv.Load()


	gin.SetMode(gin.ReleaseMode) 

	// Production:
	// nhớ cấu hình reverse proxy (X-Forwarded-Proto) nếu đứng sau CDN/Load Balancer

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	
	db, err := data.NewMongo(os.Getenv("MONGO_CONECTION"))
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	// if dbName == "" {
	// 	dbName = "tasks"
	// }

	http.TasksRoutes(r, db, dbName)
	

	defer db.Close()

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}