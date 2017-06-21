package main

import (
	"fmt"
	"github.com/BelphegorPrime/lib"
)

func doDeleteRequest() {
	values := map[string]interface{}{"url": *deleteUrl, "username": *deleteName, "password": *deletePassword, "crypto": 1}
	resp := doPostRequest(values, "delete")

	requestContent := lib.GetRequestContentFromResponse(resp)
	result, err := lib.Decrypt([]byte(requestContent["responseText"].(string)), key)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
