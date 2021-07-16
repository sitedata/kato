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

package store

import (
	"testing"

	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/gridworkz/kato/worker/server/pb"
	"github.com/stretchr/testify/assert"
)

func TestGetAppStatus(t *testing.T) {
	tests := []struct {
		name     string
		statuses map[string]string
		want     pb.AppStatus_Status
	}{
		{
			name: "nocomponent",
			want: pb.AppStatus_NIL,
		},
		{
			name: "undeploy",
			statuses: map[string]string{
				"apple":  v1.UNDEPLOY,
				"banana": v1.UNDEPLOY,
			},
			want: pb.AppStatus_NIL,
		},
		{
			name: "closed",
			statuses: map[string]string{
				"apple":  v1.UNDEPLOY,
				"banana": v1.CLOSED,
				"cat":    v1.CLOSED,
			},
			want: pb.AppStatus_CLOSED,
		},
		{
			name: "abnormal",
			statuses: map[string]string{
				"apple":  v1.ABNORMAL,
				"banana": v1.SOMEABNORMAL,
				"cat":    v1.RUNNING,
				"dog":    v1.CLOSED,
			},
			want: pb.AppStatus_ABNORMAL,
		},
		{
			name: "starting",
			statuses: map[string]string{
				"cat":  v1.RUNNING,
				"dog":  v1.CLOSED,
				"food": v1.STARTING,
			},
			want: pb.AppStatus_STARTING,
		},
		{
			name: "stopping",
			statuses: map[string]string{
				"apple":  v1.STOPPING,
				"banana": v1.CLOSED,
			},
			want: pb.AppStatus_STOPPING,
		},
		{
			name: "stopping2",
			statuses: map[string]string{
				"apple":  v1.STOPPING,
				"banana": v1.CLOSED,
				"cat":    v1.RUNNING,
			},
			want: pb.AppStatus_RUNNING,
		},
		{
			name: "running",
			statuses: map[string]string{
				"apple":  v1.RUNNING,
				"banana": v1.CLOSED,
				"cat":    v1.UNDEPLOY,
			},
			want: pb.AppStatus_RUNNING,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			status := getAppStatus(tc.statuses)
			assert.Equal(t, tc.want, status)
		})
	}
}
