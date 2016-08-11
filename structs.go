package main

type Result struct {
	timestamp       string
	id              int
	event_type      string
	computer        string
	detection       string
	disposition     string
	filename        string
	file_Sha256     string
}

type FireAMP_Event struct {
	Version string `json:"version"`
	Metadata struct {
		Links struct {
			Self string `json:"self"`
			Next string `json:"next"`
		} `json:"links"`
		Results struct {
			CurrentItemCount int `json:"current_item_count"`
			Index            int `json:"index"`
			ItemsPerPage     int `json:"items_per_page"`
			Total            int `json:"total"`
		} `json:"results"`
	} `json:"metadata"`
	Data []struct {
		ID                   int      `json:"id"`
		Timestamp            int      `json:"timestamp"`
		TimestampNanoseconds int      `json:"timestamp_nanoseconds"`
		Date        string `json:"date"`
		EventType   string `json:"event_type"`
		EventTypeID int    `json:"event_type_id"`
		Detection   string `json:"detection"`
		//DetectionID string `json:"detection_id"`  # some records have int, some have string
		GroupGuids           []string `json:"group_guids"`
		Computer struct {
			ConnectorGUID string `json:"connector_guid"`
			Hostname      string `json:"hostname"`
			User          string `json:"user"`
			Active        bool   `json:"active"`
			Links         struct {
				Computer   string `json:"computer"`
				Trajectory string `json:"trajectory"`
				Group      string `json:"group"`
			} `json:"links"`
		} `json:"computer"`
		File        struct {
			Disposition string `json:"disposition"`
			Identity    struct {
				Sha256 string `json:"sha256"`
				Sha1   string `json:"sha1"`
				MD5    string `json:"md5"`
			} `json:"identity"`
			FileName    string `json:"file_name"`
			FilePath    string `json:"file_path"`
		} `json:"file"`
		
	} `json:"data"`

}
