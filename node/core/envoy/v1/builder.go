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

//CreatOutlierDetection  create outlierDetection
//https://www.envoyproxy.io/docs/envoy/latest/api-v1/cluster_manager/cluster_outlier_detection#config-cluster-manager-cluster-outlier-detection
func CreatOutlierDetection(options map[string]interface{}) *OutlierDetection {
	if _, ok := options[KeyMaxConnections]; !ok {
		return nil
	}
	var od OutlierDetection
	od.ConsecutiveErrors = GetOptionValues(KeyConsecutiveErrors, options).(int)
	od.IntervalMS = GetOptionValues(KeyIntervalMS, options).(int64)
	od.BaseEjectionTimeMS = GetOptionValues(KeyBaseEjectionTimeMS, options).(int64)
	od.MaxEjectionPercent = GetOptionValues(KeyMaxEjectionPercent, options).(int)
	return &od
}

//CreateCircuitBreaker create circuitBreaker
//https://www.envoyproxy.io/docs/envoy/latest/api-v1/cluster_manager/cluster_circuit_breakers#config-cluster-manager-cluster-circuit-breakers-v1
func CreateCircuitBreaker(options map[string]interface{}) *CircuitBreaker {
	if _, ok := options[KeyMaxConnections]; !ok {
		return nil
	}
	var cb CircuitBreaker
	cb.Default.MaxConnections = GetOptionValues(KeyMaxConnections, options).(int)
	cb.Default.MaxRequests = GetOptionValues(KeyMaxRequests, options).(int)
	cb.Default.MaxRetries = GetOptionValues(KeyMaxActiveRetries, options).(int)
	cb.Default.MaxPendingRequests = GetOptionValues(KeyMaxPendingRequests, options).(int)
	return &cb
}
