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
