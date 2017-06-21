package main

import (
	"fmt"
	clipboard "github.com/atotto/clipboard"
	"github.com/BelphegorPrime/lib"
)

func doGetRequest() {
	cipherText, err := lib.Encrypt([]byte(*getUrl), key)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	values := map[string]interface{}{"url": string(cipherText), "crypto": 1}

	resp := doPostRequest(values, "get")

	requestContent := lib.GetRequestContentFromResponse(resp)
	returnLines := requestContent["responseText"].([]interface{})
	var result string = ""
	for i := 0; i < len(returnLines); i++ {
		resultText, err := lib.Decrypt([]byte(returnLines[i].(string)), key)
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
