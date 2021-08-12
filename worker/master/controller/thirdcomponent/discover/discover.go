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

package discover

import (
	"context"
	"fmt"
	"time"

	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	katolistersv1alpha1 "github.com/gridworkz/kato/pkg/generated/listers/kato/v1alpha1"
	"github.com/gridworkz/kato/worker/master/controller/thirdcomponent/prober"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Discover -
type Discover interface {
	GetComponent() *v1alpha1.ThirdComponent
	DiscoverOne(ctx context.Context) ([]*v1alpha1.ThirdComponentEndpointStatus, error)
	Discover(ctx context.Context, update chan *v1alpha1.ThirdComponent) ([]*v1alpha1.ThirdComponentEndpointStatus, error)
	SetProberManager(proberManager prober.Manager)
}

// NewDiscover -
func NewDiscover(component *v1alpha1.ThirdComponent,
	restConfig *rest.Config,
	lister katolistersv1alpha1.ThirdComponentLister) (Discover, error) {
	if component.Spec.EndpointSource.KubernetesService != nil {
		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			logrus.Errorf("create kube client error: %s", err.Error())
			return nil, err
		}
		return &kubernetesDiscover{
			component: component,
			client:    clientset,
		}, nil
	}
	if len(component.Spec.EndpointSource.StaticEndpoints) > 0 {
		return &staticEndpoint{
			component: component,
			lister:    lister,
		}, nil
	}
	return nil, fmt.Errorf("not support source type")
}

type kubernetesDiscover struct {
	component *v1alpha1.ThirdComponent
	client    *kubernetes.Clientset
}

func (k *kubernetesDiscover) GetComponent() *v1alpha1.ThirdComponent {
	return k.component
}
func (k *kubernetesDiscover) getNamespace() string {
	component := k.component
	namespace := component.Spec.EndpointSource.KubernetesService.Namespace
	if namespace == "" {
		namespace = component.Namespace
	}
	return namespace
}
func (k *kubernetesDiscover) Discover(ctx context.Context, update chan *v1alpha1.ThirdComponent) ([]*v1alpha1.ThirdComponentEndpointStatus, error) {
	namespace := k.getNamespace()
	component := k.component
	service, err := k.client.CoreV1().Services(namespace).Get(ctx, component.Spec.EndpointSource.KubernetesService.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("load kubernetes service failure %s", err.Error())
	}
	re, err := k.client.CoreV1().Endpoints(namespace).Watch(ctx, metav1.ListOptions{LabelSelector: labels.FormatLabels(service.Spec.Selector)})
	if err != nil {
		return nil, fmt.Errorf("watch kubernetes endpoints failure %s", err.Error())
	}
	defer re.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-re.ResultChan():
			func() {
				ctx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()
				endpoints, err := k.DiscoverOne(ctx)
				if err == nil {
					new := component.DeepCopy()
					new.Status.Endpoints = endpoints
					update <- new
				} else {
					logrus.Errorf("discover kubernetes endpoints %s change failure %s", component.Spec.EndpointSource.KubernetesService.Name, err.Error())
				}
			}()
			return k.DiscoverOne(ctx)
		}
	}
}
func (k *kubernetesDiscover) DiscoverOne(ctx context.Context) ([]*v1alpha1.ThirdComponentEndpointStatus, error) {
	component := k.component
	namespace := k.getNamespace()
	service, err := k.client.CoreV1().Services(namespace).Get(ctx, component.Spec.EndpointSource.KubernetesService.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("load kubernetes service failure %s", err.Error())
	}
	// service name must be same with endpoint name
	endpoint, err := k.client.CoreV1().Endpoints(namespace).Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("load kubernetes endpoints failure %s", err.Error())
	}
	getServicePort := func(portName string) int {
		for _, port := range service.Spec.Ports {
			if port.Name == portName {
				return int(port.Port)
			}
		}
		return 0
	}
	var es = []*v1alpha1.ThirdComponentEndpointStatus{}
	for _, subset := range endpoint.Subsets {
		for _, port := range subset.Ports {
			for _, address := range subset.Addresses {
				ed := v1alpha1.NewEndpointAddress(address.IP, int(port.Port))
				if ed != nil {
					es = append(es, &v1alpha1.ThirdComponentEndpointStatus{
						ServicePort: getServicePort(port.Name),
						Address:     *ed,
						TargetRef:   address.TargetRef,
						Status:      v1alpha1.EndpointReady,
					})
				}
			}
			for _, address := range subset.NotReadyAddresses {
				ed := v1alpha1.NewEndpointAddress(address.IP, int(port.Port))
				if ed != nil {
					es = append(es, &v1alpha1.ThirdComponentEndpointStatus{
						Address:     *ed,
						ServicePort: getServicePort(port.Name),
						TargetRef:   address.TargetRef,
						Status:      v1alpha1.EndpointReady,
					})
				}
			}
		}
	}
	return es, nil
}

func (k *kubernetesDiscover) SetProberManager(proberManager prober.Manager) {

}
