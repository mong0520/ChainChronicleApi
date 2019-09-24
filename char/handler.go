package char

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/controllers"
	sdkModel "github.com/mong0520/ChainChronicleGo/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CharQueryHandler char the quest by title
func CharQueryHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	conn := c.MustGet("databaseConn").(*mgo.Session)
	charName := c.DefaultQuery("name", "")
	query := bson.M{
		"name": bson.RegEx{Pattern: charName, Options: "i"},
	}
	results := []sdkModel.Charainfo{}
	err := controllers.GeneralQueryAll(conn, "charainfo", query, &results)
	if err != nil {
		res.Status = http.StatusNotFound
		res.Error = err.Error()
	} else {
		res.Status = http.StatusOK
		res.Data = results
	}

	c.JSON(200, res)
}

func MyCardHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	conn := c.MustGet("databaseConn").(*mgo.Session)
	charName := c.DefaultQuery("name", "")
	query := bson.M{
		"name": bson.RegEx{Pattern: charName, Options: "i"},
	}
	results := []sdkModel.Charainfo{}
	err := controllers.GeneralQueryAll(conn, "charainfo", query, &results)
	if err != nil {
		res.Status = http.StatusNotFound
		res.Error = err.Error()
	} else {
		res.Status = http.StatusOK
		res.Data = results
	}

	c.JSON(200, res)
}
