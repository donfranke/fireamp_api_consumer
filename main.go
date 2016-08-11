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
	_ "runtime/debug"
)

//const C_INTERVAL = 864000
const C_INTERVAL = 300000

const C_LIMIT = 3

func main() {
	log.Print("Starting")
	clientID := flag.String("c","", "3rd Party API Client ID")
	APIKey := flag.String("a", "", "API Key")
	flag.Parse()

	for {
		callAPI(*clientID,*APIKey)
		time.Sleep(time.Second * C_INTERVAL)
	}

	//callAPI(*clientID,*APIKey)
	// 3. iterate struct array and push to splunk forwarder
	
}

func callAPI(clientID string, APIKey string) {
	// 1. get events from FireAMP API
	apiresult := getEvents(clientID,APIKey)
	//fmt.Println(apiresult)
	//os.Exit(0)
	// 2. pass results of API call to create struct array
	results := parseJSON(apiresult)
	// 3. write to log file for splunk forwarder to pick up
	pushToSplunk(results)
}

func parseJSON(jsondata []byte) []Result{
	// sample data
	r := Result{}
	var results []Result

	//jsondata := `{}`
	res := &FireAMP_Event{}
	err := json.Unmarshal([]byte(jsondata), res)
	if err != nil {
		//debug.PrintStack()
		log.Fatal(err)
	}
	//_ = r
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

		//fmt.Println(item)
		//fmt.Printf("%s|%s|%s|%s|%s|%s|%s|%s\n", item.Date, strconv.Itoa(item.ID), item.EventType, item.Computer.Hostname, item.Detection, item.File.Disposition,r.filename,r.file_Sha256 )
	}
	return results
}

func getEvents(clientid string, apikey string) []byte {
	//encoded := ""
	//decodeCredentials(encoded)
	// ========================================
	client := &http.Client{}

	// ISO8601 2015-10-01T00:00:00+00:00
	t := time.Now()
	//then := t.AddDate(0, 0, -14)  // use this for testing purposes only
	then := t.Add(time.Second * C_INTERVAL * -1)
	start_date := then.Format(time.RFC3339)

	//fmt.Println(start_date)
	//url := "https://api.amp.cisco.com/v1/computers?limit=1"
	//url := "https://api.amp.cisco.com/v1/events?limit=1"
	url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT) + "&start_date=" + start_date + "&event_type[]=1090519054"
	//url := "https://api.amp.cisco.com/v1/events?limit=8&start_date=" + start_date
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	//req.Header.Add("Authorization", encoded)
	
	//fmt.Println(req)
	req.SetBasicAuth(clientid,apikey) // this uses base 64 encoding, which doesn't currently work

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s\n\n", body)
	//os.Exit(0)

	return body
}

func pushToSplunk(r []Result) {

	f, err := os.OpenFile("/tmp/fireamp_events.log", os.O_APPEND|os.O_WRONLY,0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

    //n4, err := w.WriteString("buffered\n")


	//var content string
	for i:=0;i<len(r);i++ {
		rr := r[i]
		s := rr.timestamp + "," + strconv.Itoa(rr.id) + "," + rr.event_type + "," + rr.computer + "," + rr.detection + "," + rr.disposition + "," + rr.filename + "," + rr.file_Sha256
		w.WriteString(s + "\n")
		fmt.Println(i,s)
	}
	_ = err
	w.Flush()
	defer f.Close()

	// iterate result struct array
	//   build content string
	//   send content string to splunk forwarder

	// example
	// timestamp,detection_timestamp,threat_type,computer,detection,disposition
}

func decodeCredentials(inEncoded string) {
	// ====== TO TEST BASE 64 ENCODING ======
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
