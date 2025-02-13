/*
Copyright 2021 The Kubernetes Authors.

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

package topology

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/topology/internal/scope"
	"sigs.k8s.io/cluster-api/internal/testtypes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetBlueprint(t *testing.T) {
	crds := []client.Object{
		testtypes.GenericInfrastructureClusterTemplateCRD,
		testtypes.GenericInfrastructureMachineTemplateCRD,
		testtypes.GenericInfrastructureMachineCRD,
		testtypes.GenericControlPlaneTemplateCRD,
		testtypes.GenericBootstrapConfigTemplateCRD,
	}

	// ignoreResourceVersion is an option to pass to cmpopts to ignore this field which is set by the fakeClient
	// TODO: Make composable version of these options in the builder package to reuse these filters across tests.
	ignoreResourceVersion := cmpopts.IgnoreFields(metav1.ObjectMeta{}, "ResourceVersion")

	// Create objects used across test cases.
	infraClusterTemplate := testtypes.NewInfrastructureClusterTemplateBuilder(metav1.NamespaceDefault, "infraclustertemplate1").
		Build()
	controlPlaneTemplate := testtypes.NewControlPlaneTemplateBuilder(metav1.NamespaceDefault, "controlplanetemplate1").
		Build()

	controlPlaneInfrastructureMachineTemplate := testtypes.NewInfrastructureMachineTemplateBuilder(metav1.NamespaceDefault, "controlplaneinframachinetemplate1").
		Build()
	controlPlaneTemplateWithInfrastructureMachine := testtypes.NewControlPlaneTemplateBuilder(metav1.NamespaceDefault, "controlplanetempaltewithinfrastructuremachine1").
		WithInfrastructureMachineTemplate(controlPlaneInfrastructureMachineTemplate).
		Build()

	workerInfrastructureMachineTemplate := testtypes.NewInfrastructureMachineTemplateBuilder(metav1.NamespaceDefault, "workerinframachinetemplate1").
		Build()
	workerBootstrapTemplate := testtypes.NewBootstrapTemplateBuilder(metav1.NamespaceDefault, "workerbootstraptemplate1").
		Build()
	machineDeployment := testtypes.NewMachineDeploymentClassBuilder(metav1.NamespaceDefault, "machinedeployment1").
		WithClass("workerclass1").
		WithLabels(map[string]string{"foo": "bar"}).
		WithAnnotations(map[string]string{"a": "b"}).
		WithInfrastructureTemplate(workerInfrastructureMachineTemplate).
		WithBootstrapTemplate(workerBootstrapTemplate).
		Build()
	mds := []clusterv1.MachineDeploymentClass{*machineDeployment}

	// Define test cases.
	tests := []struct {
		name         string
		clusterClass *clusterv1.ClusterClass
		objects      []client.Object
		want         *scope.ClusterBlueprint
		wantErr      bool
	}{
		{
			name:    "Fails if ClusterClass does not exist",
			wantErr: true,
		},
		{
			name: "Fails if ClusterClass does not have reference to the InfrastructureClusterTemplate",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "clusterclass1").
				// No InfrastructureClusterTemplate reference!
				Build(),
			wantErr: true,
		},
		{
			name: "Fails if ClusterClass references an InfrastructureClusterTemplate that does not exist",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "clusterclass1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				Build(),
			objects: []client.Object{
				// infraClusterTemplate is missing!
			},
			wantErr: true,
		},
		{
			name: "Fails if ClusterClass does not have reference to the ControlPlaneTemplate",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				// No ControlPlaneTemplate reference!
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
			},
			wantErr: true,
		},
		{
			name: "Fails if ClusterClass does not have reference to the ControlPlaneTemplate",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				// ControlPlaneTemplate is missing!
			},
			wantErr: true,
		},
		{
			name: "Should read a ClusterClass without worker classes",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplate,
			},
			want: &scope.ClusterBlueprint{
				ClusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
					WithInfrastructureClusterTemplate(infraClusterTemplate).
					WithControlPlaneTemplate(controlPlaneTemplate).
					Build(),
				InfrastructureClusterTemplate: infraClusterTemplate,
				ControlPlane: &scope.ControlPlaneBlueprint{
					Template: controlPlaneTemplate,
				},
				MachineDeployments: map[string]*scope.MachineDeploymentBlueprint{},
			},
		},
		{
			name: "Should read a ClusterClass referencing an InfrastructureMachineTemplate for the ControlPlane (but without any worker class)",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplateWithInfrastructureMachine).
				WithControlPlaneInfrastructureMachineTemplate(controlPlaneInfrastructureMachineTemplate).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplateWithInfrastructureMachine,
				controlPlaneInfrastructureMachineTemplate,
			},
			want: &scope.ClusterBlueprint{
				ClusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
					WithInfrastructureClusterTemplate(infraClusterTemplate).
					WithControlPlaneTemplate(controlPlaneTemplateWithInfrastructureMachine).
					WithControlPlaneInfrastructureMachineTemplate(controlPlaneInfrastructureMachineTemplate).
					Build(),
				InfrastructureClusterTemplate: infraClusterTemplate,
				ControlPlane: &scope.ControlPlaneBlueprint{
					Template:                      controlPlaneTemplateWithInfrastructureMachine,
					InfrastructureMachineTemplate: controlPlaneInfrastructureMachineTemplate,
				},
				MachineDeployments: map[string]*scope.MachineDeploymentBlueprint{},
			},
		},
		{
			name: "Fails if ClusterClass references an InfrastructureMachineTemplate for the ControlPlane that does not exist",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				WithControlPlaneInfrastructureMachineTemplate(controlPlaneInfrastructureMachineTemplate).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplate,
				// controlPlaneInfrastructureMachineTemplate is missing!
			},
			wantErr: true,
		},
		{
			name: "Should read a ClusterClass with a MachineDeploymentClass",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				WithWorkerMachineDeploymentClasses(mds).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplate,
				workerInfrastructureMachineTemplate,
				workerBootstrapTemplate,
			},
			want: &scope.ClusterBlueprint{
				ClusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
					WithInfrastructureClusterTemplate(infraClusterTemplate).
					WithControlPlaneTemplate(controlPlaneTemplate).
					WithWorkerMachineDeploymentClasses(mds).
					Build(),
				InfrastructureClusterTemplate: infraClusterTemplate,
				ControlPlane: &scope.ControlPlaneBlueprint{
					Template: controlPlaneTemplate,
				},
				MachineDeployments: map[string]*scope.MachineDeploymentBlueprint{
					"workerclass1": {
						Metadata: clusterv1.ObjectMeta{
							Labels:      map[string]string{"foo": "bar"},
							Annotations: map[string]string{"a": "b"},
						},
						InfrastructureMachineTemplate: workerInfrastructureMachineTemplate,
						BootstrapTemplate:             workerBootstrapTemplate,
					},
				},
			},
		},
		{
			name: "Fails if ClusterClass has a MachineDeploymentClass referencing a BootstrapTemplate that does not exist",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				WithWorkerMachineDeploymentClasses(mds).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplate,
				workerInfrastructureMachineTemplate,
				// workerBootstrapTemplate is missing!
			},
			wantErr: true,
		},
		{
			name: "Fails if ClusterClass has a MachineDeploymentClass referencing a InfrastructureMachineTemplate that does not exist",
			clusterClass: testtypes.NewClusterClassBuilder(metav1.NamespaceDefault, "class1").
				WithInfrastructureClusterTemplate(infraClusterTemplate).
				WithControlPlaneTemplate(controlPlaneTemplate).
				WithWorkerMachineDeploymentClasses(mds).
				Build(),
			objects: []client.Object{
				infraClusterTemplate,
				controlPlaneTemplate,
				workerBootstrapTemplate,
				// workerInfrastructureTemplate is missing!
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			// Set up a cluster using the ClusterClass, if any.
			cluster := testtypes.NewClusterBuilder(metav1.NamespaceDefault, "cluster1").Build()
			if tt.clusterClass != nil {
				cluster.Spec.Topology = &clusterv1.Topology{
					Class: tt.clusterClass.Name,
				}
			} else {
				cluster.Spec.Topology = &clusterv1.Topology{
					Class: "foo",
				}
			}

			// Sets up the fakeClient for the test case.
			objs := []client.Object{}
			objs = append(objs, crds...)
			objs = append(objs, tt.objects...)
			if tt.clusterClass != nil {
				objs = append(objs, tt.clusterClass)
			}
			fakeClient := fake.NewClientBuilder().
				WithScheme(fakeScheme).
				WithObjects(objs...).
				Build()

			// Calls getBlueprint.
			r := &ClusterReconciler{
				Client:                    fakeClient,
				UnstructuredCachingClient: fakeClient,
			}
			got, err := r.getBlueprint(ctx, scope.New(cluster).Current.Cluster)

			// Checks the return error.
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}

			if tt.want == nil {
				g.Expect(got).To(BeNil())
				return
			}

			// Checks the blueprint content.

			// Expect the Diff resulting from each object comparison to be empty when ignoring ObjectMeta.ResourceVersion
			// This is necessary as the FakeClient adds its own ResourceVersion on object creation.
			g.Expect(cmp.Diff(tt.want.ClusterClass, got.ClusterClass, ignoreResourceVersion)).To(Equal(""),
				cmp.Diff(tt.want.ClusterClass, got.ClusterClass, ignoreResourceVersion))
			g.Expect(cmp.Diff(tt.want.InfrastructureClusterTemplate, got.InfrastructureClusterTemplate, ignoreResourceVersion)).To(Equal(""),
				cmp.Diff(tt.want.InfrastructureClusterTemplate, got.InfrastructureClusterTemplate, ignoreResourceVersion))
			g.Expect(cmp.Diff(tt.want.ControlPlane, got.ControlPlane, ignoreResourceVersion)).To(Equal(""),
				cmp.Diff(tt.want.ControlPlane, got.ControlPlane, ignoreResourceVersion))
			g.Expect(cmp.Diff(tt.want.MachineDeployments, got.MachineDeployments, ignoreResourceVersion)).To(Equal(""),
				cmp.Diff(tt.want.MachineDeployments, got.MachineDeployments, ignoreResourceVersion))
		})
	}
}
