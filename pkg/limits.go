package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/valyala/fasthttp"
)

// LimitsResponse is the response from the limits endpoint
// https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_limits.htm?q=limits
type LimitsResponse struct {
	AnalyticsExternalDataSizeMB                 Limit               `json:"AnalyticsExternalDataSizeMB"`
	ConcurrentAsyncGetReportInstances           Limit               `json:"ConcurrentAsyncGetReportInstances"`
	ConcurrentEinsteinDataInsightsStoryCreation Limit               `json:"ConcurrentEinsteinDataInsightsStoryCreation"`
	ConcurrentEinsteinDiscoveryStoryCreation    Limit               `json:"ConcurrentEinsteinDiscoveryStoryCreation"`
	ConcurrentSyncReportRuns                    Limit               `json:"ConcurrentSyncReportRuns"`
	DailyAnalyticsDataflowJobExecutions         Limit               `json:"DailyAnalyticsDataflowJobExecutions"`
	DailyAnalyticsUploadedFilesSizeMB           Limit               `json:"DailyAnalyticsUploadedFilesSizeMB"`
	DailyApiRequests                            Limit               `json:"DailyApiRequests"`
	DailyAsyncApexExecutions                    Limit               `json:"DailyAsyncApexExecutions"`
	DailyAsyncApexTests                         Limit               `json:"DailyAsyncApexTests"`
	DailyBulkApiBatches                         Limit               `json:"DailyBulkApiBatches"`
	DailyBulkV2QueryFileStorageMB               Limit               `json:"DailyBulkV2QueryFileStorageMB"`
	DailyBulkV2QueryJobs                        Limit               `json:"DailyBulkV2QueryJobs"`
	DailyDeliveredPlatformEvents                Limit               `json:"DailyDeliveredPlatformEvents"`
	DailyDurableGenericStreamingApiEvents       Limit               `json:"DailyDurableGenericStreamingApiEvents"`
	DailyDurableStreamingApiEvents              Limit               `json:"DailyDurableStreamingApiEvents"`
	DailyEinsteinDataInsightsStoryCreation      Limit               `json:"DailyEinsteinDataInsightsStoryCreation"`
	DailyEinsteinDiscoveryOptimizationJobRuns   Limit               `json:"DailyEinsteinDiscoveryOptimizationJobRuns"`
	DailyEinsteinDiscoveryPredictAPICalls       Limit               `json:"DailyEinsteinDiscoveryPredictAPICalls"`
	DailyEinsteinDiscoveryPredictionsByCDC      Limit               `json:"DailyEinsteinDiscoveryPredictionsByCDC"`
	DailyEinsteinDiscoveryStoryCreation         Limit               `json:"DailyEinsteinDiscoveryStoryCreation"`
	DailyFunctionsApiCallLimit                  Limit               `json:"DailyFunctionsApiCallLimit"`
	DailyGenericStreamingApiEvents              Limit               `json:"DailyGenericStreamingApiEvents"`
	DailyStandardVolumePlatformEvents           Limit               `json:"DailyStandardVolumePlatformEvents"`
	DailyStreamingApiEvents                     Limit               `json:"DailyStreamingApiEvents"`
	DailyWorkflowEmails                         Limit               `json:"DailyWorkflowEmails"`
	DataStorageMB                               Limit               `json:"DataStorageMB"`
	DurableStreamingApiConcurrentClients        Limit               `json:"DurableStreamingApiConcurrentClients"`
	FileStorageMB                               Limit               `json:"FileStorageMB"`
	HourlyAsyncReportRuns                       Limit               `json:"HourlyAsyncReportRuns"`
	HourlyDashboardRefreshes                    Limit               `json:"HourlyDashboardRefreshes"`
	HourlyDashboardResults                      Limit               `json:"HourlyDashboardResults"`
	HourlyDashboardStatuses                     Limit               `json:"HourlyDashboardStatuses"`
	HourlyLongTermIdMapping                     Limit               `json:"HourlyLongTermIdMapping"`
	HourlyManagedContentPublicRequests          Limit               `json:"HourlyManagedContentPublicRequests"`
	HourlyODataCallout                          Limit               `json:"HourlyODataCallout"`
	HourlyPublishedPlatformEvents               Limit               `json:"HourlyPublishedPlatformEvents"`
	HourlyPublishedStandardVolumePlatformEvents Limit               `json:"HourlyPublishedStandardVolumePlatformEvents"`
	HourlyShortTermIdMapping                    Limit               `json:"HourlyShortTermIdMapping"`
	HourlySyncReportRuns                        Limit               `json:"HourlySyncReportRuns"`
	HourlyTimeBasedWorkflow                     Limit               `json:"HourlyTimeBasedWorkflow"`
	MassEmail                                   Limit               `json:"MassEmail"`
	MonthlyEinsteinDiscoveryStoryCreation       Limit               `json:"MonthlyEinsteinDiscoveryStoryCreation"`
	MonthlyPlatformEventsUsageEntitlement       Limit               `json:"MonthlyPlatformEventsUsageEntitlement"`
	Package2VersionCreates                      Limit               `json:"Package2VersionCreates"`
	Package2VersionCreatesWithoutValidation     Limit               `json:"Package2VersionCreatesWithoutValidation"`
	PermissionSets                              PermissionSetsLimit `json:"PermissionSets"`
	PrivateConnectOutboundCalloutHourlyLimitMB  Limit               `json:"PrivateConnectOutboundCalloutHourlyLimitMB"`
	PublishCallbackUsageInApex                  Limit               `json:"PublishCallbackUsageInApex"`
	SingleEmail                                 Limit               `json:"SingleEmail"`
	StreamingApiConcurrentClients               Limit               `json:"StreamingApiConcurrentClients"`
}

type Limit struct {
	Max       int `json:"Max"`
	Remaining int `json:"Remaining"`
}

type PermissionSetsLimit struct {
	Max          int   `json:"Max"`
	Remaining    int   `json:"Remaining"`
	CreateCustom Limit `json:"CreateCustom"`
}

func (s *SalesforceUtils) GetLimits() (*LimitsResponse, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getLimitsUrl()
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	body, statusCode, deferredFunc, err := s.sendRequest(req)
	defer deferredFunc()
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d with body: %s", statusCode, body)
	}
	response := &LimitsResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		logging.Log.WithField("body", string(body)).Error("failed to unmarshal response")
		return nil, err
	}
	return response, nil
}

func (s *SalesforceUtils) getLimitsUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/limits", s.Config.BaseUrl, s.Config.ApiVersion)
}
