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

package v1

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	//KeyPrefix request path prefix
	KeyPrefix string = "Prefix"
	//KeyHeaders request http headers
	KeyHeaders string = "Headers"
	//KeyDomains request domains
	KeyDomains string = "Domains"
	//KeyMaxConnections The maximum number of connections that Envoy will make to the upstream cluster. If not specified, the default is 1024.
	KeyMaxConnections string = "MaxConnections"
	//KeyMaxPendingRequests The maximum number of pending requests that Envoy will allow to the upstream cluster. If not specified, the default is 1024
	KeyMaxPendingRequests string = "MaxPendingRequests"
	//KeyMaxRequests The maximum number of parallel requests that Envoy will make to the upstream cluster. If not specified, the default is 1024.
	KeyMaxRequests string = "MaxRequests"
	//KeyMaxActiveRetries  The maximum number of parallel retries that Envoy will allow to the upstream cluster. If not specified, the default is 3.
	KeyMaxActiveRetries string = "MaxActiveRetries"
	//KeyUpStream upStream
	KeyUpStream string = "upStream"
	//KeyDownStream downStream
	KeyDownStream string = "downStream"
	//KeyWeight WEIGHT
	KeyWeight string = "Weight"
	//KeyWeightModel MODEL_WEIGHT
	KeyWeightModel string = "weight_model"
	//KeyPrefixModel MODEL_PREFIX
	KeyPrefixModel string = "prefix_model"
	//KeyIntervalMS IntervalMS key
	KeyIntervalMS string = "IntervalMS"
	//KeyConsecutiveErrors ConsecutiveErrors key
	KeyConsecutiveErrors string = "ConsecutiveErrors"
	//KeyBaseEjectionTimeMS BaseEjectionTimeMS key
	KeyBaseEjectionTimeMS string = "BaseEjectionTimeMS"
	//KeyMaxEjectionPercent MaxEjectionPercent key
	KeyMaxEjectionPercent string = "MaxEjectionPercent"
)

//GetOptionValues get value from options
//if not exist,return default value
func GetOptionValues(kind string, sr map[string]interface{}) interface{} {
	switch kind {
	case KeyPrefix:
		if prefix, ok := sr[KeyPrefix]; ok {
			return prefix
		}
		return "/"
	case KeyMaxConnections:
		if circuit, ok := sr[KeyMaxConnections]; ok {
			cc, err := strconv.Atoi(circuit.(string))
			if err != nil {
				logrus.Errorf("strcon circuit error")
				return 1024
			}
			return cc
		}
		return 1024
	case KeyMaxRequests:
		if maxRequest, ok := sr[KeyMaxRequests]; ok {
			mrt, err := strconv.Atoi(maxRequest.(string))
			if err != nil {
				logrus.Errorf("strcon max request error")
				return 1024
			}
			return mrt
		}
		return 1024
	case KeyMaxPendingRequests:
		if maxPendingRequests, ok := sr[KeyMaxPendingRequests]; ok {
			mpr, err := strconv.Atoi(maxPendingRequests.(string))
			if err != nil {
				logrus.Errorf("strcon max pending request error")
				return 1024
			}
			return mpr
		}
		return 1024
	case KeyMaxActiveRetries:
		if maxRetries, ok := sr[KeyMaxActiveRetries]; ok {
			mxr, err := strconv.Atoi(maxRetries.(string))
			if err != nil {
				logrus.Errorf("strcon max retry error")
				return 3
			}
			return mxr
		}
		return 3
	case KeyHeaders:
		var np []Header
		if headers, ok := sr[KeyHeaders]; ok {
			parents := strings.Split(headers.(string), ";")
			for _, h := range parents {
				headers := strings.Split(h, ":")
				//has_header:no 默认
				if len(headers) == 2 {
					if headers[0] == "has_header" && headers[1] == "no" {
						continue
					}
					ph := Header{
						Name:  headers[0],
						Value: headers[1],
					}
					np = append(np, ph)
				}
			}
		}
		return np
	case KeyDomains:
		if domain, ok := sr[KeyDomains]; ok {
			if strings.Contains(domain.(string), ",") {
				mm := strings.Split(domain.(string), ",")
				return mm
			}
			return []string{domain.(string)}
		}
		return []string{"*"}
	case KeyWeight:
		if weight, ok := sr[KeyWeight]; ok {
			w, err := strconv.Atoi(weight.(string))
			if err != nil {
				return 100
			}
			return w
		}
		return 100
	case KeyIntervalMS:
		if in, ok := sr[KeyIntervalMS]; ok {
			w, err := strconv.Atoi(in.(string))
			if err != nil {
				return int64(10000)
			}
			return int64(w)
		}
		return int64(10000)
	case KeyConsecutiveErrors:
		if in, ok := sr[KeyConsecutiveErrors]; ok {
			w, err := strconv.Atoi(in.(string))
			if err != nil {
				return 5
			}
			return w
		}
		return 5
	case KeyBaseEjectionTimeMS:
		if in, ok := sr[KeyBaseEjectionTimeMS]; ok {
			w, err := strconv.Atoi(in.(string))
			if err != nil {
				return int64(30000)
			}
			return int64(w)
		}
		return int64(30000)
	case KeyMaxEjectionPercent:
		if in, ok := sr[KeyMaxEjectionPercent]; ok {
			w, err := strconv.Atoi(in.(string))
			if err != nil || w > 100 {
				return 10
			}
			return w
		}
		return 10
	default:
		return nil
	}
}
