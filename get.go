package main

import (
	"fmt"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	clipboard "github.com/atotto/clipboard"
)

func doGetRequest(){
	url := "https://"+configuration.Server+configuration.Port+"/get"

	ciphertext, err := encrypt([]byte(*getUrl))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	values := map[string]interface{}{"url": string(ciphertext)}
	jsonValue, _ := json.Marshal(values)
	jsonStr := bytes.NewBuffer(jsonValue)

	req, err := http.NewRequest("POST", url, jsonStr)
	if(err != nil){
		fmt.Println("Can't build request: "+ err.Error())
	}
	req.SetBasicAuth(configuration.User, configuration.Password)
	resp, err := client.Do(req)
	if(err != nil){
		fmt.Println("Can't execute request: "+ err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if(err != nil){
		fmt.Println("Can't read response body: "+ err.Error())
	}

	result, err := decrypt(bodyBytes)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	bodyString := string(result)

	fmt.Println(bodyString)
	err = clipboard.WriteAll(bodyString)
	if(err != nil){
		fmt.Println("Can't write to clipboard: "+ err.Error())
	}
}
