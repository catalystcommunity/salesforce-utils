

# Salesforce Utils
This is a utility library for consuming the salesforce lightning API
## Supported Endpoints
* sobjects (object CRUD)
* query (SOQL queries)
## Usage Example
Instantiate a new instance using `NewSalesforceUtils()`. Configuration can be provided as environment variables, or in code.
```go  
package main  
  
import salesforce_utils "github.com/catalystsquad/salesforce-utils"
  
func main() {
	// if authenticate is true then authentication is done as part of instantiation. If authenticate is 	
	// false, you'll need to manually call the Authenticate() method.
	authenticate := true
	// empty config will use environment variables
	config := salesforce_utils.Config{}
  sfUtils, err := salesforce_utils.NewSalesforceUtils(authenticate, config)
  if err != nil {
		fmt.Printf("error instantiating salesforce utils: %s", err.Error())
	}
	// create an object
	myJson := `{"some": "stuff"}`
	response, sfErr := sfUtils.CreateObject("my_type_name", []byte(myJson))
	// make a soql query
	myQuery := 'select fields(all) from my_type_name'
	response, sfErr := sfUtils.ExecuteSoqlQuery(myQuery)
}  
```  
## Configuration
Configuration is handled by code or environment variables. Code variables take precedence.
### Code Example
```go
// if authenticate is true then authentication is done as part of instantiation. If authenticate is false, you'll need to manually call the Authenticate() method.
authenticate := true
config := salesforce_utils.Config{  
  BaseUrl:      "https://mydomain.my.salesforce.com",  
  ApiVersion:   "55.0",  
  ClientId:     "client_id_here",  
  ClientSecret: "client_secret_here",  
  Username:     "username_here",  
  Password:     "password_here",  
  GrantType:    "password",  
}
sfUtils, err := salesforce_utils.NewSalesforceUtils(authenticate, config)
```
### Environment Variables
|name|required|purpose|default|
|--|--|--|--|
|SALESFORCE_BASE_URL|yes|Set the salesforce base url, i.e. https://mydomain.my.salesforce.com | ""
|SALESFORCE_CLIENT_ID|yes|Set the connected app client id|  ""
|SALESFORCE_CLIENT_SECRET|yes|Set the connected app client secret|  ""
|SALESFORCE_USERNAME|yes|User to authenticate as|  ""
|SALESFORCE_PASSWORD|yes|Password to authenticate with |  ""
|SALESFORCE_GRANT_TYPE|no|Grant type, we advise not setting this and letting it use the default|  "password"
|SALESFORCE_API_VERSION|no|Salesforce api version to use|  "55.0"
