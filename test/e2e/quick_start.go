/*
Copyright 2020 The Kubernetes Authors.

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

package e2e

import (
	"context"
	"fmt"
	infrastructurev1alpha4 "github.com/vmware-tanzu/cluster-api-provider-byoh/api/v1alpha4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// QuickStartSpecInput is the input for QuickStartSpec.
type QuickStartSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	ClusterctlConfigPath  string
	BootstrapClusterProxy framework.ClusterProxy
	ArtifactFolder        string
	SkipCleanup           bool
}

// QuickStartSpec implements a spec that mimics the operation described in the Cluster API quick start, that is
// creating a workload cluster.
// This test is meant to provide a first, fast signal to detect regression; it is recommended to use it as a PR blocker test.
func QuickStartSpec(ctx context.Context, inputGetter func() QuickStartSpecInput) {
	var (
		specName         = "quick-start"
		input            QuickStartSpecInput
		namespace        *corev1.Namespace
		cancelWatches    context.CancelFunc
		clusterResources *clusterctl.ApplyClusterTemplateAndWaitResult
	)

	BeforeEach(func() {
		Expect(ctx).NotTo(BeNil(), "ctx is required for %s spec", specName)
		input = inputGetter()
		Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil when calling %s spec", specName)
		Expect(input.ClusterctlConfigPath).To(BeAnExistingFile(), "Invalid argument. input.ClusterctlConfigPath must be an existing file when calling %s spec", specName)
		Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil when calling %s spec", specName)
		Expect(os.MkdirAll(input.ArtifactFolder, 0755)).To(Succeed(), "Invalid argument. input.ArtifactFolder can't be created for %s spec", specName)

		Expect(input.E2EConfig.Variables).To(HaveKey(KubernetesVersion))

		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace, cancelWatches = setupSpecNamespace(ctx, specName, input.BootstrapClusterProxy, input.ArtifactFolder)
		clusterResources = new(clusterctl.ApplyClusterTemplateAndWaitResult)
	})

	FIt("Should create a workload cluster", func() {
		By("Creating a workload cluster")

		flavor := clusterctl.DefaultFlavor
		if input.E2EConfig.GetVariable(IPFamily) == "IPv6" {
			flavor = "ipv6"
		}

		clusterName := fmt.Sprintf("%s-%s", specName, "happybdayanusha")
		//clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		//	ClusterProxy: input.BootstrapClusterProxy,
		//	ConfigCluster: clusterctl.ConfigClusterInput{
		//		LogFolder:                filepath.Join(input.ArtifactFolder, "clusters", input.BootstrapClusterProxy.GetName()),
		//		ClusterctlConfigPath:     input.ClusterctlConfigPath,
		//		KubeconfigPath:           input.BootstrapClusterProxy.GetKubeconfigPath(),
		//		InfrastructureProvider:   "docker:v0.4.99",
		//		Flavor:                   flavor,
		//		Namespace:                namespace.Name,
		//		ClusterName:              clusterName,
		//		KubernetesVersion:        input.E2EConfig.GetVariable(KubernetesVersion),
		//		ControlPlaneMachineCount: pointer.Int64Ptr(1),
		//		WorkerMachineCount:       pointer.Int64Ptr(0),
		//	},
		//	WaitForClusterIntervals:      input.E2EConfig.GetIntervals(specName, "wait-cluster"),
		//	WaitForControlPlaneIntervals: input.E2EConfig.GetIntervals(specName, "wait-control-plane"),
		//}, clusterResources)

		By("create a ByoHost")
		ByoHost := &infrastructurev1alpha4.ByoHost{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ByoHost",
				APIVersion: "infrastructure.cluster.x-k8s.io/v1alpha4",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jaime.com",
				Namespace: namespace.Name,
			},
			Spec: infrastructurev1alpha4.ByoHostSpec{
				Foo: "Baz",
			},
		}
		client := input.BootstrapClusterProxy.GetClient()
		Expect(client.Create(ctx, ByoHost)).Should(Succeed())

		clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
			ClusterProxy: input.BootstrapClusterProxy,
			ConfigCluster: clusterctl.ConfigClusterInput{
				LogFolder:                filepath.Join(input.ArtifactFolder, "clusters", input.BootstrapClusterProxy.GetName()),
				ClusterctlConfigPath:     input.ClusterctlConfigPath,
				KubeconfigPath:           input.BootstrapClusterProxy.GetKubeconfigPath(),
				InfrastructureProvider:   "byoh:v0.4.0",
				Flavor:                   flavor,
				Namespace:                namespace.Name,
				ClusterName:              clusterName,
				KubernetesVersion:        input.E2EConfig.GetVariable(KubernetesVersion),
				ControlPlaneMachineCount: pointer.Int64Ptr(1),
				WorkerMachineCount:       pointer.Int64Ptr(1),
			},
			WaitForClusterIntervals:      input.E2EConfig.GetIntervals(specName, "wait-cluster"),
			WaitForControlPlaneIntervals: input.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			WaitForMachineDeployments: input.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
		}, clusterResources)

		//
		//clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		//	ClusterProxy: input.BootstrapClusterProxy,
		//	ConfigCluster: clusterctl.ConfigClusterInput{
		//		LogFolder:                filepath.Join(input.ArtifactFolder, "clusters", input.BootstrapClusterProxy.GetName()),
		//		ClusterctlConfigPath:     input.ClusterctlConfigPath,
		//		KubeconfigPath:           input.BootstrapClusterProxy.GetKubeconfigPath(),
		//		InfrastructureProvider:   "byoh:v0.4.0",
		//		Flavor:                   flavor,
		//		Namespace:                namespace.Name,
		//		ClusterName:              fmt.Sprintf("%s-%s", specName, util.RandomString(6)),
		//		KubernetesVersion:        input.E2EConfig.GetVariable(KubernetesVersion),
		//		ControlPlaneMachineCount: pointer.Int64Ptr(0),
		//		WorkerMachineCount:       pointer.Int64Ptr(1),
		//	},
		//	WaitForClusterIntervals:      input.E2EConfig.GetIntervals(specName, "wait-cluster"),
		//	WaitForControlPlaneIntervals: input.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
		//}, clusterResources)

		By("PASSED!")
	})

	AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		dumpSpecResourcesAndCleanup(ctx, specName, input.BootstrapClusterProxy, input.ArtifactFolder, namespace, cancelWatches, clusterResources.Cluster, input.E2EConfig.GetIntervals, input.SkipCleanup)
	})
}
