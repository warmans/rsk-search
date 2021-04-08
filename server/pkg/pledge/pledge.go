package pledge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"net/http"
)

type Config struct {
	APIKey            string
	DonationEmail     string
	DonationFirstName string
	DonationLastName  string
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.APIKey, prefix, "pledge-api-key", "", "API key")
	flag.StringVarEnv(fs, &c.DonationEmail, prefix, "default-donation-email", "", "Email to use for anonymous donation request")
	flag.StringVarEnv(fs, &c.DonationFirstName, prefix, "default-donation-first-name", "", "First name to use for anonymous donation request")
	flag.StringVarEnv(fs, &c.DonationLastName, prefix, "default-donation-last-name", "", "Last name to use for anonymous donation request")
}

type Client struct {
	cfg Config
}

func (c *Client) CreateAnonymousDonation(donationDetails AnonymousDonationRequest) (*Donation, error) {
	return c.CreateDonation(DonationRequest{
		Email:          c.cfg.DonationEmail,
		FirstName:      c.cfg.DonationFirstName,
		LastName:       c.cfg.DonationFirstName,
		OrganizationID: donationDetails.OrganizationID,
		Amount:         donationDetails.Amount,
		Metadata:       donationDetails.Metadata,
	})
}

func (c *Client) CreateDonation(donationDetails DonationRequest) (*Donation, error) {
	b, err := json.Marshal(donationDetails)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode request")
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.pledge.to/v1/donations", bytes.NewBuffer(b))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	result := &Donation{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return result, nil
}

func (c *Client) ListOrganizations() (*OrganizationList, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.pledge.to/v1/organizations", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	list := &OrganizationList{}
	if err := json.NewDecoder(resp.Body).Decode(list); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return list, nil
}

func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 400:
		return fmt.Errorf("bad request")
	case 401:
		return fmt.Errorf("API key was invalid")
	case 404:
		return fmt.Errorf("resource does not exist")
	case 422:
		return fmt.Errorf("service could not process request")
	case 500:
		return fmt.Errorf("internal error")
	}
	return nil
}

type OrganizationList struct {
	Page       int32        `json:"page"`
	Per        int32        `json:"per"`
	TotalCount int32        `json:"total_count"`
	Results    Organization `json:"results"`
}

type Organization struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Mission    string `json:"mission"`
	WebsiteURL string `json:"website_url"`
	LogoURL    string `json:"logo_url"`
}

type DonationRequest struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	OrganizationID string `json:"organization_id"`
	Amount         string `json:"amount"`
	Metadata       string `json:"metadata"`
}

type AnonymousDonationRequest struct {
	OrganizationID string `json:"organization_id"`
	Amount         string `json:"amount"`
	Metadata       string `json:"metadata"`
}

type Donation struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Amount           string `json:"amount"`
	Status           string `json:"status"`
	Metadata         string `json:"metadata"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
