package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	lastRun = 0
	http.HandleFunc("/livecount", streamStats)
	if err := http.ListenAndServe(":8092", nil); err != nil {
		panic(err)
	}
}

var lastRun int64
var userCount int

func countLiveUsers() int {
	fmt.Println("Disk access")
	file, _ := os.Open("/var/log/nginx/access.log")
	scanner := bufio.NewScanner(file)
	var ips []string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), ".m3u8") {
			layout := "02/Jan/2006:15:04:05"
			t, _ := time.Parse(layout, strings.Split(scanner.Text(), " ")[3][1:len(strings.Split(scanner.Text(), " ")[3])])
			if (time.Now().Unix() - t.Unix()) < 300 {
				ip := strings.Split(scanner.Text(), " ")[0]
				if !stringInSlice(ip, ips) {
					ips = append(ips, ip)
				}
			}

		}
	}
	return len(ips)
}

func streamStats(w http.ResponseWriter, r *http.Request) {
	if (time.Now().Unix() - lastRun) > 10 {
		lastRun = time.Now().Unix()
		userCount = countLiveUsers()
	}

	message := strconv.Itoa(userCount)
	w.Write([]byte(message))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
