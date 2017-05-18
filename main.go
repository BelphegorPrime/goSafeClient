package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
)

var (
	app = kingpin.New("safeClient", "A command-line tool to safe new passwords and usernames.")

	save         = app.Command("save", "save a new password")
	saveUrl      = save.Arg("url", "the Url where the password is used").Required().String()
	saveName     = save.Arg("name", "username/login").Required().String()
	savePassword = save.Arg("password", "password for the site").Required().String()

	get    = app.Command("get", "Get username and password for an Url.")
	getUrl = get.Arg("url", "url you want the credentials for").Required().String()

	delete 		= app.Command("delete", "Delete an entry")
	deleteUrl 	= delete.Arg("url", "url you want to delete from").Required().String()
	deleteName 	= delete.Arg("name", "username/login").Required().String()
	deletePassword 	= delete.Arg("password", "password for the site").Required().String()
)

var configuration Configuration
var client *http.Client
var key []byte

type Configuration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	Key      string `json:"key"`
}

func init() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println("Konfigurations Lesefehler: " + err.Error())
	}
	jsonDecoder := json.NewDecoder(configFile)
	configuration = Configuration{}
	jsonDecoder.Decode(&configuration)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	key = []byte(configuration.Key) // 32 bytes
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case save.FullCommand():
		doSaveRequest()
	case get.FullCommand():
		doGetRequest()
	case delete.FullCommand():
		doDeleteRequest()
	}
}
