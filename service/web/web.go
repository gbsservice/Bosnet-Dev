package web

import (
	"bytes"
	"api_kino/config/constant"
	"time"

	"github.com/gin-gonic/gin"
)

type Broadcast struct {
	Self     bool
	Schedule bool
}

type H struct {
	Broadcast Broadcast
	Data      interface{}
	Records   int64
	Pages     int64
	Error     string
}

const (
	TypeImage = "image/png"
	TypePDF   = "application/pdf"
)

type HData struct {
	Broadcast   Broadcast
	Data        []byte
	Headers     map[string]string
	ContentType string
	Error       string
}

func Response(c *gin.Context, statusCode int, response H) {
	start := c.MustGet(constant.RequestTime).(time.Time)
	if statusCode > 300 {
		rsError := response.Error
		if rsError == "" {
			rsError = "something wrong happen"
		}
		c.AbortWithStatusJSON(200, gin.H{
			"process_time": time.Since(start).Seconds(),
			"success":      false,
			"error":        rsError,
			"extras":       response.Data,
		})
		return
	}
	if response.Broadcast.Self || response.Broadcast.Schedule {
		_, exists := c.Get(constant.Auth)
		_, exists = c.Get(constant.Activity)
		if exists {
			//auth := c.MustGet(constant.Auth).(*model_base.Users)
			//activity := c.MustGet(constant.Activity).(model_base.UsersActivities)
			//channel := ws.Channel{
			//	ID:    auth.ID,
			//	Group: ws.ChannelUser,
			//	Gin:   c,
			//}
			//if response.Broadcast.Schedule {
			//	channel = ws.Channel{
			//		ID:    *auth.ScheduleID,
			//		Group: ws.ChannelSchedule,
			//		Gin:   c,
			//	}
			//}
			//ws.BroadCastMessage(&ws.Response{
			//	Topic:  ws.TopicCreatedNotification,
			//	Option: ws.Option{Code: 200},
			//	Data:   activity,
			//}, &channel)
		}
	}
	c.JSON(200, gin.H{
		"process_time": time.Since(start).Seconds(),
		"success":      true,
		"data":         response.Data,
		"records":      response.Records,
		"pages":        response.Pages,
		"error":        nil,
	})
}

func ResponseData(c *gin.Context, statusCode int, response HData) {
	reader := bytes.NewReader(response.Data)
	contentLength := int64(len(response.Data))
	contentType := response.ContentType
	c.DataFromReader(statusCode, contentLength, contentType, reader, response.Headers)
}
