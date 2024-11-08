package pkg

import (
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/catalystcommunity/app-utils-go/env"
	"github.com/valyala/fasthttp"
)

// SalesforceUtils is the struct that holds config and credentials
type SalesforceUtils struct {
	Config         Config
	Credentials    SalesforceCredentials
	FastHTTPClient *fasthttp.Client
}

type Config struct {
	BaseUrl        string `valid:"url,required"`
	ApiVersion     string `valid:"required"`
	ClientId       string `valid:"required"`
	ClientSecret   string `valid:"required"`
	Username       string `valid:"required"`
	Password       string `valid:"required"`
	GrantType      string `valid:"required"`
	FastHTTPClient *fasthttp.Client
}

// NewSalesforceUtils creates a new instance of SalesforceUtils with the given configuration. If any configuration is
// not set, it will look up the value from environment variables instead. If the configuration is inavlid, an error
// will be returned.
func NewSalesforceUtils(authenticate bool, config Config) (*SalesforceUtils, error) {
	if config.BaseUrl == "" {
		config.BaseUrl = os.Getenv("SALESFORCE_BASE_URL")
	}
	if config.ApiVersion == "" {
		config.ApiVersion = env.GetEnvOrDefault("SALESFORCE_API_VERSION", "55.0")
	}
	if config.ClientId == "" {
		config.ClientId = os.Getenv("SALESFORCE_CLIENT_ID")
	}
	if config.ClientSecret == "" {
		config.ClientSecret = os.Getenv("SALESFORCE_CLIENT_SECRET")
	}
	if config.Username == "" {
		config.Username = os.Getenv("SALESFORCE_USERNAME")
	}
	if config.Password == "" {
		config.Password = os.Getenv("SALESFORCE_PASSWORD")
	}
	if config.GrantType == "" {
		config.GrantType = env.GetEnvOrDefault("SALESFORCE_GRANT_TYPE", "password")
	}
	// validate the config
	_, err := govalidator.ValidateStruct(config)
	if err != nil {
		return nil, err
	}
	utils := &SalesforceUtils{Config: config}
	utils.Config = config
	// allow passing a custom fasthttp client, default to empty
	if config.FastHTTPClient != nil {
		utils.FastHTTPClient = config.FastHTTPClient
	} else {
		utils.FastHTTPClient = &fasthttp.Client{}
	}
	// authenticate
	if authenticate {
		err = utils.Authenticate()
	}
	return utils, err
}
