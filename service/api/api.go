package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"api_kino/config/constant"
)

type Request struct {
	BaseUrl string
	Params  interface{}
	Header  http.Header
}

func ShouldBindJson(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(&target)
}

func ShouldBindXml(resp *http.Response, target interface{}) error {
	//byteValue, _ := ioutil.ReadAll(resp.Body)
	//return xml.Unmarshal(byteValue, &target)
	return xml.NewDecoder(resp.Body).Decode(&target)
}

func PutBind(route string, req Request, target interface{}) error {
	res, err := Send(http.MethodPut, route, req)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
		}
		return err
	}
	err = ShouldBindJson(res, target)
	if res != nil {
		defer res.Body.Close()
	}
	return err
}

func GetBind(route string, req Request, target interface{}) error {
	res, err := Send(http.MethodGet, route, req)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
		}
		return err
	}
	err = ShouldBindJson(res, target)
	if res != nil {
		defer res.Body.Close()
	}
	return err
}

func PostBind(route string, req Request, target interface{}) error {
	res, err := Send(http.MethodPost, route, req)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
		}
		return err
	}
	err = ShouldBindJson(res, target)
	if res != nil {
		defer res.Body.Close()
	}
	return err
}

func Send(method string, route string, req Request) (*http.Response, error) {
	var body io.Reader
	if req.Params != nil {
		jsonBytes, err := json.Marshal(req.Params)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonBytes)
		//log.Printf("Params: %s", jsonBytes)
	}
	url := constant.BaseUrl
	if req.BaseUrl != "" {
		url = req.BaseUrl
	}
	url += route
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if req.Header != nil {
		request.Header = req.Header
	} else {
		request.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	//log.Printf("Header: %s", request.Header)
	//log.Printf("Url: %s", url)
	client := &http.Client{}
	return client.Do(request)
}

func GetBindXml(route string, req Request, target interface{}) error {
	res, err := SendXml(http.MethodGet, route, req)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
		}
		return err
	}
	err = ShouldBindXml(res, target)
	if res != nil {
		defer res.Body.Close()
	}
	return err
}

func SendXml(method string, route string, req Request) (*http.Response, error) {
	var body io.Reader
	if req.Params != nil {
		jsonBytes, err := json.Marshal(req.Params)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonBytes)
		//log.Printf("Params: %s", jsonBytes)
	}
	url := constant.BaseUrl
	if req.BaseUrl != "" {
		url = req.BaseUrl
	}
	url += route
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if req.Header != nil {
		request.Header = req.Header
	} else {
		request.Header.Add("Content-Type", "application/xml;charset=utf-8")
	}
	//log.Printf("Header: %s", request.Header)
	//log.Printf("Url: %s", url)
	client := &http.Client{}
	return client.Do(request)
}
