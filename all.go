package main

import (
"fmt"
)

func doAllRequest() {
	resp := doPostRequest(map[string]string{}, "all")

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
	fmt.Println(string(result))
}

