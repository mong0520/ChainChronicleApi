package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/char"
	"github.com/mong0520/ChainChronicleApi/quest"
	"github.com/mong0520/ChainChronicleApi/session"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// APIMiddleware will add the db connection to the context
func APIMiddleware(db *mgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
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

	router.GET("/login", session.LoginHandler)
	router.GET("/status", session.StatusHandler)
	router.GET("/quest", quest.QuestQueryHandler)
	router.GET("/playquest", quest.QuestPlayHandler)
	router.GET("/char", char.CharQueryHandler)

	router.Run(":8000")
}
