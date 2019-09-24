package quest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/controllers"
	"gopkg.in/mgo.v2"

	sdkClients "github.com/mong0520/ChainChronicleGo/clients"
	sdkQuest "github.com/mong0520/ChainChronicleGo/clients/quest"
)

// QuestQueryHandler query the quest by title
func QuestQueryHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	conn := c.MustGet("databaseConn").(*mgo.Session)
	questName := c.DefaultQuery("name", "")
	questInfo, err := controllers.GetQuestsByName(conn, questName)
	if err != nil {
		res.Status = http.StatusNotFound
		res.Error = err.Error()
	} else {
		res.Status = http.StatusOK
		res.Data = questInfo
	}

	c.JSON(200, res)
}

// QuestPlayHandler query the quest by title
func QuestPlayHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	questInfo := sdkQuest.NewQuest()

	sid := c.Query("sid")
	questInfo.QuestId, _ = strconv.Atoi(c.Query("qid"))
	questInfo.Type, _ = strconv.Atoi(c.Query("qtype"))
	questInfo.Pt, _ = strconv.Atoi(c.DefaultQuery("pt", "0"))
	questInfo.AutoSell = false
	questInfo.AutoBuy = false
	questInfo.AutoRaid = false
	questInfo.AutoRaidRecover = false
	metadata := &sdkClients.Metadata{Sid: sid}
	_, err := questInfo.StartQeust(metadata)

	if err != 0 {
		res.Status = http.StatusNotFound
		res.Error = string(err)
		c.JSON(200, res)
	}

	// finish quest
	resp, err := questInfo.EndQeustV2(metadata)
	if err != 0 {
		res.Status = http.StatusNotFound
		res.Error = string(err)
		c.JSON(200, res)
	}

	res.Status = http.StatusOK
	res.Data = resp

	c.JSON(200, res)
}
