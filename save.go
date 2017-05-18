package main

import (
	"fmt"
)

func doSaveRequest() {
	values := map[string]string{"url": *saveUrl, "username": *saveName, "password": *savePassword}
	resp := doPostRequest(values, "save")

	requestContent := getRequestContentFromResponse(resp)
	result, err := decrypt([]byte(requestContent["responseText"].(string)))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
