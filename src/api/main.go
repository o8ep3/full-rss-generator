package main

import (
	"fullRssAPI/config"
	"fullRssAPI/model/db"
	"fullRssAPI/server/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	sqlDB, err := connectToDB(
		config.DBInfo.Host,
		config.DBInfo.Port,
		config.DBInfo.User,
		config.DBInfo.Password,
		config.DBInfo.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}
	data := db.NewSQLDataStorage(sqlDB)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowMethods = []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"DELETE",
	}
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Content-Type"}
	r.Use(cors.New(config))

	r.POST("/api/feedinfo", handlers.Register(data))
	r.GET("/api/feedinfo", handlers.GetFeedInfo(data))
	r.DELETE("/api/feedinfo", handlers.Delete(data))
	r.GET("/api/rss/:id", handlers.GetRSS)
	r.PUT("/api/feedinfo", handlers.Refresh(data))
	r.Run()
}
