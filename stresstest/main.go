package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/fatih/color"
	skynet "github.com/autisticvegan/go-skynet"

)

/*
 * This program does mass upload/download using sia skynet
 * Use it to test banning bad IPs, load balancing, rate-limiting, whatever
 *
 * Example usage
 * stress siasky.net u 69 proxylist.txt payload.txt

 * explanation:
 * stress - the exe
 * siasky.net - the portal to use
 * u - u for upload (use d for download)
 * 69 - the count of simultaneous connections to open per entry in proxylist
 * proxylist.txt - the list of proxies to use
* payload.txt - the file to upload


 * example for download:
 * stress siasky.net d 6 proxylist.txt VACKuQGhq6HN15CEmzRIXi5PDz9KGczdEFrnC_RcFWC4sg

 * explanation:
 * stress - the exe
 * siasky.net - the portal to use
 * d - d for download
 * 6 - the count of simultaneous connections to open per entry in proxylist
 * proxylist.txt - the list of proxies to use
 * VACKuQGhq6HN15CEmzRIXi5PDz9KGczdEFrnC_RcFWC4sg - the file (skylink) to download


 * If no proxies are supplied, only the current connection is used
 future improvements: yaml config file?
*/

func uploadWorker(wg *sync.WaitGroup, portal string, file string, proxy string) {
	defer wg.Done()
	client := skynet.NewCustom(portal, skynet.Options{})
	_, err := client.UploadFile(file, skynet.DefaultUploadOptions, proxy)
	if err != nil {
		color.HiRed("err in upload")
		fmt.Println(err)
	}
	//color.HiYellow("closing out uploadWorker")
}

func downloadWorker(wg *sync.WaitGroup, portal string, file string, proxy string) {
	defer wg.Done()
	client := skynet.NewCustom(portal, skynet.Options{})
	genPath := "tempdownload"
	err := client.DownloadFile(genPath, file, skynet.DefaultDownloadOptions, proxy)
	if err != nil {
		color.HiRed("err in download")
		fmt.Println(err)
	}
	//color.HiYellow("closing out downloadWorker")
}

func stress(portal string, isUpload bool, connCount int, proxies []string, file string) {
//for each proxie, open X connections to download each Y file

	var wg sync.WaitGroup
	for i, p := range proxies {
		if p == "" {
			color.HiGreen("Testing with proxy #" + strconv.Itoa(i) + " (default no proxy)")
		} else {
			color.HiGreen("Testing with proxy #" + strconv.Itoa(i) + " " + p)
		}

		color.HiRed("Opening " + strconv.Itoa(connCount) + " connections...")
		for j := 0; j < connCount; j++ {
			wg.Add(1)
			if isUpload {
				go uploadWorker(&wg, portal, file, p)
			} else {
				go downloadWorker(&wg, portal, file, p)
			}
		}
	}
	wg.Wait()
}

func parseListFromFile(filepath string) []string {
	res := []string{}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func main() {

	if len(os.Args) != 5 || len(os.Args) != 6 {
		fmt.Println("See comments for usage")
		return
	}

	portal := os.Args[1]
	connCount, _ := strconv.Atoi(os.Args[3])
	proxies := parseListFromFile(os.Args[4])
	file := os.Args[5]
	if os.Args[2] == "d" {
		color.HiCyan("Starting up download stressor...")
		stress(portal, false, connCount, proxies, file)
	} else if os.Args[2] == "u" {
		color.HiCyan("Starting up upload stressor...")
		stress(portal, true, connCount, proxies, file)
	}

/*
	debug shit:
	proxies := parseListFromFile("proxylist.txt")
	stress("http://siasky.net", true, 2, proxies, "upload.txt")
	fmt.Scanln()
	*/
}