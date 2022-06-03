package salesforce_utils

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

// sendRequest sends a configured request, returning the body, status code, and error
func (s SalesforceUtils) sendRequest(req *fasthttp.Request) ([]byte, int, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Credentials.AccessToken))
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	err := fasthttp.Do(req, res)
	return res.Body(), res.StatusCode(), err
}
