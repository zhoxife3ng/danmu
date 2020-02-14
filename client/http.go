package client

import (
	"bytes"
	"github.com/x554462/weuse/utils"
	"io"
	"io/ioutil"
	"net/http"
)

const HttpMethodGet = "GET"
const HttpMethodPost = "POST"

func HttpReq(method, url string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		b, err := utils.JsonEncodeByte(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
