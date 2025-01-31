/*
Copyright 2024.

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

package webhook

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	kueue "github.com/openshift/kueue-operator/pkg/apis/kueueoperator/v1alpha1"
)

// These are webhooks that we only want to enable if the integration requests it
// Otherwise we will remove them from the kueue manifest
func dangerousWebhooks() []string {
	return []string{"vdeployment.kb.io", "vpod.kb.io", "vstatefulset.kb.io"}
}

func ModifyPodBasedValidatingWebhook(kueueCfg kueue.KueueConfiguration, currentWebhook *admissionregistrationv1.ValidatingWebhookConfiguration) *admissionregistrationv1.ValidatingWebhookConfiguration {
	// if there is a pod based webhook, we need to safely enabled it
	// For now we will not modify this
	for _, val := range kueueCfg.Integrations.Frameworks {
		if val == "pod" || val == "deployment" || val == "statefulset" {
			return currentWebhook
		}
	}
	newWebHook := currentWebhook.DeepCopy()
	newWebHook.Webhooks = []admissionregistrationv1.ValidatingWebhook{}
	webhooksToSearchFor := dangerousWebhooks()

	for _, val := range currentWebhook.Webhooks {
		removeWebhook := false
		for _, name := range webhooksToSearchFor {
			if val.Name == name {
				removeWebhook = true
			}
		}
		if !removeWebhook {
			newWebHook.Webhooks = append(newWebHook.Webhooks, val)
		}
	}
	return newWebHook

}

func ModifyPodBasedMutatingWebhook(kueueCfg kueue.KueueConfiguration, currentWebhook *admissionregistrationv1.MutatingWebhookConfiguration) *admissionregistrationv1.MutatingWebhookConfiguration {
	// if there is a pod based webhook, we need to safely enabled it
	// For now we will not modify this
	for _, val := range kueueCfg.Integrations.Frameworks {
		if val == "pod" || val == "deployment" || val == "statefulset" {
			return currentWebhook
		}
	}
	newWebHook := currentWebhook.DeepCopy()
	newWebHook.Webhooks = []admissionregistrationv1.MutatingWebhook{}
	// Remove dangerous webhooks for now
	webhooksToSearchFor := dangerousWebhooks()

	for _, val := range currentWebhook.Webhooks {
		removeWebhook := false
		for _, name := range webhooksToSearchFor {
			if val.Name == name {
				removeWebhook = true
			}
		}
		if !removeWebhook {
			newWebHook.Webhooks = append(newWebHook.Webhooks, val)
		}
	}
	return newWebHook
}
