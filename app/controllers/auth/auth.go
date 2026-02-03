package auth

import (
	"api_kino/app/provider"
	"api_kino/config/api"
	"api_kino/service/web"
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Param struct {
	Token    string `json:"token"`
	ClientID string `json:"client_id"`
}

type TokenStruct struct {
	Token string `json:"access_token"`
}

func Authenticate(c *gin.Context, clientID string) (string, error) {
	var tokenData TokenStruct

	var payload map[string]string

	if clientID == "26" {
		payload = map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": os.Getenv("PARAM_CLIENTSECRET"),
			"scope":         "*",
		}
	} else if clientID == "35" {
		payload = map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": os.Getenv("PARAM_CLIENTSECRET2"),
			"scope":         "*",
		}
	} else if clientID == "41" {
		payload = map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": os.Getenv("PARAM_CLIENTSECRET3"),
			"scope":         "*",
		}
	} else if clientID == "42" {
		payload = map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": os.Getenv("PARAM_CLIENTSECRET4"),
			"scope":         "*",
		}
	} else if clientID == "93" {
		payload = map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": os.Getenv("PARAM_CLIENTSECRET5"),
			"scope":         "*",
		}
	}

	url := os.Getenv("URL_AUTH")

	resp, err := api.Client().R().
		EnableTrace().
		SetBody(payload).
		Post(url)

	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return "", err
	}

	respStr := resp.String()
	if len(respStr) > 1000 {
		respStr = respStr[:1000]
	}

	cleanResp := strings.ReplaceAll(respStr, "'", "")

	if resp.StatusCode() != http.StatusOK {
		query := fmt.Sprintf("INSERT INTO TransactionLog(Grup, SubGrup, StartTime, EndTime, Status, Message, RequestStr, RespondStr) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			"T002-KINO."+clientID, "AUTH", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), "E", "Access Denied!", payload, cleanResp)
		provider.QueryRun(query)
	}

	bodyBytes := resp.Body()

	err = json.Unmarshal(bodyBytes, &tokenData)
	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return "", err
	}

	return tokenData.Token, nil
}

func DraftSOKino(p Param) (*resty.Response, error) {
	var url string

	url = os.Getenv("URL_DRAFTSO") + "?SHIPMENTNO=" + os.Getenv("PARAM_SHIPMENTNO") + "&INTERFACEID=" +
		os.Getenv("PARAM_INTERFACEID") + "&CLIENTID=" + p.ClientID + "&"

	return api.Client().R().
		SetHeader("Authorization", "Bearer "+p.Token).
		SetHeader("Content-Type", "application/json").
		EnableTrace().
		Get(url)
}
