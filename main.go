package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/handlers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// APIMiddleware will add the db connection to the context
func APIMiddleware(db *mgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.New()
		c.Set("databaseConn", db)
		c.Set("logger", logger)
		c.Next()
	}
}

func main() {

	router := gin.New()
	conn := &mgo.Session{}
	conn, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	router.Use(APIMiddleware(conn))

	router.GET("/login", handlers.LoginHandler)
	router.GET("/status", handlers.StatusHandler)
	router.GET("/query_quest", handlers.QuestQueryHandler)
	router.GET("/play_quest", handlers.QuestPlayHandler)
	router.GET("/char", handlers.CharQueryHandler)

	router.Run(":5000")
}
