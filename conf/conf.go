package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	WLAN      string = "wlan"
	LTE       string = "lte"
	TEST      string = "test"
	STREAMING string = "streaming"
	LATENCY   string = "latency"
	BANDWIDTH string = "bandwidth"
	NORMAL    string = "normal"
)

type Configuration struct {
	Env               string
	HandoverThreshold int
	Database          struct {
		Host      string `json:"host"`
		Port      int    `json:"port"`
		DBName    string `json:"dbname"`
		TableName string `json:"tablename"`
		User      string `json:"user"`
		Password  string `json:"password"`
	}
	Enbs              []string
	mme               string
	WiFiCLientIPOctet int
}

func LoadConfiguration(env string) Configuration {
	file := fmt.Sprintf("src/openapi/conf/conf.%s.json", env)
	var configuration Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println("ERROR in LoadConfiguration()", err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configuration)

	return configuration
}

func GeteNBonPosition(n int) (result string) {
	myConf := LoadConfiguration("env")
	if len(myConf.Enbs) >= n {
		return myConf.Enbs[n]
	} else {
		log.Println("ERROR: Position [%v]in eNB-array unknown", n)
	}
	return ""
}
