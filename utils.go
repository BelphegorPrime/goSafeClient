package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func doPostRequest(values map[string]string, paramter string) *http.Response {
	url := "https://" + configuration.Server + configuration.Port + "/" + paramter
	jsonValue, _ := json.Marshal(values)
	jsonStr := bytes.NewBuffer(jsonValue)

	req, err := http.NewRequest("POST", url, jsonStr)
	if err != nil {
		fmt.Println("Can't build request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(configuration.User, configuration.Password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't execute request: " + err.Error())
	}

	return resp
}
