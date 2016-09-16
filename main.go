// pulls event data from Cisco FireAMP cloud API
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"flag"
	"os"
	"time"
	"strconv"
	"bufio"

)

const C_INTERVAL = 600 // every 10 minutes
const C_LIMIT = 10

func main() {
	log.Print("Starting")
	ClientID := flag.String("clientid","", "3rd Party API Client ID")
	APIKey := flag.String("apikey", "", "API Key")
	LogFilePath := flag.String("log", "", "Log File Path")
	var i int
	
	flag.IntVar(&i, "interval", 60, "Interval between automatic runs")
	flag.Parse()

	// stay resident, call at interval
	for {
		//strPointer := Interval
		///ivl := strconv.ParseInt(strPointer,10,32)
		callAPI(*ClientID,*APIKey,*LogFilePath)
		time.Sleep(time.Second * C_INTERVAL)
	}	
}

func callAPI(clientID string, APIKey string, LogFilePath string) {
	// 1. get events from FireAMP API
	apiresult := getEvents(clientID,APIKey)
	// 2. pass results of API call to create struct array
	results := parseJSON(apiresult)
	// 3. write to log file for splunk forwarder to pick up
	//fmt.Println(results)
	pushToSplunk(results,LogFilePath)
}

func parseJSON(jsondata []byte) []Result{
	// sample data
	r := Result{}
	var results []Result
	res := &FireAMP_Event{}
	err := json.Unmarshal([]byte(jsondata), res)
	if err != nil {
		log.Fatal(err)
	}
	a1 := res.Data
	for _, item := range a1 {
		r.timestamp = item.Date
		r.id = item.ID
		r.event_type = item.EventType
		r.computer = item.Computer.Hostname
		r.computer = strings.ToLower(r.computer)
		r.detection = item.Detection
		r.disposition = item.File.Disposition
		r.filename = item.File.FileName
		r.file_Sha256 = item.File.Identity.Sha256
		results = append(results, r)
	}
	return results
}

func getEvents(clientid string, apikey string) []byte {
	client := &http.Client{}

	// ISO8601 2015-10-01T00:00:00+00:00
	t := time.Now()
	then := t.Add(time.Second * C_INTERVAL * -1)
	start_date := then.Format(time.RFC3339)
	//url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT) + "&start_date=" + start_date + "&event_type[]=1090519054"
	url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT) + "&start_date=" + start_date
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	
	req.SetBasicAuth(clientid,apikey) // this uses base 64 encoding, which doesn't currently work

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func pushToSplunk(r []Result, LogFilePath string) {
	f, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY,0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// iterate results and write to log file
	for i:=0;i<len(r);i++ {
		rr := r[i]
		s := rr.timestamp + "," + strconv.Itoa(rr.id) + "," + rr.event_type + "," + rr.computer + "," + rr.detection + "," + rr.disposition + "," + rr.filename + "," + rr.file_Sha256
		w.WriteString(s + "\n")
		fmt.Println(i,s)
	}
	_ = err
	w.Flush()
	defer f.Close()
}

// ====== TO TEST BASE 64 ENCODING ======
func decodeCredentials(inEncoded string) {
	decoded, err := base64.StdEncoding.DecodeString(inEncoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	s := strings.Split(string(decoded), ":")
	u := s[0]
	p := s[1]
	fmt.Printf("client id: %s, api credential: %s\n", u, p)
}
