package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icza/dyno"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/clients"
	sdkClients "github.com/mong0520/ChainChronicleGo/clients"
	"github.com/mong0520/ChainChronicleGo/clients/card"
	"github.com/mong0520/ChainChronicleGo/clients/gacha"
	"github.com/mong0520/ChainChronicleGo/controllers"
	sdkModels "github.com/mong0520/ChainChronicleGo/models"
	"github.com/mong0520/ChainChronicleGo/utils"
	"github.com/oleiade/reflections"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GachaHandler char the quest by title
func GachaHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	sid := c.Query("sid")
	gachaID := c.Query("gacha_id")
	gachaCount := c.Query("gacha_count")
	logger := c.MustGet("logger").(*log.Logger)
	conn := c.MustGet("databaseConn").(*mgo.Session)
	metadata := &sdkClients.Metadata{
		Sid: sid,
		DB:  conn,
	}

	defaultFilteredKeys := map[string]string{
		"rarity": "Rarity",
		"title":  "Title",
		"name":   "Name",
	}

	var result strings.Builder
	gachaInfo := gacha.NewGachaInfo()
	gachaInfo.Type, _ = strconv.Atoi(gachaID)
	gachaInfo.Count, _ = strconv.Atoi(gachaCount)
	gachaInfo.AutoSellRarityThreshold = 4
	gachaInfo.Verbose = true
	gachaInfo.AutoSell = true
	if resp, ret := gachaInfo.Gacha(metadata); ret == 0 {
		gachaResult := processGachaResult(resp, conn, logger)
		filteredResult := make([]map[string]interface{}, 0, 0)
		for _, card := range gachaResult["char"].([]sdkModels.GachaResultChara) {

			myCard := sdkModels.Charainfo{}
			query := bson.M{"cid": card.ID}
			controllers.GeneralQuery(metadata.DB, "charainfo", query, &myCard)
			msg := fmt.Sprintf("得到 %d星卡: %s-%s", myCard.Rarity, myCard.Title, myCard.Name)
			result.WriteString(msg)
			if gachaInfo.AutoSell && myCard.Rarity <= gachaInfo.AutoSellRarityThreshold {
				doSellItem(metadata, card.Idx, "", logger)
			}

			tempResult := make(map[string]interface{})
			for _, v := range defaultFilteredKeys {
				tempResult[v], _ = reflections.GetField(myCard, v)
			}
			filteredResult = append(filteredResult, tempResult)
		}
		res.Data = filteredResult
		res.Status = http.StatusOK
	} else {
		result.WriteString(utils.Map2JsonString(resp))
	}

	res.Message = result.String()

	c.JSON(200, res)
}

func doSellItem(metadata *clients.Metadata, cid int, section string, logger *log.Logger) {
	if ret, err := card.Sell(metadata, cid); err != 0 {
		logger.Info("\t-> 賣出卡片失敗", utils.Map2JsonString(ret))
	} else {
		logger.Info("\t-> 賣出卡片成功")
	}
}

func processGachaResult(resp map[string]interface{}, session *mgo.Session, logger *log.Logger) (gachaResult map[string]interface{}) {
	gachaData, _ := dyno.GetSlice(resp, "body")
	//logger.Info(utils.Map2JsonString(resp))
	gachaResult = map[string]interface{}{
		"char":   []sdkModels.GachaResultChara{},
		"item":   []sdkModels.GachaResultItem{},
		"weapon": []sdkModels.GachaResultWeapon{},
	}

	gachaResult["char"] = []sdkModels.GachaResultChara{}
	gachaResult["item"] = []sdkModels.GachaResultItem{}
	gachaResult["weapon"] = []sdkModels.GachaResultWeapon{}

	charList := make([]sdkModels.GachaResultChara, 0)
	itemList := make([]sdkModels.GachaResultItem, 0)
	weaponList := make([]sdkModels.GachaResultWeapon, 0)

	for i, data := range gachaData {
		if i == 0 {
			continue
		}
		dataType, _ := dyno.GetFloat64(data, "type")
		switch dataType {
		case 15:
			logger.Info(i, "Type 15", data)
		case 1:
			//logger.Info(i, "得到角色")
			list := data.(map[string]interface{})["data"].([]interface{})
			for _, item := range list {
				tmpItem := &sdkModels.GachaResultChara{}
				tmpDBItem := &sdkModels.Charainfo{}
				if err := utils.Map2Struct(item.(map[string]interface{}), tmpItem); err != nil {
					logger.Info("Unable to convert to struct", err)
				} else {
					query := bson.M{"cid": tmpItem.ID}
					if err := controllers.GeneralQuery(session, "charainfo", query, tmpDBItem); err != nil {
						logger.Info(i, "得到", tmpItem.ID)
					} else {
						logger.Infof("得到 %d星卡: %s-%s", tmpDBItem.Rarity, tmpDBItem.Title, tmpDBItem.Name)
					}
					charList = append(charList, *tmpItem)
				}
			}
		case 2:
			//logger.Info(i, "得到成長卡/冶鍊卡", data)
			list := data.(map[string]interface{})["data"].([]interface{})
			for _, item := range list {
				tmpItem := &sdkModels.GachaResultItem{}
				tmpDBItem := &sdkModels.Chararein{}
				if err := utils.Map2Struct(item.(map[string]interface{}), tmpItem); err != nil {
					logger.Info("Unable to convert to struct", err)
				} else {
					query := bson.M{"id": tmpItem.ItemID}
					controllers.GeneralQuery(session, "chararein", query, tmpDBItem)
					logger.Info(i, "得到", tmpDBItem.Name)
					itemList = append(itemList, *tmpItem)
				}
			}
		case 14:
			continue
		case 105:
			list := data.(map[string]interface{})["data"].([]interface{})
			for _, item := range list {
				tmpItem := &sdkModels.GachaResultWeapon{}
				tmpDBItem := &sdkModels.Weapon{}
				if err := utils.Map2Struct(item.(map[string]interface{}), tmpItem); err != nil {
					logger.Info("Unable to convert to struct", err)
				} else {
					query := bson.M{"id": tmpItem.ItemID}
					if err := controllers.GeneralQuery(session, "evolve", query, tmpDBItem); err != nil {
						logger.Info(i, "得到", tmpItem.ItemID)
					} else {
						logger.Info(i, "得到", tmpDBItem.Name)
					}
					weaponList = append(weaponList, *tmpItem)
				}
			}
		default:
			logger.Info(dataType)
		}
	}
	gachaResult["char"] = charList
	gachaResult["item"] = itemList

	return gachaResult
}
