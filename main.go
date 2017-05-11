package main

import (
	"fmt"
	"net/url"
	"net/http"
	"crypto/tls"
	"os"
	"encoding/json"
	"io/ioutil"
)

var configuration Configuration
type Configuration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
}

func init() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println("Konfigurations Lesefehler: "+err.Error())
	}
	jsonDecoder := json.NewDecoder(configFile)
	configuration = Configuration{}
	jsonDecoder.Decode(&configuration)
}

func main(){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	payload := url.Values{}
	//payload.Add("api_key", "myapikey")
	req, err := http.NewRequest("GET", "https://"+configuration.Server+configuration.Port+"?" + payload.Encode(), nil)
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
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
