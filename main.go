package main

import (
	"fmt"
	"net/url"
	"net/http"
	"crypto/tls"
	"os"
	"encoding/json"
	"io/ioutil"

	clipboard "github.com/atotto/clipboard"
	"gopkg.in/alecthomas/kingpin.v2"
	"bytes"
)

var (
	app      = kingpin.New("safeClient", "A command-line tool to safe new passwords and usernames.")

	save     = app.Command("save", "save a new password")
	saveUrl = save.Arg("url", "the Url where the password is used").Required().String()
	saveName = save.Arg("name", "username/login").Required().String()
	savePassword = save.Arg("password", "password for the site").Required().String()

	get        = app.Command("get", "Get username and password for an Url.")
	getUrl   = get.Arg("url", "url you want the credentials for").Required().String()
)


var configuration Configuration
var client *http.Client
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
}

func doGetRequest(){

	payload := url.Values{}
	payload.Add("url", *getUrl)
	req, err := http.NewRequest("GET", "https://"+configuration.Server+configuration.Port+"/get?" + payload.Encode(), nil)
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

	err = clipboard.WriteAll(bodyString)
	if(err != nil){
		fmt.Println("Can't write to clipboard: "+ err.Error())
	}
}

func doSaveRequest(){
	url := "https://"+configuration.Server+configuration.Port+"/save"
	fmt.Println("URL:>", url)

	values := map[string]string{"url": *saveUrl,"username": *saveName, "password": *savePassword}
	jsonValue, _ := json.Marshal(values)
	jsonStr := bytes.NewBuffer(jsonValue)
	req, err := http.NewRequest("POST", url, jsonStr)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(configuration.User, configuration.Password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't make save request: ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main(){
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		case save.FullCommand():
			doSaveRequest()
		case get.FullCommand():
			doGetRequest()
	}
}
