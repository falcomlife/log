package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
)

func (p *HttpClient) Get(path string, result interface{}) error {
	body, err := getMetricByHttp(p.Protocol, p.Host, p.Port, path)
	if err != nil {
		return err
	}
	err = resultToMap(body, result)
	if err != nil {
		return err
	}
	return nil
}

func (p *HttpClient) Post(path string, msg string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", getUrl(p.Protocol, p.Host, p.Port, path), bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		klog.Warning(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Warning(err)
		return nil, err
	}
	return body, nil
}

func getUrl(protocol string, host string, port string, path string) string {
	var url string = ""
	if port == "" {
		url =
			protocol +
				"://" +
				host +
				path
	} else {
		url =
			protocol +
				"://" +
				host +
				":" +
				port +
				path
	}
	return url
}

func getMetricByHttp(protocol string, host string, port string, path string) ([]byte, error) {
	resp, err := http.Get(getUrl(protocol, host, port, path))
	if err != nil {
		klog.Error("get request by http fail, message is " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Error(err)
	}
	return body, err
}

func resultToMap(b []byte, result interface{}) error {
	err := json.Unmarshal(b, result)
	if err != nil {
		klog.Warning(err)
	}
	return nil
}
