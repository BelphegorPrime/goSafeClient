package main

import (
	"fmt"
	clipboard "github.com/atotto/clipboard"
)

func doGetRequest() {
	cipherText, err := encrypt([]byte(*getUrl))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	values := map[string]string{"url": string(cipherText)}

	resp := doPostRequest(values, "get")

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
