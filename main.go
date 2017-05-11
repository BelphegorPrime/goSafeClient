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

func doGetRequest(){
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


	err = clipboard.WriteAll(bodyString)
	if(err != nil){
		fmt.Println("Can't write to clipboard: "+ err.Error())
	}
	fmt.Println(clipboard.ReadAll())
}

func main(){
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		// Register user
		case save.FullCommand():
			fmt.Println(*saveUrl)
			fmt.Println(*saveName)
			fmt.Println(*savePassword)
		// Post message
		case get.FullCommand():
			fmt.Println(*getUrl)
			doGetRequest()
	}
}
