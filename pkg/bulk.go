package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
)

type BulkJobRecord struct {
	ID                              string  `json:"id"`
	Operation                       string  `json:"operation"`
	Object                          string  `json:"object"`
	CreatedById                     string  `json:"createdById"`
	CreatedDate                     string  `json:"createdDate"`
	SystemModstamp                  string  `json:"systemModstamp"`
	State                           string  `json:"state"`
	ConcurrencyMode                 string  `json:"concurrencyMode"`
	ContentType                     string  `json:"contentType"`
	ApiVersion                      float64 `json:"apiVersion"`
	LineEnding                      string  `json:"lineEnding"`
	ColumnDelimiter                 string  `json:"columnDelimiter"`
	NumberRecordsProcessed          int64   `json:"numberRecordsProcessed"`
	Retries                         int64   `json:"retries"`
	TotalProcessingTimeMilliseconds int64   `json:"totalProcessingTime"`
}

func (s *SalesforceUtils) CreateBulkQueryJob(query string) (response BulkJobRecord, err error) {
	queryBody := map[string]string{
		"operation": "query",
		"query":     query,
	}
	queryBodyBytes, err := json.Marshal(queryBody)
	if err != nil {
		err = errorx.Decorate(err, "failed to marshal query")
		return
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getBulkUrl()
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(queryBodyBytes)
	responseBody, statusCode, requestErr := s.sendRequest(req)
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, responseBody)
		return
	}
	err = json.Unmarshal(responseBody, &response)
	return
}

func (s *SalesforceUtils) getBulkUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/jobs/query", s.Config.BaseUrl, s.Config.ApiVersion)
}

func (s *SalesforceUtils) GetBulkQueryJob(queryJobID string) (response BulkJobRecord, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getBulkQueryJobInfoUrl(queryJobID)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	body, statusCode, requestErr := s.sendRequest(req)
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) getBulkQueryJobInfoUrl(queryJobID string) string {
	return fmt.Sprintf("%s/%s", s.getBulkUrl(), queryJobID)
}

type GetBulkQueryJobResultsResponse struct {
	NumberOfRecords int
	Locator         string
	Body            []byte
}

func (s *SalesforceUtils) GetBulkQueryJobResults(queryJobID string, locator string) (response GetBulkQueryJobResultsResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getBulkQueryJobResultsUrl(queryJobID, locator)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)

	// send the request without the sendRequest helper, so that we can get
	// headers from the response
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Credentials.AccessToken))
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	err = s.FastHTTPClient.Do(req, res)
	if err != nil {
		return
	}
	if res.StatusCode() != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", res.StatusCode(), res.Body())
		return
	}

	locator = string(res.Header.Peek("Sforce-Locator"))
	if locator == "null" {
		locator = ""
	}

	response.Locator = locator
	response.NumberOfRecords, _ = strconv.Atoi(string(res.Header.Peek("Sforce-NumberOfRecords")))
	response.Body = res.Body()
	return
}

func (s *SalesforceUtils) getBulkQueryJobResultsUrl(queryJobID string, locator string) string {
	uri := fmt.Sprintf("%s/%s/results", s.getBulkUrl(), queryJobID)
	if locator != "" {
		params := url.Values{}
		params.Add("locator", locator)
		uri = fmt.Sprintf("%s?%s", uri, params.Encode())
	}
	return uri
}

type ListBulkJobsResponse struct {
	Done           bool            `json:"done"`
	Records        []BulkJobRecord `json:"records"`
	NextRecordsUrl string          `json:"nextRecordsUrl"`
}

func (s *SalesforceUtils) ListBulkJobs() (response ListBulkJobsResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getListBulkJobsUrl()
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	body, statusCode, requestErr := s.sendRequest(req)
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) getListBulkJobsUrl() string {
	return fmt.Sprintf("%s/", s.getBulkUrl())
}
