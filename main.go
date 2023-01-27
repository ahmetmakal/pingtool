package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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

var detectOS = runtime.GOOS

type Data struct {
	IpShow  bool   `json:"ip_show"`
	IpList  string `json:"ip_list"`
	TabSize uint   `json:"tab_size"`
	Dot     bool   `json:"dot"`
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

	content, err := os.ReadFile(configFile)
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

	tab := ""
	switch {
	case payload.TabSize == 0:
		tab = " "
	case payload.TabSize == 1:
		tab = "\t"
	case payload.TabSize == 2:
		tab = "\t\t"
	case payload.TabSize == 3:
		tab = "\t\t\t"
	}

	for {
		fs := pingAt(ips, payload.IpShow, payload.Dot)
		for i, v := range fs {
			if i > 0 {
				fmt.Print(tab)
			}
			fmt.Print(v)
		}
		fmt.Print("\n")
		time.Sleep(1 * time.Second)
	}
}

func pingAt(ipAdresi []string, ipShow bool, dot bool) []string {

	var my_slice []string
	for _, v := range ipAdresi {

		var shell string

		switch detectOS {
		case "darwin":
			if dot {
				shell = "ping " + v + " -c 1 -i 1 -t 1 | grep icmp_seq | awk '{print $7}' | cut -d= -f2 | tr -d '\n'"
			} else {
				shell = "ping " + v + " -c 1 -i 1 -t 1 | grep icmp_seq | awk '{print $7}' | cut -d= -f2 | cut -d. -f1 | tr -d '\n'"
			}
			shellOut, err := exec.Command("sh", "-c", shell).Output()
			if err != nil {
				fmt.Println("error shell")
				log.Fatal(err)
			}

			if string(shellOut) == "" {
				if ipShow {
					my_slice = append(my_slice, v+": "+Red+"err"+Reset)
				} else {
					my_slice = append(my_slice, Red+"err"+Reset)
				}
			} else {
				if ipShow {
					my_slice = append(my_slice, v+": "+Green+string(shellOut)+Reset)
				} else {
					my_slice = append(my_slice, Green+string(shellOut)+Reset)
				}
			}
		case "windows":
			Reset := "\x1b[0m"
			Red = "\x1b[31m"
			Green = "\x1b[32m"
			shellOut, _ := exec.Command("ping", "-n", "1", v).Output()
			convShellOut := string(shellOut[:])
			check := "Reply from"
			var outputPing string
			if strings.Contains(convShellOut, check) {
				cutTime := regexp.MustCompile(`time=(\d+ms)`)
				time := cutTime.FindAllString(convShellOut, 1)
				cutDigit := strings.Join(time, " ")
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
				digit := re.FindStringSubmatch(cutDigit)
				outputPing = strings.Join(digit, " ")
			} else {
				fmt.Print("Ping error.")
			}
			if outputPing == "" {
				if ipShow {
					my_slice = append(my_slice, v+": "+Red+"err"+Reset)
				} else {
					my_slice = append(my_slice, Red+"err"+Reset)
				}
			} else {
				if ipShow {
					my_slice = append(my_slice, v+": "+Green+outputPing+Reset)
				} else {
					my_slice = append(my_slice, Green+outputPing+Reset)
				}
			}
		}
	}

	return my_slice
}
