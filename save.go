package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func doSaveRequest() {
	url := "https://" + configuration.Server + configuration.Port + "/save"
	values := map[string]string{"url": *saveUrl, "username": *saveName, "password": *savePassword}
	jsonValue, _ := json.Marshal(values)
	jsonStr := bytes.NewBuffer(jsonValue)
	req, err := http.NewRequest("POST", url, jsonStr)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(configuration.User, configuration.Password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't make save request: ", err)
	}
	defer resp.Body.Close()

	//body, _ := ioutil.ReadAll(resp.Body)
	requestContent := getRequestContentFromResponse(resp)
	result, err := decrypt([]byte(requestContent["responseText"].(string)))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
