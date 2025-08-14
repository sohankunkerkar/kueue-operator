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

package configmap

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	configapi "sigs.k8s.io/kueue/apis/config/v1beta1"

	kueue "github.com/openshift/kueue-operator/pkg/apis/kueueoperator/v1"
)

func TestBuildConfigMap(t *testing.T) {
	testCases := map[string]struct {
		configuration kueue.KueueConfiguration
		wantCfgMap    *corev1.ConfigMap
		wantErr       error
	}{
		"batch job example": {
			configuration: kueue.KueueConfiguration{
				Integrations: kueue.Integrations{
					Frameworks: []kueue.KueueIntegration{kueue.KueueIntegrationBatchJob},
				},
				Resources: &configapi.Resources{
					ExcludeResourcePrefixes: []string{"example.com/exclude"},
					Transformations: []configapi.ResourceTransformation{
						{
							Input:    "example.com/gpu-type1",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Replace}[0],
							Outputs: corev1.ResourceList{
								"example.com/gpu-memory": resource.MustParse("5Gi"),
								"example.com/credits":    resource.MustParse("10"),
							},
						},
						{
							Input:    "cpu",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Retain}[0],
							Outputs: corev1.ResourceList{
								"example.com/credits": resource.MustParse("1"),
							},
						},
					},
				},
			},
			wantCfgMap: &corev1.ConfigMap{
				Data: map[string]string{
					"controller_manager_config.yaml": `apiVersion: config.kueue.x-k8s.io/v1beta1
clientConnection:
  burst: 100
  qps: 50
controller:
  groupKindConcurrency:
    ClusterQueue.kueue.x-k8s.io: 1
    Job.batch: 5
    LocalQueue.kueue.x-k8s.io: 1
    Pod: 5
    ResourceFlavor.kueue.x-k8s.io: 1
    Workload.kueue.x-k8s.io: 5
fairSharing:
  enable: false
featureGates:
  HierarchicalCohorts: false
  VisibilityOnDemand: false
health:
  healthProbeBindAddress: :8081
integrations:
  frameworks:
  - batch/job
internalCertManagement:
  enable: false
kind: Configuration
leaderElection:
  leaderElect: true
  leaseDuration: 2m17s
  renewDeadline: 1m47s
  resourceLock: ""
  resourceName: ""
  resourceNamespace: ""
  retryPeriod: 26s
manageJobsWithoutQueueName: false
managedJobsNamespaceSelector:
  matchLabels:
    kueue.openshift.io/managed: "true"
metrics:
  bindAddress: :8443
  enableClusterQueueResources: true
resources:
  excludeResourcePrefixes:
  - example.com/exclude
  transformations:
  - input: example.com/gpu-type1
    outputs:
      example.com/credits: "10"
      example.com/gpu-memory: 5Gi
    strategy: Replace
  - input: cpu
    outputs:
      example.com/credits: "1"
    strategy: Retain
waitForPodsReady: {}
webhook:
  port: 9443
`,
				},
			},
			wantErr: nil,
		},
		"rhoai example": {
			configuration: kueue.KueueConfiguration{
				Integrations: kueue.Integrations{
					Frameworks: []kueue.KueueIntegration{kueue.KueueIntegrationRayJob, kueue.KueueIntegrationRayCluster, kueue.KueueIntegrationPyTorchJob},
				},
				GangScheduling: kueue.GangScheduling{
					Policy: kueue.GangSchedulingPolicyByWorkload,
					ByWorkload: &kueue.ByWorkload{
						Admission: kueue.GangSchedulingWorkloadAdmissionParallel,
					},
				},
				Preemption: kueue.Preemption{PreemptionPolicy: kueue.PreemptionStrategyClassical},
				Resources: &configapi.Resources{
					ExcludeResourcePrefixes: []string{"nvidia.com/exclude"},
					Transformations: []configapi.ResourceTransformation{
						{
							Input:    "nvidia.com/gpu",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Replace}[0],
							Outputs: corev1.ResourceList{
								"nvidia.com/gpu-memory": resource.MustParse("16Gi"),
								"example.com/credits":   resource.MustParse("20"),
							},
						},
					},
				},
			},
			wantCfgMap: &corev1.ConfigMap{
				Data: map[string]string{
					"controller_manager_config.yaml": `apiVersion: config.kueue.x-k8s.io/v1beta1
clientConnection:
  burst: 100
  qps: 50
controller:
  groupKindConcurrency:
    ClusterQueue.kueue.x-k8s.io: 1
    Job.batch: 5
    LocalQueue.kueue.x-k8s.io: 1
    Pod: 5
    ResourceFlavor.kueue.x-k8s.io: 1
    Workload.kueue.x-k8s.io: 5
fairSharing:
  enable: false
featureGates:
  HierarchicalCohorts: false
  VisibilityOnDemand: false
health:
  healthProbeBindAddress: :8081
integrations:
  frameworks:
  - ray.io/rayjob
  - ray.io/raycluster
  - kubeflow.org/pytorchjob
internalCertManagement:
  enable: false
kind: Configuration
leaderElection:
  leaderElect: true
  leaseDuration: 2m17s
  renewDeadline: 1m47s
  resourceLock: ""
  resourceName: ""
  resourceNamespace: ""
  retryPeriod: 26s
manageJobsWithoutQueueName: false
managedJobsNamespaceSelector:
  matchLabels:
    kueue.openshift.io/managed: "true"
metrics:
  bindAddress: :8443
  enableClusterQueueResources: true
resources:
  excludeResourcePrefixes:
  - nvidia.com/exclude
  transformations:
  - input: nvidia.com/gpu
    outputs:
      example.com/credits: "20"
      nvidia.com/gpu-memory: 16Gi
    strategy: Replace
waitForPodsReady:
  blockAdmission: false
  enable: true
webhook:
  port: 9443
`,
				},
			},
			wantErr: nil,
		},
		"ibm example": {
			configuration: kueue.KueueConfiguration{
				Integrations: kueue.Integrations{
					Frameworks: []kueue.KueueIntegration{kueue.KueueIntegrationAppWrapper},
				},
				GangScheduling:     kueue.GangScheduling{Policy: kueue.GangSchedulingPolicyNone},
				WorkloadManagement: kueue.WorkloadManagement{LabelPolicy: kueue.LabelPolicyNone},
				Preemption:         kueue.Preemption{PreemptionPolicy: kueue.PreemptionStrategyFairsharing},
				Resources: &configapi.Resources{
					Transformations: []configapi.ResourceTransformation{
						{
							Input:    "codeflare.dev/appwrapper-cpu",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Retain}[0],
							Outputs: corev1.ResourceList{
								"example.com/compute-units": resource.MustParse("2"),
							},
						},
					},
				},
			},
			wantCfgMap: &corev1.ConfigMap{
				Data: map[string]string{
					"controller_manager_config.yaml": `apiVersion: config.kueue.x-k8s.io/v1beta1
clientConnection:
  burst: 100
  qps: 50
controller:
  groupKindConcurrency:
    ClusterQueue.kueue.x-k8s.io: 1
    Job.batch: 5
    LocalQueue.kueue.x-k8s.io: 1
    Pod: 5
    ResourceFlavor.kueue.x-k8s.io: 1
    Workload.kueue.x-k8s.io: 5
fairSharing:
  enable: true
  preemptionStrategies:
  - LessThanOrEqualToFinalShare
  - LessThanInitialShare
featureGates:
  HierarchicalCohorts: false
  VisibilityOnDemand: false
health:
  healthProbeBindAddress: :8081
integrations:
  frameworks:
  - workload.codeflare.dev/appwrapper
internalCertManagement:
  enable: false
kind: Configuration
leaderElection:
  leaderElect: true
  leaseDuration: 2m17s
  renewDeadline: 1m47s
  resourceLock: ""
  resourceName: ""
  resourceNamespace: ""
  retryPeriod: 26s
manageJobsWithoutQueueName: true
managedJobsNamespaceSelector:
  matchLabels:
    kueue.openshift.io/managed: "true"
metrics:
  bindAddress: :8443
  enableClusterQueueResources: true
resources:
  transformations:
  - input: codeflare.dev/appwrapper-cpu
    outputs:
      example.com/compute-units: "2"
    strategy: Retain
waitForPodsReady: {}
webhook:
  port: 9443
`,
				},
			},
			wantErr: nil,
		},

		"serving workloads": {
			configuration: kueue.KueueConfiguration{
				Integrations: kueue.Integrations{
					Frameworks: []kueue.KueueIntegration{kueue.KueueIntegrationDeployment, kueue.KueueIntegrationPod, kueue.KueueIntegrationStatefulSet, kueue.KueueIntegrationAppWrapper, kueue.KueueIntegrationLeaderWorkerSet},
				},
				Resources: &configapi.Resources{
					ExcludeResourcePrefixes: []string{"ephemeral-storage"},
					Transformations: []configapi.ResourceTransformation{
						{
							Input:    "memory",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Retain}[0],
							Outputs: corev1.ResourceList{
								"example.com/memory-credits": resource.MustParse("1"),
							},
						},
					},
				},
			},
			wantCfgMap: &corev1.ConfigMap{
				Data: map[string]string{
					"controller_manager_config.yaml": `apiVersion: config.kueue.x-k8s.io/v1beta1
clientConnection:
  burst: 100
  qps: 50
controller:
  groupKindConcurrency:
    ClusterQueue.kueue.x-k8s.io: 1
    Job.batch: 5
    LocalQueue.kueue.x-k8s.io: 1
    Pod: 5
    ResourceFlavor.kueue.x-k8s.io: 1
    Workload.kueue.x-k8s.io: 5
fairSharing:
  enable: false
featureGates:
  HierarchicalCohorts: false
  VisibilityOnDemand: false
health:
  healthProbeBindAddress: :8081
integrations:
  frameworks:
  - deployment
  - pod
  - statefulset
  - workload.codeflare.dev/appwrapper
  - leaderworkerset.x-k8s.io/leaderworkerset
internalCertManagement:
  enable: false
kind: Configuration
leaderElection:
  leaderElect: true
  leaseDuration: 2m17s
  renewDeadline: 1m47s
  resourceLock: ""
  resourceName: ""
  resourceNamespace: ""
  retryPeriod: 26s
manageJobsWithoutQueueName: false
managedJobsNamespaceSelector:
  matchLabels:
    kueue.openshift.io/managed: "true"
metrics:
  bindAddress: :8443
  enableClusterQueueResources: true
resources:
  excludeResourcePrefixes:
  - ephemeral-storage
  transformations:
  - input: memory
    outputs:
      example.com/memory-credits: "1"
    strategy: Retain
waitForPodsReady: {}
webhook:
  port: 9443
`,
				},
			},
			wantErr: nil,
		},
		"sequential gang admission": {
			configuration: kueue.KueueConfiguration{
				Integrations: kueue.Integrations{
					Frameworks: []kueue.KueueIntegration{kueue.KueueIntegrationDeployment, kueue.KueueIntegrationPod, kueue.KueueIntegrationStatefulSet, kueue.KueueIntegrationAppWrapper, kueue.KueueIntegrationLeaderWorkerSet},
				},
				GangScheduling: kueue.GangScheduling{
					Policy: kueue.GangSchedulingPolicyByWorkload,
					ByWorkload: &kueue.ByWorkload{
						Admission: kueue.GangSchedulingWorkloadAdmissionSequential,
					},
				},
				Resources: &configapi.Resources{
					Transformations: []configapi.ResourceTransformation{
						{
							Input:    "cpu",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Replace}[0],
							Outputs: corev1.ResourceList{
								"example.com/cpu-credits": resource.MustParse("5"),
							},
						},
						{
							Input:    "memory",
							Strategy: &[]configapi.ResourceTransformationStrategy{configapi.Replace}[0],
							Outputs: corev1.ResourceList{
								"example.com/memory-credits": resource.MustParse("2"),
							},
						},
					},
				},
			},
			wantCfgMap: &corev1.ConfigMap{
				Data: map[string]string{
					"controller_manager_config.yaml": `apiVersion: config.kueue.x-k8s.io/v1beta1
clientConnection:
  burst: 100
  qps: 50
controller:
  groupKindConcurrency:
    ClusterQueue.kueue.x-k8s.io: 1
    Job.batch: 5
    LocalQueue.kueue.x-k8s.io: 1
    Pod: 5
    ResourceFlavor.kueue.x-k8s.io: 1
    Workload.kueue.x-k8s.io: 5
fairSharing:
  enable: false
featureGates:
  HierarchicalCohorts: false
  VisibilityOnDemand: false
health:
  healthProbeBindAddress: :8081
integrations:
  frameworks:
  - deployment
  - pod
  - statefulset
  - workload.codeflare.dev/appwrapper
  - leaderworkerset.x-k8s.io/leaderworkerset
internalCertManagement:
  enable: false
kind: Configuration
leaderElection:
  leaderElect: true
  leaseDuration: 2m17s
  renewDeadline: 1m47s
  resourceLock: ""
  resourceName: ""
  resourceNamespace: ""
  retryPeriod: 26s
manageJobsWithoutQueueName: false
managedJobsNamespaceSelector:
  matchLabels:
    kueue.openshift.io/managed: "true"
metrics:
  bindAddress: :8443
  enableClusterQueueResources: true
resources:
  transformations:
  - input: cpu
    outputs:
      example.com/cpu-credits: "5"
    strategy: Replace
  - input: memory
    outputs:
      example.com/memory-credits: "2"
    strategy: Replace
waitForPodsReady:
  blockAdmission: true
  enable: true
webhook:
  port: 9443
`,
				},
			},
			wantErr: nil,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			got, err := BuildConfigMap("test", tc.configuration)
			if diff := cmp.Diff(got.Data["controller_manager_config.yaml"], tc.wantCfgMap.Data["controller_manager_config.yaml"]); len(diff) != 0 {
				t.Errorf("Unexpected buckets (-want,+got):\n%s", diff)
			}
			if err != nil && tc.wantErr == nil {
				t.Errorf("Unexpected error: want=%v, got=%v", tc.wantErr, err)
			}
		})
	}
}
