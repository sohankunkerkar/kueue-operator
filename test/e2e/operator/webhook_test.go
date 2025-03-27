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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Webhook Configuration", func() {
	Context("When the operator is installed", func() {
		It("Should have valid webhook configurations", func() {
			Eventually(func() error {
				return validateWebhooks()
			}, operatorReadyTimeout, operatorPollInterval).Should(Succeed())
		})
	})
})

func validateWebhooks() error {
	vwh, err := f.KubeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(
		context.TODO(),
		"kueue-validating-webhook-configuration",
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("getting validating webhook: %w", err)
	}

	mwh, err := f.KubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(
		context.TODO(),
		"kueue-mutating-webhook-configuration",
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("getting mutating webhook: %w", err)
	}

	if len(vwh.Webhooks) == 0 || len(mwh.Webhooks) == 0 {
		return fmt.Errorf("webhook configurations are empty")
	}

	return nil
}
