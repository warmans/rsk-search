package models

type IncompleteTranscription struct {
	Publication string  `json:"publication"`
	Series      int32   `json:"series"`
	Episode     int32   `json:"episode"`
	Chunks      []Chunk `json:"chunks"`
}

type Chunk struct {
	Raw         string `json:"raw"`
	StartSecond int64  `json:"start_second"`
	EndSecond   int64  `json:"end_second"`
}
