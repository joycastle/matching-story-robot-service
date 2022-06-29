package lib

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func PostFormToJson(url string, data url.Values, output interface{}) error {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return err
	}

	return nil
}
