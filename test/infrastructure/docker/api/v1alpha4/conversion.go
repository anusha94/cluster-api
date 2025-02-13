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

package v1alpha4

import (
	"sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *DockerCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.DockerCluster)

	return Convert_v1alpha4_DockerCluster_To_v1beta1_DockerCluster(src, dst, nil)
}

func (dst *DockerCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.DockerCluster)

	return Convert_v1beta1_DockerCluster_To_v1alpha4_DockerCluster(src, dst, nil)
}

func (src *DockerClusterTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.DockerClusterTemplate)

	return Convert_v1alpha4_DockerClusterTemplate_To_v1beta1_DockerClusterTemplate(src, dst, nil)
}

func (dst *DockerClusterTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.DockerClusterTemplate)

	return Convert_v1beta1_DockerClusterTemplate_To_v1alpha4_DockerClusterTemplate(src, dst, nil)
}

func (src *DockerMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.DockerMachine)

	return Convert_v1alpha4_DockerMachine_To_v1beta1_DockerMachine(src, dst, nil)
}

func (dst *DockerMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.DockerMachine)

	return Convert_v1beta1_DockerMachine_To_v1alpha4_DockerMachine(src, dst, nil)
}

func (src *DockerMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.DockerMachineTemplate)

	return Convert_v1alpha4_DockerMachineTemplate_To_v1beta1_DockerMachineTemplate(src, dst, nil)
}

func (dst *DockerMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.DockerMachineTemplate)

	return Convert_v1beta1_DockerMachineTemplate_To_v1alpha4_DockerMachineTemplate(src, dst, nil)
}
