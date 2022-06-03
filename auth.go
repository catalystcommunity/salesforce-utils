package salesforce_utils

import (
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
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
	params := url.Values{}
	params.Add("client_id", s.Config.ClientId)
	params.Add("client_secret", s.Config.ClientSecret)
	params.Add("username", s.Config.Username)
	params.Add("password", s.Config.Password)
	params.Add("grant_type", s.Config.GrantType)
	return fmt.Sprintf("%s/services/oauth2/token?%s", s.Config.BaseUrl, params.Encode())
}
