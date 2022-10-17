package pkg

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// sendRequest sends a configured request, returning the body, status code, and error
func (s *SalesforceUtils) sendRequest(req *fasthttp.Request) ([]byte, int, func(), error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Credentials.AccessToken))
	res := fasthttp.AcquireResponse()
	err := s.FastHTTPClient.Do(req, res)
	return res.Body(), res.StatusCode(), func() { fasthttp.ReleaseResponse(res) }, err
}
