/*
 * @Date: 2022-02-22 14:48:09
 * @LastEditors: recar
 * @LastEditTime: 2022-02-24 16:09:22
 */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
)

var (
	flagBarkApi string
	flagSctApi  string
	configPah   = "/root/.toBtg/conf.yaml"
)

type Config struct {
	Barkapi string `yaml:"barkapi"`
	Sctapi  string `yaml:"sctapi"`
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func loadApi() []string {
	var barkApi string
	var sctApi string
	apiList := []string{}
	barkApi = os.Getenv("BARK_API")
	sctApi = os.Getenv("SCT_API")
	if barkApi == "" && sctApi == "" {
		// 读conf.yml
		isExists, err := PathExists(configPah)
		if err != nil {
			log.Println(err.Error())
		}
		if isExists {
			var conf Config
			yamlFile, err := ioutil.ReadFile(configPah)
			if err != nil {
				log.Println(err.Error())
			}
			err1 := yaml.Unmarshal(yamlFile, &conf)
			if err1 != nil {
				log.Println(err1.Error())
			}
			barkApi = conf.Barkapi
			sctApi = conf.Sctapi
		} else {
			log.Println("file not exists")
			barkApi = flagBarkApi
			sctApi = flagSctApi
		}

	}

	if barkApi != "" {
		log.Println("use barkA")
		apiList = append(apiList, barkApi)
	}
	if sctApi != "" {
		log.Println("use sct")
		apiList = append(apiList, sctApi)
	}
	return apiList
}

func send(api string, text string) {
	url := fmt.Sprintf("%s%s", api, text)
	_, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
}

func init() {
	flag.StringVar(&flagBarkApi, "b", "", "Bark api")
	flag.StringVar(&flagSctApi, "s", "", "Sct api")
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	apiList := loadApi()
	if len(apiList) == 0 {
		log.Println("需要配置api")
		return
	}
	if !terminal.IsTerminal(0) {
		// 从管道
		b, err := ioutil.ReadAll(os.Stdin)
		if err == nil {
			text := string(b)
			text = strings.Replace(text, "\n", "%0a", -1)
			for _, api := range apiList {
				send(api, text)
			}
			log.Println("toBtg send successfully")
		} else {
			log.Println("toBtg send error")
		}
	}

}
