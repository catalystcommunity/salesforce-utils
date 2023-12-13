package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
)

// CompositeObject is a type used for input into CreateObjects. Contains all
// the information that the library needs to build a CompositeRequest to send
// to salesforce
type CompositeObject struct {
	// SalesforceId is only used for UpdateObjects and DeleteObjects
	SalesforceId string
	// ReferenceId is used to correlate the response to the request
	ReferenceId string
	// ObjectType is the salesforce object type, e.g. "Account"
	ObjectType string
	// Body is the json body to send to salesforce
	Body []byte
}

// CompositeRequest is used by the CreateObjects method when interacting with
// the salesforce API
//
// ref:
// https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/requests_composite.htm
type CompositeRequest struct {
	AllOrNone        bool                  `json:"allOrNone"`
	CompositeRequest []CompositeSubRequest `json:"compositeRequest"`
}

type CompositeSubRequest struct {
	// the use of json.RawMessage here is to avoid the need to unmarshal the
	// body from json and remarshal it.
	Body        json.RawMessage `json:"body"`
	HttpHeaders interface{}     `json:"httpHeaders"`
	Method      string          `json:"method"`
	ReferenceId string          `json:"referenceId"`
	Url         string          `json:"url"`
}

type CompositeResponse struct {
	CompositeResponse []CompositeSubResponse `json:"compositeResponse"`
}

type CompositeSubResponse struct {
	Body           CompositeSubResponseBody `json:"body"`
	HttpStatusCode int                      `json:"httpStatusCode"`
	ReferenceId    string                   `json:"referenceId"`
}

type CompositeSubResponseBody struct {
	Id      string   `json:"id"`
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func (s *SalesforceUtils) CreateObjects(objects []CompositeObject) (response CompositeResponse, err error) {
	compositeReq := s.convertToCompositeCreateRequest(objects)
	return s.doCompositeRequest(compositeReq)
}

func (s *SalesforceUtils) UpdateObjects(objects []CompositeObject) (response CompositeResponse, err error) {
	compositeReq := s.convertToCompositeUpdateRequest(objects)
	return s.doCompositeRequest(compositeReq)
}

func (s *SalesforceUtils) UpsertObjects(objects []CompositeObject) (response CompositeResponse, err error) {
	compositeReq := s.convertToCompositeUpsertRequest(objects)
	return s.doCompositeRequest(compositeReq)
}

func (s *SalesforceUtils) DeleteObjects(objects []CompositeObject) (response CompositeResponse, err error) {
	compositeReq := s.convertToCompositeDeleteRequest(objects)
	return s.doCompositeRequest(compositeReq)
}

// convertToCompositeCreateRequest converts a list of CompositeObjects intended
// for creation into a single CompositeRequest object that can be marshalled to
// json to send to salesforce
func (s *SalesforceUtils) convertToCompositeCreateRequest(objects []CompositeObject) CompositeRequest {
	compositeReq := CompositeRequest{
		// always set to true, so that the api behaves transactionally.
		AllOrNone:        true,
		CompositeRequest: []CompositeSubRequest{},
	}

	for _, obj := range objects {
		compositeReq.CompositeRequest = append(compositeReq.CompositeRequest, CompositeSubRequest{
			Body:        obj.Body,
			Method:      "POST",
			ReferenceId: obj.ReferenceId,
			Url:         s.getTypePath(obj.ObjectType),
		})
	}

	return compositeReq
}

// convertToCompositeUpdateRequest converts a list of CompositeObjects intended
// for update into a single CompositeRequest object that can be marshalled to
// json to send to salesforce
func (s *SalesforceUtils) convertToCompositeUpdateRequest(objects []CompositeObject) CompositeRequest {
	compositeReq := CompositeRequest{
		// always set to true, so that the api behaves transactionally.
		AllOrNone:        true,
		CompositeRequest: []CompositeSubRequest{},
	}

	for _, obj := range objects {
		compositeReq.CompositeRequest = append(compositeReq.CompositeRequest, CompositeSubRequest{
			Body:        obj.Body,
			Method:      "PATCH",
			ReferenceId: obj.ReferenceId,
			Url:         s.getObjectIdPath(obj.ObjectType, obj.SalesforceId),
		})
	}

	return compositeReq
}

// convertToCompositeUpsertRequest converts a list of CompositeObjects intended
// for either a create or update into a single CompositeRequest object that can
// be marshalled to json to send to the salesforce composite api. determines
// update vs create based on the existence of the salesforce id
func (s *SalesforceUtils) convertToCompositeUpsertRequest(objects []CompositeObject) CompositeRequest {
	compositeReq := CompositeRequest{
		// always set to true, so that the api behaves transactionally.
		AllOrNone:        true,
		CompositeRequest: []CompositeSubRequest{},
	}

	for _, obj := range objects {
		if obj.SalesforceId == "" {
			compositeReq.CompositeRequest = append(compositeReq.CompositeRequest, CompositeSubRequest{
				Body:        obj.Body,
				Method:      "POST",
				ReferenceId: obj.ReferenceId,
				Url:         s.getTypePath(obj.ObjectType),
			})
		} else {
			compositeReq.CompositeRequest = append(compositeReq.CompositeRequest, CompositeSubRequest{
				Body:        obj.Body,
				Method:      "PATCH",
				ReferenceId: obj.ReferenceId,
				Url:         s.getObjectIdPath(obj.ObjectType, obj.SalesforceId),
			})
		}
	}

	return compositeReq
}

// convertToCompositeDeleteRequest converts a list of CompositeObjects intended
// for deletion into a single CompositeRequest object that can be marshalled to
// json to send to salesforce
func (s *SalesforceUtils) convertToCompositeDeleteRequest(objects []CompositeObject) CompositeRequest {
	compositeReq := CompositeRequest{
		// always set to true, so that the api behaves transactionally.
		AllOrNone:        true,
		CompositeRequest: []CompositeSubRequest{},
	}

	for _, obj := range objects {
		compositeReq.CompositeRequest = append(compositeReq.CompositeRequest, CompositeSubRequest{
			Method:      "DELETE",
			ReferenceId: obj.ReferenceId,
			Url:         s.getObjectIdPath(obj.ObjectType, obj.SalesforceId),
		})
	}

	return compositeReq
}

func (s *SalesforceUtils) doCompositeRequest(compositeRequest CompositeRequest) (response CompositeResponse, err error) {
	reqBodyBytes, err := json.Marshal(compositeRequest)
	if err != nil {
		return response, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getCompositeUrl()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(reqBodyBytes)

	body, statusCode, deferredFunc, err := s.sendRequest(req)
	defer deferredFunc()
	if err != nil {
		return response, err
	}
	if statusCode != fasthttp.StatusOK {
		return response, errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, errorx.IllegalState.New("failed to unmarshal response: %s", err)
	}

	return response, nil
}

func (s *SalesforceUtils) getCompositeUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/composite", s.Config.BaseUrl, s.Config.ApiVersion)
}
