package salesforce_utils

import "github.com/valyala/fasthttp"

// sendRequest sends a configured request, returning the body, status code, and error
func sendRequest(req *fasthttp.Request) ([]byte, int, error) {
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	err := fasthttp.Do(req, res)
	return res.Body(), res.StatusCode(), err
}
