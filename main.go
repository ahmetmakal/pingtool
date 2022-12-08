package main

import (
	"fmt"
	"log"
	"os/exec"
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

func main() {

	// content, err := ioutil.ReadFile("ping-ips.txt")
	// if err != nil {
	// 	fmt.Println("ping-ips.txt dosyası hatalı")
	// }
	// content = []byte(strings.Replace(string(content), "\n", "", -1))
	// lines := strings.Split(string(content), ",")

	address := []string{
		"8.8.8.8",
		"77.88.8.8",
		"1.1.1.1",
	}

	Ping := make([]string, len(address))
	for i, v := range address {
		// fmt.Println(i, v)
		if v != "\r\n" {
			Ping[i] = v
		}
	}

	fmt.Println(Ping)

	for {
		fmt.Println(pingAt(Ping))
		time.Sleep(1 * time.Second)
	}
}

func pingAt(ipAdresi []string) []string {

	var my_slice []string
	for _, v := range ipAdresi {

		shell := "ping " + v + " -c 1 -i 1 -t 1 | grep icmp_seq | awk '{print $7}' | cut -d= -f2 | cut -d. -f1 | tr -d '\n'"
		shellOut, err := exec.Command("sh", "-c", shell).Output()
		if err != nil {
			fmt.Println("error shell")
			log.Fatal(err)
		}

		if string(shellOut) == "" {
			my_slice = append(my_slice, Red+"err"+Reset)
		}

		my_slice = append(my_slice, Green+string(shellOut)+Reset)
	}

	return my_slice
}
