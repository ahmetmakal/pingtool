package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type Data struct {
	IpShow  bool   `json:"ip_show"`
	IpList  string `json:"ip_list"`
	TabSize uint   `json:"tab_size"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	configFile := "./ping.config.json"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		fmt.Println(configFile + " dosyasini duzenleyebilirsiniz")
		d1 := []byte("{\n    \"ip_show\": true,\n    \"tab_size\": 1,\n    \"ip_list\": \"8.8.8.8,77.88.8.8,208.67.222.222\"\n}")
		err := os.WriteFile(configFile, d1, 0644)
		check(err)
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	ips := strings.Split(payload.IpList, ",")

	for {
		fmt.Println(pingAt(ips, payload.IpShow, payload.TabSize))
		time.Sleep(1 * time.Second)
	}
}

func pingAt(ipAdresi []string, ipShow bool, tabSize uint) []string {

	tab := ""
	switch {
	case tabSize == 0:
		tab = ""
	case tabSize == 1:
		tab = "\t"
	case tabSize == 2:
		tab = "\t\t"
	case tabSize == 3:
		tab = "\t\t\t"
	}

	var my_slice []string
	for _, v := range ipAdresi {

		shell := "ping " + v + " -c 1 -i 1 -t 1 | grep icmp_seq | awk '{print $7}' | cut -d= -f2 | cut -d. -f1 | tr -d '\n'"
		shellOut, err := exec.Command("sh", "-c", shell).Output()
		if err != nil {
			fmt.Println("error shell")
			log.Fatal(err)
		}

		if string(shellOut) == "" {
			if ipShow {
				my_slice = append(my_slice, v+": "+Red+"err"+Reset+tab)
			} else {
				my_slice = append(my_slice, Red+"err"+Reset+tab)
			}
		} else {
			if ipShow {
				my_slice = append(my_slice, v+": "+Green+string(shellOut)+Reset+tab)
			} else {
				my_slice = append(my_slice, Green+string(shellOut)+Reset+tab)
			}
		}
	}

	return my_slice
}
