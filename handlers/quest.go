package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/controllers"
	"github.com/oleiade/reflections"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	sdkClients "github.com/mong0520/ChainChronicleGo/clients"
	sdkQuest "github.com/mong0520/ChainChronicleGo/clients/quest"
)

// QuestQueryHandler query the quest by title
func QuestQueryHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	conn := c.MustGet("databaseConn").(*mgo.Session)
	logger := c.MustGet("logger").(*log.Logger)
	defaultFilteredKeys := map[string]string{
		"questid":   "QuestID",
		"questtype": "QuestType",
		"name":      "Name",
	}
	defaultSelector := bson.M{}
	for k := range defaultFilteredKeys {
		defaultSelector[k] = 1
	}
	questName := c.DefaultQuery("name", "")
	questInfo, err := controllers.GetQuestsByName(conn, questName, &defaultSelector)

	filteredResult := make([]map[string]interface{}, 0, 0)

	for _, quest := range questInfo {
		tempResult := make(map[string]interface{})
		for _, v := range defaultFilteredKeys {
			tempResult[v], err = reflections.GetField(quest, v)
			if err != nil {
				logger.Error(err)
			}
		}
		filteredResult = append(filteredResult, tempResult)
	}

	if err != nil {
		res.Status = http.StatusNotFound
		res.Error = err.Error()
	} else {
		res.Status = http.StatusOK
		res.Data = filteredResult
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
		res.Message = "關卡進入失敗"
		c.JSON(http.StatusNotFound, res)
	} else {
		// finish quest
		resp, err := questInfo.EndQeustV2(metadata)
		if err != 0 {
			res.Status = http.StatusNotFound
			res.Error = string(err)
			res.Message = "關卡完成失敗"
			c.JSON(http.StatusNotFound, res)
		} else {
			res.Status = http.StatusOK
			res.Data = resp
			res.Message = "關卡完成"
			c.JSON(200, res)
		}
	}

}
