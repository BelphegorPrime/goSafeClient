package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	clipboard "github.com/atotto/clipboard"
	"net/http"
)

func doGetRequest() {
	url := "https://" + configuration.Server + configuration.Port + "/get"

	cipherText, err := encrypt([]byte(*getUrl))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	values := map[string]interface{}{"url": string(cipherText)}
	jsonValue, _ := json.Marshal(values)
	jsonStr := bytes.NewBuffer(jsonValue)

	req, err := http.NewRequest("POST", url, jsonStr)
	if err != nil {
		fmt.Println("Can't build request: " + err.Error())
	}
	req.SetBasicAuth(configuration.User, configuration.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't execute request: " + err.Error())
	}

	requestContent := getRequestContentFromResponse(resp)
	returnLines := requestContent["responseText"].([]interface{})
	var result string = ""
	for i := 0; i < len(returnLines); i++ {
		resultText, err := decrypt([]byte(returnLines[i].(string)))
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
		result = result + string(resultText)
	}

	bodyString := string(result)

	fmt.Println(bodyString)
	err = clipboard.WriteAll(bodyString)
	if err != nil {
		fmt.Println("Can't write to clipboard: " + err.Error())
	}
}
