package common

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
)

func (p *HttpClient) Get(path string) (map[string]interface{}, error) {
	body, err := getMetricByHttp(p.Host, p.Port, path)
	if err != nil {
		return nil, err
	}
	m, err := resultToMap(body)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getMetricByHttp(host string, port string, path string) ([]byte, error) {
	var url string = ""
	if port == "" {
		url =
			"http://" +
				host +
				path
	} else {
		url =
			"http://" +
				host +
				":" +
				port +
				path
	}
	resp, err := http.Get(url)
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

func resultToMap(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	json.Unmarshal(b, &m)
	return m, nil
}
