package models

const EndSecondEOF = -1

type Tscript struct {
	Publication string  `json:"publication"`
	Series      int32   `json:"series"`
	Episode     int32   `json:"episode"`
	Chunks      []Chunk `json:"chunks"`
}

func (i Tscript) ID() string {
	return IncompleteTranscriptID(i)
}

type Chunk struct {
	ID          string `json:"id"`
	Raw         string `json:"raw"`
	StartSecond int64  `json:"start_second"`
	EndSecond   int64  `json:"end_second"`
}
