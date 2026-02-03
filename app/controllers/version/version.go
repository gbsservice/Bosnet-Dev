package version

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"api_kino/config/constant"
	"api_kino/service/web"
	"time"
)

type GetLatestRequest struct {
	Name string `form:"name"`
}

type ApiVersionResponse struct {
	ApiVersion string `json:"version"`
}

func GetApi(c *gin.Context) {
	data := ApiVersionResponse{
		ApiVersion: constant.ApiVersion,
	}
	web.Response(c, http.StatusOK, web.H{
		Data: data,
	})
}

func GetTime(c *gin.Context) {
	time := time.Now()
	web.Response(c, http.StatusOK, web.H{
		Data: time,
	})
}
