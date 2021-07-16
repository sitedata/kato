# envoy_discover_service

`PREFIX`  URL prefix path configuration, for example /api

`DOMAINS` Intranet requests domain name configuration, and forwards it to downstream applications based on the configured domain name

`LIMITS`  TCP rate limit, the configuration range is 0-2048, fill in the number in the box, if it is configured to be 0, it will be fuse

`MaxPendingRequests` HTTP suspend request, the configuration range is 0-2048, fill in the number in the box, configure 0 to suspend the request immediately

`WEIGHT` Forwarding weight setting, ranging from 1 to 100, this parameter will judge multiple downstream services with the same domain name for weight distribution, the sum of weights must be 100, otherwise it will be inaccessible

`HEADERS` HTTP request header setting, in k:v format, multiple separated by ";", for example header1:mm;header2:nn

`MaxRequests` The maximum number of requests limit is 1024 by default, set 0 to 0 requests

`MaxRetries` The maximum number of retries is 3 by default, set 0 to 0 to retry
