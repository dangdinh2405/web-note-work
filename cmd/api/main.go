package main

import (
	"os"
	"log"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

	"github.com/dangdinh2405/web-note-work/internal/data"
	"github.com/dangdinh2405/web-note-work/internal/http"
)

func main() {
	godotenv.Load()

	// Production:

	// gin.SetMode(gin.ReleaseMode) #Deloy must turn on
	// nhớ cấu hình reverse proxy (X-Forwarded-Proto) nếu đứng sau CDN/Load Balancer

	r := gin.Default()
	r.SetTrustedProxies(nil)

	port := os.Getenv("PORT")
	
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