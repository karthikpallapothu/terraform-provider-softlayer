/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package session

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/softlayer/softlayer-go/config"
	"github.com/softlayer/softlayer-go/sl"
)

// DefaultEndpoint is the default endpoint for API calls, when no override
// is provided.
const DefaultEndpoint = "https://api.softlayer.com/rest/v3"

type TransportHandlerFunc func(
	sess *Session,
	service string,
	method string,
	args []interface{},
	options *sl.Options,
	pResult interface{}) error

// Session stores the information required for communication with the SoftLayer
// API
type Session struct {
	// UserName is the name of the SoftLayer API user
	UserName string

	// ApiKey is the secret for making API calls
	APIKey string

	// Endpoint is the SoftLayer API endpoint to communicate with
	Endpoint string

	// Debug controls logging of request details (URI, parameters, etc.)
	Debug bool

	// The function that will be called for each API request.  Handles the
	// request and any response parsing specific to the desired protocol
	// (e.g., REST).  Set automatically for a new Session, based on the
	// provided Endpoint.
	TransportHandler TransportHandlerFunc
}

// New creates and returns a pointer to a new session object.  It takes up to
// three parameters, all of which are optional.  If specified, they will be
// interpreted in the following sequence:
//
// 1. UserName
// 2. Api Key
// 3. Endpoint
//
// If one or more are omitted, New() will attempt to retrieve these values from
// the environment, and the ~/.softlayer config file, in that order.
func New(args ...interface{}) *Session {
	keys := map[string]int{"username": 0, "api_key": 1, "endpoint_url": 2}
	values := []string{"", "", ""}

	for i := 0; i < len(args); i++ {
		values[i] = args[i].(string)
	}

	// Default to the environment variables
	envFallback("SOFTLAYER_USERNAME", &values[keys["username"]])
	envFallback("SOFTLAYER_API_KEY", &values[keys["api_key"]])
	envFallback("SOFTLAYER_ENDPOINT_URL", &values[keys["endpoint_url"]])

	// Read ~/.softlayer for configuration
	u, err := user.Current()
	if err != nil {
		panic("session: Could not determine current user.")
	}

	configPath := fmt.Sprintf("%s/.softlayer", u.HomeDir)
	if _, err = os.Stat(configPath); !os.IsNotExist(err) {
		// config file exists
		file, err := config.LoadFile(configPath)
		if err != nil {
			log.Println(fmt.Sprintf("[WARN] session: Could not parse %s : %s", configPath, err))
		} else {
			for k, v := range keys {
				value, ok := file.Get("softlayer", k)
				if ok && values[v] == "" {
					values[v] = value
				}
			}
		}
	}

	endpointURL := values[keys["endpoint_url"]]
	if endpointURL == "" || !strings.Contains(endpointURL, "/rest/") {
		endpointURL = DefaultEndpoint
	}

	return &Session{
		UserName:         values[keys["username"]],
		APIKey:           values[keys["api_key"]],
		Endpoint:         endpointURL,
		TransportHandler: doRestRequest,
	}
}

// DoRequest hands off the processing to the assigned transport handler. It is
// normally called internally by the service objects, but is exported so that it can
// be invoked directly by client code in exceptional cases where direct control is
// needed over one of the parameters.
//
// service and method are the SoftLayer service name and method name, exactly as they
// are documented at http://sldn.softlayer.com/reference/softlayerapi (i.e., with the
// 'SoftLayer_' prefix and properly cased.
//
// args is a slice of arguments required for the service method being invoked.  The
// types of each argument varies. See the method definition in the services package
// for the expected type of each argument.
//
// options is an sl.Options struct, containing any mask, filter, or result limit values
// to be applied.
//
// pResult is a pointer to a variable to be populated with the result of the API call.
// The type of the variable pointed to determines how the response is handled.  E.g.,
// for simple integer or string types, a type conversion is attempted (e.g.,
// strconv.Atoi()). For a map or struct type, the response is unmarshaled into pResult
//
// A sl.Error is returned, and can be (with a type assertion) inspected for details of
// the error (http code, API error message, etc.), or simply handled as a generic error,
// (in which case no type assertion would be necessary)
func (r *Session) DoRequest(service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	if r.TransportHandler == nil {
		r.TransportHandler = doRestRequest
	}

	return r.TransportHandler(r, service, method, args, options, pResult)
}

func envFallback(keyName string, value *string) {
	if *value == "" {
		*value = os.Getenv(keyName)
	}
}
