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

package operator_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ssv1 "github.com/openshift/kueue-operator/pkg/apis/kueueoperator/v1alpha1"
	ssscheme "github.com/openshift/kueue-operator/pkg/generated/clientset/versioned/scheme"
	"github.com/openshift/kueue-operator/test/e2e/bindata"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"
)

var _ = Describe("Kueue Operator Installation", Ordered, func() {
	Context("When installing the operator", func() {
		Expect(deployOperator()).To(Succeed(), "operator deployment should not fail")
		It("Should have running operator pods", func() {
			Eventually(func() error {
				pods, err := f.KubeClient.CoreV1().Pods(operatorNamespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return fmt.Errorf("listing pods: %w", err)
				}

				for _, pod := range pods.Items {
					if strings.HasPrefix(pod.Name, operatorNamespace+"-") &&
						pod.Status.Phase == corev1.PodRunning &&
						pod.DeletionTimestamp == nil {
						return nil
					}
				}
				return fmt.Errorf("no running operator pods found")
			}, operatorReadyTimeout, operatorPollInterval).Should(Succeed())
		})

		It("Should have valid CRDs installed", func() {
			Eventually(func() error {
				crds, err := f.APIExtClient.CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return fmt.Errorf("listing CRDs: %w", err)
				}

				for _, crd := range crds.Items {
					if strings.Contains(crd.Name, "kueue") && crd.Name != "kueues.operator.openshift.io" {
						for _, version := range crd.Spec.Versions {
							if strings.HasPrefix(version.Name, "v1alpha") {
								return fmt.Errorf("unexpected v1alpha CRD version found in %s", crd.Name)
							}
						}
					}
				}
				return nil
			}, operatorReadyTimeout, operatorPollInterval).Should(Succeed())
		})
	})
})

func deployOperator() error {
	kubeClient := f.KubeClient
	apiExtClient := f.APIExtClient
	ssClient := f.KueueClient

	eventRecorder := events.NewKubeRecorder(kubeClient.CoreV1().Events(operatorNamespace), "test-e2e", &corev1.ObjectReference{}, clock.RealClock{})

	ctx, cancelFnc := context.WithCancel(context.TODO())
	defer cancelFnc()

	assets := []struct {
		path           string
		readerAndApply func(objBytes []byte) error
	}{
		{
			path: "assets/00_kueue-operator.crd.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyCustomResourceDefinitionV1(ctx, apiExtClient, eventRecorder, resourceread.ReadCustomResourceDefinitionV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/01_namespace.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyNamespace(ctx, kubeClient.CoreV1(), eventRecorder, resourceread.ReadNamespaceV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/02_clusterrole.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRole(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/03_clusterrolebinding.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRoleBinding(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleBindingV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/04_serviceaccount.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyServiceAccount(ctx, kubeClient.CoreV1(), eventRecorder, resourceread.ReadServiceAccountV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/05_clusterrole_kueue-batch.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRole(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/06_clusterrole_kueue-admin.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRole(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleV1OrDie(objBytes))
				return err
			},
		},

		{
			path: "assets/07_deployment.yaml",
			readerAndApply: func(objBytes []byte) error {
				required := resourceread.ReadDeploymentV1OrDie(objBytes)
				operatorImage := os.Getenv("OPERATOR_IMAGE")
				kueueImage := os.Getenv("KUEUE_IMAGE")

				required.Spec.Template.Spec.Containers[0].Image = operatorImage
				required.Spec.Template.Spec.Containers[0].Env[2].Value = kueueImage
				_, _, err := resourceapply.ApplyDeployment(
					ctx,
					kubeClient.AppsV1(),
					eventRecorder,
					required,
					1000, // any random high number
				)
				return err
			},
		},
		{
			path: "assets/08_kueue_default.yaml",
			readerAndApply: func(objBytes []byte) error {
				requiredObj, err := runtime.Decode(ssscheme.Codecs.UniversalDecoder(ssv1.SchemeGroupVersion), objBytes)
				if err != nil {
					klog.Errorf("Unable to decode assets/08_kueue_default.yaml: %v", err)
					return err
				}
				requiredSS := requiredObj.(*ssv1.Kueue)
				_, err = ssClient.KueueV1alpha1().Kueues(requiredSS.Namespace).Create(ctx, requiredSS, metav1.CreateOptions{})
				return err
			},
		},
	}

	Eventually(func() error {
		for _, asset := range assets {
			klog.Infof("Creating %v", asset.path)
			if err := asset.readerAndApply(bindata.MustAsset(asset.path)); err != nil {
				return err
			}
		}
		return nil
	}, operatorReadyTimeout, operatorPollInterval).Should(Succeed(), "assets should be deployed")

	return nil
}
