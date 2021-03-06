package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/clients/general"
	"github.com/mong0520/ChainChronicleGo/clients/web"
	"github.com/mong0520/ChainChronicleGo/utils"
	log "github.com/sirupsen/logrus"
)

// EventsHandler char the quest by title
func EventsHandler(c *gin.Context) {
	res := models.GeneralResponse{}
	sid := c.Query("sid")
	logger := c.MustGet("logger").(*log.Logger)

	param := map[string]interface{}{}
	ret, _ := general.GeneralAction("data/eventPortal", sid, param)
	res.Data = ret
	events := models.Events{}
	err := utils.Map2Struct(ret, &events)
	if err != nil {
		res.Data = err
	}

	results := make([]models.Banner, 0, 0)
	now := time.Now() // current local time
	sec := int(now.Unix())
	for _, d := range events.EventPortal {
		for _, event := range d.Events {
			for _, banner := range event.Banner {
				if banner.Type == "PremiumGacha" && sec >= banner.Start && sec <= banner.End {
					// logger.Debug("--------------------------")
					// logger.Debug("gacha type = ", banner.GachaType)
					// logger.Debug("gacha id = ", banner.GachaID)
					banner.BannerURL = fmt.Sprintf("http://content.cc.mobimon.com.tw/3818/Prod/Resource/Banner/%s", banner.Name)
					banner.InfoURL = fmt.Sprintf("http://v3818.cc.mobimon.com.tw/web/gacha?type=%s&gacha_id=%d", banner.GachaType, banner.GachaID)

					gachasInfo, err := web.GetGachaInfo(banner.GachaType, sid, banner.GachaID)
					if err != nil {
						logger.Error(err)
					}

					banner.GachaInfo = gachasInfo

					results = append(results, banner)
				}
			}
		}
	}

	// FIXME: very hard code to show Majin Gacha info
	gachasInfo, err := web.GetGachaInfo("raidex", sid, 0)
	if err != nil {
		logger.Error(err)
	}

	tmpBanner := models.Banner{}
	tmpBanner.GachaInfo = gachasInfo

	results = append(results, tmpBanner)

	res.Data = results

	c.JSON(200, res)
}
