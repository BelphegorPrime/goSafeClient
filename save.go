package main

import (
	"fmt"
	"github.com/BelphegorPrime/lib"
)

func doSaveRequest() {
	values := map[string]interface{}{"url": *saveUrl, "username": *saveName, "password": *savePassword, "crypto": 1}
	resp := doPostRequest(values, "save")

	requestContent := lib.GetRequestContentFromResponse(resp)
	result, err := lib.Decrypt([]byte(requestContent["responseText"].(string)), key)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
