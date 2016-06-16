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
)

type Result struct {
	timestamp 			string
	event_type          string
	computer            string
	detection           string
	disposition         string
}

type FireAMP struct {
	Data []struct {
		Computer struct {
			Active        bool   `json:"active"`
			ConnectorGUID string `json:"connector_guid"`
			Hostname      string `json:"hostname"`
			Links         struct {
				Computer   string `json:"computer"`
				Group      string `json:"group"`
				Trajectory string `json:"trajectory"`
			} `json:"links"`
			User string `json:"user"`
		} `json:"computer"`
		Date        string `json:"date"`
		Detection   string `json:"detection"`
		DetectionID int    `json:"detection_id"`
		EventType   string `json:"event_type"`
		EventTypeID int    `json:"event_type_id"`
		File        struct {
			Disposition string `json:"disposition"`
			FileName    string `json:"file_name"`
			FilePath    string `json:"file_path"`
			Identity    struct {
				Sha256 string `json:"sha256"`
			} `json:"identity"`
			Parent struct {
				Disposition string `json:"disposition"`
				FileName    string `json:"file_name"`
				Identity    struct {
					Sha256 string `json:"sha256"`
				} `json:"identity"`
			} `json:"parent"`
		} `json:"file"`
		GroupGuids           []string `json:"group_guids"`
		ID                   int      `json:"id"`
		Timestamp            int      `json:"timestamp"`
		TimestampNanoseconds int      `json:"timestamp_nanoseconds"`
	} `json:"data"`
	Metadata struct {
		Links struct {
			Next string `json:"next"`
			Self string `json:"self"`
		} `json:"links"`
		Results struct {
			CurrentItemCount int `json:"current_item_count"`
			Index            int `json:"index"`
			ItemsPerPage     int `json:"items_per_page"`
			Total            int `json:"total"`
		} `json:"results"`
	} `json:"metadata"`
	Version string `json:"version"`
}

func main() {
	log.Print("Starting")
	// 1. get events from FireAMP API
	//getEvents()
	// 2. pass results of API call to create struct array
	parseJSON()
	// 3. iterate struct array and push to splunk forwarder
	//pushEventsToSplunk()
}

func parseJSON() {

	// sample data
	//r := Result{}
	//var results []Result

	jsondata := `{}`
	res := &FireAMP{}
	err := json.Unmarshal([]byte(jsondata), res)
	if err != nil {
		log.Fatal(err)
	}
	a1 := res.Data
	for _, item := range a1 {
		//r.timestamp = item.Date
		//r.event_type = item.EventType
		//r.computer = item.Computer.Hostname
		//r.detection = item.Detection
		//r.disposition = item.Disposition
		//results.append(r)
		
		fmt.Printf("%s|%s|%s|%s|%s\n", item.Date, item.EventType, item.Computer.Hostname, item.Detection, item.File.Disposition)
	}
}

func getEvents() {
	encoded := "NmM1MzE1M2M2MWUwNDMzODNiMTU6YTc3MjFhMjYtNTAxOS00MWU0LWFmZjItZWVhNzgxYzA2NjlmCg=="
	decodeCredentials(encoded)
	// ========================================
	client := &http.Client{}
	url := "https://api.amp.cisco.com/v1/computers"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", encoded)

	//fmt.Println(req)
	//req.SetBasicAuth(u,p) // this uses base 64 encoding, which doesn't currently work

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", body)
}

func pushToSplunk() {
	//var content string

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
