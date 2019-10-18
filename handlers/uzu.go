package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/clients"
	"github.com/mong0520/ChainChronicleGo/clients/general"
	"github.com/mong0520/ChainChronicleGo/clients/uzu"
	"github.com/mong0520/ChainChronicleGo/utils"
	log "github.com/sirupsen/logrus"
)

// UzuQueryHandler query the quest by title
func UzuQueryHandler(c *gin.Context) {
	sid := c.Query("sid")
	res := models.GeneralResponse{}
	logger := c.MustGet("logger").(*log.Logger)

	uzuData, _ := uzu.GetUzuInfo(sid)

	logger.Debug("enter UzuQueryHandler")

	type parsedUzuData struct {
		ID         int
		Name       string
		ScheduleID int
	}
	uzuResult := []parsedUzuData{}
	for _, uzu := range uzuData.Uzu {
		currentScheduleID := uzuData.GetCurrentScheduleID(uzu.UzuID)
		parsedUzuData := parsedUzuData{
			ID:         uzu.UzuID,
			Name:       uzu.Name,
			ScheduleID: currentScheduleID,
		}
		uzuResult = append(uzuResult, parsedUzuData)
	}

	res.Status = http.StatusOK
	res.Data = uzuResult

	c.JSON(200, res)
}

// UzuPlayHandler query the quest by title
func UzuPlayHandler(c *gin.Context) {
	// var logBuffer strings.Builder
	sid := c.Query("sid")
	uzid := c.Query("uzid")
	scid := c.Query("scid")
	res := models.GeneralResponse{}
	res.Status = http.StatusOK
	logger := c.MustGet("logger").(*log.Logger)

	apiEntry := "uzu/entry"
	apiRecover := "user/recover_uzu"
	apiResult := "uzu/result"

	stStart := 1
	stEnd := 12
	param := map[string]interface{}{
		"pt":    0,
		"fid":   1965350,
		"htype": 0,
		"uzid":  uzid,
		"scid":  scid,
	}
	autoRecover := true

	for i := stStart; i <= stEnd; i++ {
		param["st"] = strconv.Itoa(i)
		resp, ret := general.GeneralAction(apiEntry, sid, param)
		switch ret {
		case 0:
			logger.Debugf("開始挑戰天魔ID: %s, 第 %s 層", param["uzid"], param["st"])
		case 2803:
			logger.Debugf("無挑戰天魔ID: %s, 第 %s 層，挑戰權不足", param["uzid"], param["st"])
			if autoRecover {
				param := map[string]interface{}{}
				param["type"] = 1
				param["item_id"] = 35
				_, ret := general.GeneralAction(apiRecover, sid, param)
				if ret != 0 {
					logger.Error("回復挑戰權失敗")
					return
				}
				logger.Debug("挑戰權回復成功，重新挑戰")
				i--
				continue
			}
		case 2809:
			res.Message = resp["msg"].(string)
			logger.Debugf("無法挑戰天魔ID: %s, 第 %s 層", param["uzid"], param["st"])
			c.JSON(200, res)
			return
		default:
			res.Message = resp["msg"].(string)
			logger.Debugf("未知的錯誤")
			c.JSON(200, res)
			return
		}

		// End Quest
		paramResult := map[string]interface{}{}
		paramResult["res"] = 1
		paramResult["uzid"] = param["uzid"]
		// logger.Debugf("End UZU with Options %+v", paramResult)
		requestUrl := fmt.Sprintf("%s/%s", clients.HOST, apiResult)
		resp, _ = utils.PostV2(
			requestUrl,
			"wvt=%5b%7b%22wave_num%22%3a1%2c%22time%22%3a1383%7d%2c%7b%22wave_num%22%3a2%2c%22time%22%3a2849%7d%2c%7b%22wave_num%22%3a3%2c%22time%22%3a4316%7d%2c%7b%22wave_num%22%3a4%2c%22time%22%3a6856%7d%5d&mission=%7b%22cid%22%3a%5b2282%2c295%2c7634%2c5275%2c1245%2c8194%5d%2c%22sid%22%3a%5b0%2c0%2c296%2c5014%2c8131%2c8192%5d%2c%22fid%22%3a%5b8900%5d%2c%22hrid%22%3a%5b7206%5d%2c%22ms%22%3a0%2c%22md%22%3a3057%2c%22sc%22%3a%7b%220%22%3a0%2c%221%22%3a0%2c%222%22%3a0%2c%223%22%3a0%2c%224%22%3a0%7d%2c%22es%22%3a0%2c%22at%22%3a0%2c%22he%22%3a0%2c%22da%22%3a0%2c%22ba%22%3a0%2c%22bu%22%3a0%2c%22job%22%3a%7b%220%22%3a2%2c%221%22%3a4%2c%222%22%3a0%2c%223%22%3a2%2c%224%22%3a0%7d%2c%22weapon%22%3a%7b%220%22%3a1%2c%221%22%3a0%2c%222%22%3a0%2c%223%22%3a2%2c%224%22%3a0%2c%225%22%3a3%2c%228%22%3a0%2c%229%22%3a1%2c%2210%22%3a1%7d%2c%22box%22%3a0%2c%22um%22%3a%7b%221%22%3a0%2c%222%22%3a0%2c%223%22%3a0%7d%2c%22fj%22%3a0%2c%22fw%22%3a0%2c%22fo%22%3a0%2c%22mlv%22%3a100%2c%22mbl%22%3a445%2c%22udj%22%3a0%2c%22sdmg%22%3a98973%2c%22tp%22%3a0%2c%22gma%22%3a8%2c%22gmr%22%3a4%2c%22gmp%22%3a0%2c%22stp%22%3a0%2c%22auto%22%3a1%2c%22uh%22%3a%7b%226%22%3a1%2c%224%22%3a1%2c%222%22%3a1%2c%229%22%3a1%2c%223%22%3a2%2c%225%22%3a1%2c%221%22%3a1%7d%2c%22cc%22%3a1%2c%22bf_atk%22%3a0%2c%22bf_hp%22%3a0%2c%22bf_spd%22%3a0%7d&nature=cnt%3d16b68ab30b7%26mission%3d%257b%2522cid%2522%253a%255b2282%252c295%252c7634%252c5275%252c1245%252c8194%255d%252c%2522sid%2522%253a%255b0%252c0%252c296%252c5014%252c8131%252c8192%255d%252c%2522fid%2522%253a%255b8900%255d%252c%2522hrid%2522%253a%255b7206%255d%252c%2522ms%2522%253a0%252c%2522md%2522%253a3057%252c%2522sc%2522%253a%257b%25220%2522%253a0%252c%25221%2522%253a0%252c%25222%2522%253a0%252c%25223%2522%253a0%252c%25224%2522%253a0%257d%252c%2522es%2522%253a0%252c%2522at%2522%253a0%252c%2522he%2522%253a0%252c%2522da%2522%253a0%252c%2522ba%2522%253a0%252c%2522bu%2522%253a0%252c%2522job%2522%253a%257b%25220%2522%253a2%252c%25221%2522%253a4%252c%25222%2522%253a0%252c%25223%2522%253a2%252c%25224%2522%253a0%257d%252c%2522weapon%2522%253a%257b%25220%2522%253a1%252c%25221%2522%253a0%252c%25222%2522%253a0%252c%25223%2522%253a2%252c%25224%2522%253a0%252c%25225%2522%253a3%252c%25228%2522%253a0%252c%25229%2522%253a1%252c%252210%2522%253a1%257d%252c%2522box%2522%253a0%252c%2522um%2522%253a%257b%25221%2522%253a0%252c%25222%2522%253a0%252c%25223%2522%253a0%257d%252c%2522fj%2522%253a0%252c%2522fw%2522%253a0%252c%2522fo%2522%253a0%252c%2522mlv%2522%253a100%252c%2522mbl%2522%253a445%252c%2522udj%2522%253a0%252c%2522sdmg%2522%253a98973%252c%2522tp%2522%253a0%252c%2522gma%2522%253a8%252c%2522gmr%2522%253a4%252c%2522gmp%2522%253a0%252c%2522stp%2522%253a0%252c%2522auto%2522%253a1%252c%2522uh%2522%253a%257b%25226%2522%253a1%252c%25224%2522%253a1%252c%25222%2522%253a1%252c%25229%2522%253a1%252c%25223%2522%253a2%252c%25225%2522%253a1%252c%25221%2522%253a1%257d%252c%2522cc%2522%253a1%252c%2522bf_atk%2522%253a0%252c%2522bf_hp%2522%253a0%252c%2522bf_spd%2522%253a0%257d%26res%3d1%26uzid%3d5%26wvt%3d%255b%257b%2522wave_num%2522%253a1%252c%2522time%2522%253a1383%257d%252c%257b%2522wave_num%2522%253a2%252c%2522time%2522%253a2849%257d%252c%257b%2522wave_num%2522%253a3%252c%2522time%2522%253a4316%257d%252c%257b%2522wave_num%2522%253a4%252c%2522time%2522%253a6856%257d%255d",
			paramResult,
			sid)
		// logger.Info(utils.Map2JsonString(resp))
		ret = int(resp["res"].(float64))
		if ret != 0 {
			logger.Debug("魔神戰挑戰失敗")
			res.Message = resp["msg"].(string)
		} else {
			logger.Debug("魔神戰完成")
			res.Message = "魔神戰完成"
		}
	}

	c.JSON(200, res)
}

// QuestPlayHandler query the quest by title
// func QuestPlayHandler(c *gin.Context) {
// 	res := models.GeneralResponse{}
// 	questInfo := sdkQuest.NewQuest()

// 	sid := c.Query("sid")
// 	questInfo.QuestId, _ = strconv.Atoi(c.Query("qid"))
// 	questInfo.Type, _ = strconv.Atoi(c.Query("qtype"))
// 	questInfo.Pt, _ = strconv.Atoi(c.DefaultQuery("pt", "0"))
// 	questInfo.AutoSell = false
// 	questInfo.AutoBuy = false
// 	questInfo.AutoRaid = false
// 	questInfo.AutoRaidRecover = false
// 	metadata := &sdkClients.Metadata{Sid: sid}
// 	_, err := questInfo.StartQeust(metadata)

// 	if err != 0 {
// 		res.Status = http.StatusNotFound
// 		res.Error = string(err)
// 		res.Message = "關卡進入失敗"
// 		c.JSON(http.StatusNotFound, res)
// 	} else {
// 		// finish quest
// 		resp, err := questInfo.EndQeustV2(metadata)
// 		if err != 0 {
// 			res.Status = http.StatusNotFound
// 			res.Error = string(err)
// 			res.Message = "關卡完成失敗"
// 			c.JSON(http.StatusNotFound, res)
// 		} else {
// 			res.Status = http.StatusOK
// 			res.Data = resp
// 			res.Message = "關卡完成"
// 			c.JSON(200, res)
// 		}
// 	}

// }
