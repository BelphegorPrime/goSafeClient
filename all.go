package main

import (
	"fmt"
	"github.com/BelphegorPrime/lib"
)

func doAllRequest() {
	resp := doPostRequest(map[string]interface{}{"crypto": 1}, "all")

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
	fmt.Println(string(result))
}

