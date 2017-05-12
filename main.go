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
	"crypto/aes"
	"io"
	"crypto/cipher"
	"errors"
	"crypto/rand"
	"encoding/base64"
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
var key []byte

type Configuration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	Key 	 string `json:"key"`
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
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
	key = []byte(configuration.Key) // 32 bytes
}

func doGetRequest(){

	payload := url.Values{}

	ciphertext, err := encrypt(key, []byte(*getUrl))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	payload.Add("url", string(ciphertext))
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

	result, err := decrypt(key, bodyBytes)
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

func doSaveRequest(){
	url := "https://"+configuration.Server+configuration.Port+"/save"
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

	body, _ := ioutil.ReadAll(resp.Body)
	result, err := decrypt(key, body)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Println(string(result))
}

func main(){
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		case save.FullCommand():
			doSaveRequest()
		case get.FullCommand():
			doGetRequest()
	}
}
