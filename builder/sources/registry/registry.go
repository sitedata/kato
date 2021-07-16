// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package registry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//LogfCallback LogfCallback
type LogfCallback func(format string, args ...interface{})

//Quiet Quiet
/*
 * Discard log messages silently.
 */
func Quiet(format string, args ...interface{}) {
	/* discard logs */
}

//Log print log
/*
 * Pass log messages along to Go's "log" module.
 */
func Log(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

//Registry the client for  image repostory
type Registry struct {
	URL    string
	Client *http.Client
	Logf   LogfCallback
}

//New new ssl registry client
/*
 * Create a new Registry with the given URL and credentials, then Ping()s it
 * before returning it to verify that the registry is available.
 *
 * You can, alternately, construct a Registry manually by populating the fields.
 * This passes http.DefaultTransport to WrapTransport when creating the
 * http.Client.
 */
func New(registryURL, username, password string) (*Registry, error) {
	transport := http.DefaultTransport
	return newFromTransport(registryURL, username, password, transport, Log)
}

//NewInsecure new insecure skip verify tls client
/*
 * Create a new Registry, as with New, using an http.Transport that disables
 * SSL certificate verification.
 */
func NewInsecure(registryURL, username, password string) (*Registry, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return newFromTransport(registryURL, username, password, transport, Log)
}

// WrapTransport returns an existing http.RoundTripper such as http.DefaultTransport,
// build the transport stack necessary to authenticate to the Docker registry API.
// This adds in support for OAuth bearer tokens and HTTP Basic auth, and sets up
// error handling this library relies on.
func WrapTransport(transport http.RoundTripper, url, username, password string) http.RoundTripper {
	tokenTransport := &TokenTransport{
		Transport: transport,
		Username:  username,
		Password:  password,
	}
	basicAuthTransport := &BasicTransport{
		Transport: tokenTransport,
		URL:       url,
		Username:  username,
		Password:  password,
	}
	errorTransport := &ErrorTransport{
		Transport: basicAuthTransport,
	}
	return errorTransport
}

func newFromTransport(registryURL, username, password string, transport http.RoundTripper, logf LogfCallback) (*Registry, error) {
	url := strings.TrimSuffix(registryURL, "/")
	containsScheme := strings.HasPrefix(url, "http")
	if !containsScheme {
		// use https by default
		url = fmt.Sprintf("https://%s", registryURL)
	}

	if username != "" {
		transport = WrapTransport(transport, url, username, password)
	}
	registry := &Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: logf,
	}

	if err := registry.Ping(); err != nil {
		if errors.Is(err, ErrRegistryNotFound) && !containsScheme {
			// try again with http url
			registry.URL = strings.Replace(url, "https", "http", 1)
			if err := registry.Ping(); err != nil {
				return nil, errors.Wrap(err, "registry url "+url)
			}
			return registry, nil
		}
		return nil, errors.Wrap(err, "registry url "+url)
	}

	return registry, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

//Ping ping registry server
func (r *Registry) Ping() error {
	url := r.url("/v2/")
	r.Logf("registry.ping url=%s", url)
	resp, err := r.Client.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return errors.Wrap(ErrRegistryNotFound, err.Error())
		}
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return errors.Wrap(ErrRegistryNotFound, "not found")
		}
	}
	return nil
}
