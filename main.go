package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/handlers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// APIMiddleware will add the db connection to the context
func APIMiddleware(db *mgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.New()
		logger.SetLevel(log.DebugLevel)
		c.Set("databaseConn", db)
		c.Set("logger", logger)
		c.Next()
	}
}

func main() {
	router := gin.New()
	// router.Use(cors.New(cors.Config{
	// 	AllowOriginFunc:  func(origin string) bool { return true },
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	router.Use(cors.Default())
	dbMode := os.Getenv("DBMode")

	if dbMode == "docker" {
		conn, err := mgo.Dial("mongodb:27017")
		if err != nil {
			log.Fatal(err)
		}
		router.Use(APIMiddleware(conn))
	} else if dbMode == "local" {
		conn, err := mgo.Dial("localhost:27017")
		if err != nil {
			log.Fatal(err)
		}
		router.Use(APIMiddleware(conn))
	} else if dbMode == "remote" {
		dbURI := os.Getenv("DBUri")
		dialInfo, _ := mgo.ParseURL(dbURI)
		tlsConfig := &tls.Config{}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		if conn, err := mgo.DialWithInfo(dialInfo); err != nil {
			log.Errorf("Unable to connect DB %s, Dialnfo = %+v", err, dialInfo)
		} else {
			router.Use(APIMiddleware(conn))
		}
	} else {
		log.Fatalf("unsupported DBMode: %s", dbMode)
	}

	// router.Static("/web", "./web")
	router.GET("/login", handlers.LoginHandler)
	router.GET("/status", handlers.StatusHandler)
	router.GET("/query_quest", handlers.QuestQueryHandler)
	router.GET("/play_quest", handlers.QuestPlayHandler)
	router.GET("/char", handlers.CharQueryHandler)
	router.GET("/query_uzu", handlers.UzuQueryHandler)
	router.GET("/play_uzu", handlers.UzuPlayHandler)
	router.GET("/gacha", handlers.GachaHandler)
	router.GET("/events", handlers.EventsHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	router.Run(addr)
}
