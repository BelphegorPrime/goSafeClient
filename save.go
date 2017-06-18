package main

import (
	"fmt"
	"github.com/BelphegorPrime/lib"
)

func doSaveRequest() {
	values := map[string]string{"url": *saveUrl, "username": *saveName, "password": *savePassword}
	resp := doPostRequest(values, "save")

	requestContent := lib.GetRequestContentFromResponse(resp)
	result, err := lib.Decrypt([]byte(requestContent["responseText"].(string)), key)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
