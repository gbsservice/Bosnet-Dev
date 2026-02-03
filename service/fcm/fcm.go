package fcm

type Params struct {
	Data  interface{} `json:"data"`
	Topic string      `json:"to"`
}

type Response struct {
	MessageId int64 `json:"message_id"`
}

//func SendFcmByAuth(topic string, data interface{}, auth *model_base.User) (*Response, error) {
//	if strings.Contains(topic, "topics/") {
//		topic = "/" + topic
//		subOrganizationID := *auth.SubOrganizationID
//		if subOrganizationID == "9193fbd4-6ab2-11eb-b4e8-00155d585831" {
//			topic = topic + "-infoter"
//		}
//	}
//	header := http.Header{}
//	header.Add("Content-Type", "application/json")
//	header.Add("Authorization", "key="+app.Config().FcmKey)
//	fcmResponse := Response{}
//	if err := api.PostBind("send",
//		api.Request{
//			BaseUrl: constant.BaseUrlFcm,
//			Header:  header,
//			Params: Params{
//				Data:  data,
//				Topic: topic,
//			},
//		},
//		&fcmResponse,
//	); err != nil {
//		return nil, err
//	}
//	return &fcmResponse, nil
//}
//
//func Send(topic string, data interface{}, c *gin.Context) (*Response, error) {
//	auth := c.MustGet(constant.JwtClaim).(*jwt_auth.Claim)
//	if strings.Contains(topic, "topics/") {
//		topic = "/" + topic
//		subOrganizationID := *auth.SubOrganizationID
//		if subOrganizationID == "9193fbd4-6ab2-11eb-b4e8-00155d585831" {
//			topic = topic + "-infoter"
//		}
//	}
//	header := http.Header{}
//	header.Add("Content-Type", "application/json")
//	header.Add("Authorization", "key="+app.Config().FcmKey)
//	fcmResponse := Response{}
//	if err := api.PostBind("send",
//		api.Request{
//			BaseUrl: constant.BaseUrlFcm,
//			Header:  header,
//			Params: Params{
//				Data:  data,
//				Topic: topic,
//			},
//		},
//		&fcmResponse,
//	); err != nil {
//		return nil, err
//	}
//	return &fcmResponse, nil
//}
