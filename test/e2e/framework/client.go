/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"os"

	kueueclient "github.com/openshift/kueue-operator/pkg/generated/clientset/versioned"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	upstreamkueueclient "sigs.k8s.io/kueue/client-go/clientset/versioned"
)

type Framework struct {
	KubeClient          kubernetes.Interface
	APIExtClient        apiextv1.ApiextensionsV1Interface
	KueueClient         kueueclient.Interface
	UpstreamKueueClient upstreamkueueclient.Interface
	DynamicClient       dynamic.Interface
}

func New(kubeconfigPath string) *Framework {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %v", err)
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating kubernetes client: %v", err)
		os.Exit(1)
	}

	apiExtClient, err := apiextv1.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating API extension client: %v", err)
		os.Exit(1)
	}

	kueueClient, err := kueueclient.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating Kueue client: %v", err)
		os.Exit(1)
	}

	upstreamKueueClient, err := upstreamkueueclient.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating upstream Kueue client: %v", err)
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating dynamic client: %v", err)
		os.Exit(1)
	}

	return &Framework{
		KubeClient:          kubeClient,
		APIExtClient:        apiExtClient,
		KueueClient:         kueueClient,
		UpstreamKueueClient: upstreamKueueClient,
		DynamicClient:       dynamicClient,
	}
}
