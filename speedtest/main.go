package main

import (
	"bufio"
	"io/ioutil"
	//"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	skynet "github.com/NebulousLabs/go-skynet"
	"github.com/PuerkitoBio/goquery"
	//	skydb_auvega "github.com/autisticvegan/skynetskydb"
	"github.com/fatih/color"

)

// PortalPingFileSize is used as the file size that will be used for testing speed of portals to client
type PortalPingFileSize string
const (
	Small PortalPingFileSize = "SMALL" // 1KB
	Medium PortalPingFileSize = "MEDIUM" // 10 MB
	Large PortalPingFileSize = "LARGE" //300 MB
	portalsLink string = "https://siasetup.info/tools/portals"
	txtFileForPortals string = "portals.txt"
)

type portalAndTime struct  {
	name string
	time int
}

// TODO(autisticvegan): should timestamps be sth else? Should we have the errors too?
type resultsObject struct {
	fastestPortal string
	slowestPortal string
	fastestTime string
	slowestTime string
	resultsTable map[string]int // contains mapping of portals to their median time, -69 if it had an error
}

/*
 * This function will open the siasetup page, and scrape the portals by looking at
 * the <td> elements with anchors in them.  This seems kind of loosy-goosey and may be a terrible
 * idea but lets go with it for now.
 */
 func scrapePageForPortals(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
			log.Fatal(err)
	}
	listOfPortals := []string{}

	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		s.Find("section").Each(func(j int, s2 *goquery.Selection) {
			s2.Find("table").Each(func(k int, s3 *goquery.Selection) {
				s3.Find("tbody").Each(func(k int, s4 *goquery.Selection) {
					s4.Find("tr").Each(func(k int, s5 *goquery.Selection) {
						s5.Find("td").Each(func(k int, s6 *goquery.Selection) {
							if link := s6.Find("a"); link != nil && link.Length() > 0 {
								if link.Length() == 1 {
									listOfPortals = append(listOfPortals, link.Text())
								} else {
									link.Each(func(k int, s7 *goquery.Selection) {
										listOfPortals = append(listOfPortals, s7.Text())
									})
								}
							}
						})
					})
				})
			})
		})
	})
	return listOfPortals
}

/*
 * This function is used as a fallback if the scraping of the siasetup page doesn't work
 * It will open "portals.txt" and get the list of portals from there, reading line by line
 */
 func getHardCodedList(path string) []string {
	file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	listOfPortals := []string{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        listOfPortals = append(listOfPortals, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
	return listOfPortals
}

/*
 * Function to upload, then download a file from a given portal
 * returns the milliseconds the whole process took
 */
func skynetPortalPing(port string, file string) (int, error) {
	t := time.Now()
	err := uploadAndDownload(file, port)
	if err != nil {
		return -1, err
	}
	return int(time.Since(t).Milliseconds()), nil
}

/*
 * Shuffle is used for putting the portals in a shuffled order
 * when we ping them to try to make the test more realistic
 */
 func shuffle(vals []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
	  n := len(vals)
	  randIndex := r.Intn(n)
	  vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	  vals = vals[:n-1]
	}
}

/*
 * Function to test portal speed.  It takes in a size, and runs over a list of portals, testing each one 3 times,
 * and taking the median of the times.  It returns a list of portals sorted by time, from fastest to slowest.
 * In case errors are encountered, that portal is put in the "err list" which will be the in the second returned value.
 */
 func portalPing(portalPingFileSize PortalPingFileSize) ([]portalAndTime, []string) {
	//either scrape, or use backup list
	color.HiYellow("Scraping portal list from " + portalsLink + " ...")
	listOfPortals := scrapePageForPortals(portalsLink)
	if len(listOfPortals) == 0 {
		color.HiYellow("Something went wrong with scraping " + portalsLink + " ... so, falling back to " + txtFileForPortals)
		listOfPortals = getHardCodedList("portals.txt")
	}
	shuffle(listOfPortals)
	color.HiBlue("Portals loaded:")
	for _, p := range listOfPortals {
		color.HiYellow(p)
	}
	portalToTime := make(map[string]int)
	timeToPortal := make(map[int]string)
	portalsThatErrored := []string{}
	for _, port := range listOfPortals {
		color.Cyan("========================================")
		color.HiGreen("Testing portal: " + port)
		timesTaken := []int{}
		//do the test 3 times, take the median
		for i := 0; i < 3; i++ {
			var err error
			var val int
			switch portalPingFileSize {
			case Small:
				val, err = skynetPortalPing(port, "sample_small.txt")
			case Medium:
				val, err = skynetPortalPing(port, "sample_medium.txt")
			case Large:
				val, err = skynetPortalPing(port, "sample_large.txt")
			}
			if err == nil {
				timesTaken = append(timesTaken, val)
			} else {
				color.Red("error encountered with portal " + port + " ... skipping this one for the test")
				//color.HiYellow(err.Error())
				portalsThatErrored = append(portalsThatErrored, port)
				break
			}
		}

		if len(timesTaken) == 0 || timesTaken[0] == -1 {
			continue
		}
		sort.Ints(timesTaken)
		for _, time := range timesTaken {
			color.White(strconv.Itoa(time) + " ms")
		}
		medianTime := timesTaken[len(timesTaken)/2]
		color.Cyan("Portal " + port + " took a median time of " + strconv.Itoa(medianTime) + " ms")
		portalToTime[port] = medianTime
		timeToPortal[medianTime] = port
	}
	times := []int{}
	for _, v := range portalToTime {
		times = append(times, v)
	}
	sort.Ints(times)

	portalsAndTimes := []portalAndTime{}
	for _, t := range times {
		portalsAndTimes = append(portalsAndTimes, portalAndTime{
			timeToPortal[t],
			t,
		})
	}
	return portalsAndTimes, portalsThatErrored
}

/*
 * Function basically from https://siasky.net/ that uploads and downloads a file given a portal
 */
func uploadAndDownload(path string, portal string) (error) {
	if len(path) == 0 || len(portal) == 0 {
		return nil
	}
	var client = skynet.NewCustom(portal, skynet.Options{})
	skylink, err := client.UploadFile(path, skynet.DefaultUploadOptions)
    if err != nil {
		fmt.Println("Error in upload " + err.Error())
		return err
    }
	err = client.DownloadFile(path + "_downloaded", skylink, skynet.DefaultDownloadOptions)
    if err != nil {
		fmt.Println("Error in download " + err.Error())
		return err
	}
	return nil
}

func uploadFile() {}

func encodeResultsObjToStr(resObj resultsObject) string {
	str := "fastest: " + resObj.fastestPortal + " at " + resObj.fastestTime + "\n"
	str += "slowest: " + resObj.slowestPortal + " at " + resObj.slowestTime + "\n"
	return str
}

func printResultsObj(resObj resultsObject) {
	color.HiMagenta("fastest was:")
	color.HiGreen(resObj.fastestPortal)
	color.HiMagenta("with time of " + resObj.fastestTime + " ms")
	color.HiMagenta("slowest was:")
	color.HiYellow(resObj.slowestPortal)
	color.HiMagenta("with time of " + resObj.slowestTime + " ms")
	//name  median time   comment
	color.HiCyan("Table of results:")
	for k, v := range resObj.resultsTable {
		time := strconv.Itoa(v)
		if v == -69 {
			color.HiRed(k + "\t\t\t" +  "\t" + "\t\t\t" + "ERR")
		} else {
			color.HiGreen(k + "\t\t\t" + time + " ms\t\t\t" + "OK")
		}
		
	}
}

func writeResultsToTxtFile(resObj resultsObject) {
	fileName := "test_results_" + strconv.FormatInt(time.Now().Unix(), 10)
	resStr := "Results:\n"
	for k,v := range resObj.resultsTable {
		time := strconv.Itoa(v)
		if v == -69 {
			time = "" 
		}
		resStr += k + "\t" + time + " ms\n"
	}
	ioutil.WriteFile(fileName, []byte(resStr), 0644)
}
func uploadResultsToSkyDB(resObj resultsObject, portal string) {
	color.HiCyan("Saving results to registry in SkyDB ...")

	/*
	Since this doesn't seem to be working as intended, use michiel_post 's C# implementation in a separate program
	dataKey := "SPEEDTEST_RESULTS"
	//Generate a new key if old one doesn't exist
    _, err := os.Stat("keys.txt")
    if os.IsNotExist(err) {
		skydb_auvega.WriteNewKeyPairToFile()
    }
	pub, priv, _ := skydb_auvega.ParseKeysFromFile("keys.txt")
	bytesToSend := []byte(encodeResultsObjToStr(resObj))
	err = skydb_auvega.SetBytesToRegistry(string(priv), string(pub), dataKey, bytesToSend, portal)
	fmt.Println(err)
	*/
}

func parseResults(portalsAndTimes []portalAndTime, errList []string) resultsObject {
	fastestPortalName := portalsAndTimes[0].name
	slowestPortalName := portalsAndTimes[len(portalsAndTimes) - 1].name
	fastestPortalTime := strconv.Itoa(portalsAndTimes[0].time)
	slowestPortalTime := strconv.Itoa(portalsAndTimes[len(portalsAndTimes) - 1].time)
	resTab := make(map[string]int)
	for _, p := range errList {
		resTab[p] = -69
	}
	for _, p := range portalsAndTimes {
		resTab[p.name] = p.time
	}
	return resultsObject{
		fastestPortal: fastestPortalName,
		slowestPortal: slowestPortalName,
		fastestTime: fastestPortalTime,
		slowestTime: slowestPortalTime,
		resultsTable: resTab,
	}
}

func main() {
	color.HiMagenta("Welcome to a tool for testing sia skynet portals :)")
	color.HiMagenta("")
	portalsAndTimes, errList := portalPing(Small)
	if len(portalsAndTimes) == 0 {
		panic("something is wrong")
	}
	color.HiMagenta("Test complete")
	color.HiMagenta("")
	resObj := parseResults(portalsAndTimes, errList)
	printResultsObj(resObj)
	writeResultsToTxtFile(resObj)
	uploadResultsToSkyDB(resObj, "https://www.skyportal.xyz")
	color.White("Press ENTER to exit...")
	fmt.Scanln()
}