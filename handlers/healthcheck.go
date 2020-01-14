package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
)

func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, models.GeneralResponse{})
}
