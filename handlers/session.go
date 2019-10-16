package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mong0520/ChainChronicleApi/models"
	"github.com/mong0520/ChainChronicleGo/clients/session"
)

// LoginHandler returns session ID
func LoginHandler(c *gin.Context) {

	uid := c.DefaultQuery("uid", "iOS7069F115-D7DC-48F3-8658-04191F4AA949")
	token := c.DefaultQuery("token", "")
	res := models.GeneralResponse{}

	if ret, err := session.Login(uid, token, false); err != nil {
		res.Status = http.StatusNotFound
		res.Error = err.Error()
	} else {
		res.Status = http.StatusOK
		res.Data = ret
	}
	c.JSON(200, res)
}

// StatusHandler returns users status highlight
func StatusHandler(c *gin.Context) {
	sid := c.DefaultQuery("sid", "")

	res := models.GeneralResponse{}
	result := session.GetSummaryStatus(sid)
	// data, _ := json.Marshal(result)
	res.Data = result

	c.JSON(200, res)
}
