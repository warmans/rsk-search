package assemblyai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	APIBaseURL string = "https://api.assemblyai.com/v2"
)

type TranscribeRequest struct {
	AudioURL          string `json:"audio_url"`
	SpeakerLabels     bool   `json:"speaker_labels"`
	EntityDetection   bool   `json:"entity_detection"`
	SentimentAnalysis bool   `json:"sentiment_analysis"`
	Summarization     bool   `json:"summarization"`
	SummaryModel      string `json:"summary_model"`
	SummaryType       string `json:"summary_type"`
}

type Entity struct {
	EntityType string `json:"entity_type"`
	Text       string `json:"text"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
}

type Sentiment struct {
	Text       string  `json:"text"`
	Start      int64   `json:"start"`
	End        int64   `json:"end"`
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
}

type TranscriptWord struct {
	Text  string `json:"text"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	//Words  []*TranscriptWord `json:"words"` // don't care
}

type TranscriptUtterance struct {
	Pos        int64   `json:"pos"`
	Speaker    string  `json:"speaker"`
	Text       string  `json:"text"`
	Start      int64   `json:"start"`
	End        int64   `json:"end"`
	Confidence float64 `json:"confidence"`
}

type TranscriptionStatusResponse struct {
	ID                       string                 `json:"id"`
	Status                   string                 `json:"status"`
	Text                     string                 `json:"text"`
	Utterances               []*TranscriptUtterance `json:"utterances"`
	Summary                  string                 `json:"summary"`
	SentimentAnalysisResults []*Sentiment           `json:"sentiment_analysis_results"`
	Entities                 []*Entity              `json:"entities"`
}

func NewClient(logger *zap.Logger, httpClient *http.Client, cfg *Config) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     cfg.AccessToken,
		logger:     logger,
	}
}

type Client struct {
	httpClient *http.Client
	apiKey     string
	logger     *zap.Logger
}

func (c *Client) Transcribe(ctx context.Context, req *TranscribeRequest) (*TranscriptionStatusResponse, error) {
	c.logger.Info("Submitting job...")
	job, err := c.submitQuery(req)
	if err != nil {
		return nil, err
	}

	if job.ID == "" {
		return nil, fmt.Errorf("no job ID returned (status: %s)", job.Status)
	}

	c.logger.Info("Job submitted OK, awaiting result...", zap.String("id", job.ID), zap.String("status", job.Status))
	result, err := c.awaitResult(ctx, job.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) submitQuery(reqBody *TranscribeRequest) (*TranscriptionStatusResponse, error) {

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(reqBody); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", APIBaseURL, "transcript"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", c.apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, resp.Body); err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, buff.String())
	}

	status := &TranscriptionStatusResponse{}
	if err := json.NewDecoder(buff).Decode(status); err != nil {
		return nil, err
	}

	if status.ID == "" {
		_ = os.WriteFile("response.dump.json", buff.Bytes(), 0644)
		return nil, errors.New("unexpected result, dumping response to response.dump.json")
	}

	return status, nil
}

func (c *Client) awaitResult(ctx context.Context, jobID string) (*TranscriptionStatusResponse, error) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			resp, err := c.getStatus(jobID)
			if err != nil {
				return nil, err
			}
			switch resp.Status {
			case "completed":
				return resp, nil
			case "error":
				return resp, fmt.Errorf("job did not complete")
			default:
				c.logger.Info("Pending...", zap.String("id", resp.ID), zap.String("status", resp.Status))
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for transcript to complete")

		}
	}
}

func (c *Client) getStatus(jobID string) (*TranscriptionStatusResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", APIBaseURL, "transcript", jobID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", c.apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	status := &TranscriptionStatusResponse{}
	if err := json.NewDecoder(resp.Body).Decode(status); err != nil {
		return nil, err
	}

	// add a position property to assist in debugging other scripts
	for k := range status.Utterances {
		status.Utterances[k].Pos = int64(k) + 1
	}

	return status, nil
}
