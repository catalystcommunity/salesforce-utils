package salesforce_utils

import (
	"github.com/asaskevich/govalidator"
	"github.com/catalystsquad/app-utils-go/env"
	"os"
)

// SalesforceUtils is the struct that holds config and credentials
type SalesforceUtils struct {
	Config      Config
	Credentials SalesforceCredentials
}

type Config struct {
	BaseUrl      string `valid:"url,required"`
	ApiVersion   string `valid:"alphanum,required"`
	ClientId     string `valid:"alphanum,required"`
	ClientSecret string `valid:"alphanum,required"`
	Username     string `valid:"alphanum,required"`
	Password     string `valid:"alphanum,required"`
	GrantType    string `valid:"alphanum,required"`
}

// NewSalesforceUtils creates a new instance of SalesforceUtils with the given configuration. If any configuration is
// not set, it will look up the value from environment variables instead. If the configuration is inavlid, an error
// will be returned.
func NewSalesforceUtils(authenticate bool, config Config) (utils SalesforceUtils, err error) {
	if config.BaseUrl == "" {
		config.BaseUrl = os.Getenv("SALESFORCE_BASE_URL")
	}
	if config.ApiVersion == "" {
		config.ApiVersion = os.Getenv("SALESFORCE_API_VERSION")
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
	_, err = govalidator.ValidateStruct(config)
	if err != nil {
		return
	}
	utils.Config = config
	// authenticate
	if authenticate {
		err = utils.Authenticate()
	}
	return
}
