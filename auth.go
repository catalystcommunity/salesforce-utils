package salesforce_utils

import (
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
	"net/http"
)

// SalesforceCredentials represents the response from salesforce's /services/oauth2/token endpoint to get an access token
type SalesforceCredentials struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
	Id          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    int    `json:"issued_at,string"`
	Signature   string `json:"signature"`
}

// Authenticate authenticates with salesforce, storing the resulting credentials on the SalesforceUtils object
func (s SalesforceUtils) Authenticate() error {
	body, statusCode, err := s.getSalesforceAccessToken()
	if err != nil || statusCode != 200 {
		return errorx.Decorate(err, "error getting access token with status code: %d and body: %s", statusCode, body)
	}
	var creds SalesforceCredentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return errorx.Decorate(err, "error unmarshalling response from salesforce into credentials object")
	}
	s.Credentials = creds
	return nil
}

// getSalesforceAccessToken makes an http request to the salesforce api to get an access token
func (s SalesforceUtils) getSalesforceAccessToken() ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getAuthUrl()
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPost)
	return sendRequest(req)
}

// getAuthUrl gets a formatted url to the token endpoint
func (s SalesforceUtils) getAuthUrl() string {
	return fmt.Sprintf("%s/services/oauth2/token?client_id=%s&client_secret=%s&username=%s&password=%s&grant_type=%s", s.Config.BaseUrl, s.Config.ClientId, s.Config.ClientSecret, s.Config.Username, s.Config.Password, s.Config.GrantType)
}
