package main

import (
	"fmt"
)

func doDeleteRequest() {
	values := map[string]string{"url": *deleteUrl, "username": *deleteName, "password": *deletePassword}
	resp := doPostRequest(values, "delete")

	requestContent := getRequestContentFromResponse(resp)
	result, err := decrypt([]byte(requestContent["responseText"].(string)))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}
