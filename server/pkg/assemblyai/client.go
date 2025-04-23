package assemblyai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const (
	APIBaseURL string = "https://api.assemblyai.com/v2"
)

type TranscribeRequest struct {
	AudioURL      string `json:"audio_url"`
	SpeakerLabels bool   `json:"speaker_labels"`
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
	ID     string `json:"id"`
	Status string `json:"status"`
	Text   string `json:"text"`
	//Words      []*TranscriptWord      `json:"words"` // don't care
	Utterances []*TranscriptUtterance `json:"utterances"`
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

	status := &TranscriptionStatusResponse{}
	if err := json.NewDecoder(resp.Body).Decode(status); err != nil {
		return nil, err
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
