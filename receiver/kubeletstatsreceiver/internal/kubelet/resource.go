// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
)

func fillNodeResource(dest pcommon.Resource, s stats.NodeStats) {
	dest.Attributes().UpsertString(conventions.AttributeK8SNodeName, s.NodeName)
}

func fillPodResource(dest pcommon.Resource, s stats.PodStats) {
	dest.Attributes().UpsertString(conventions.AttributeK8SPodUID, s.PodRef.UID)
	dest.Attributes().UpsertString(conventions.AttributeK8SPodName, s.PodRef.Name)
	dest.Attributes().UpsertString(conventions.AttributeK8SNamespaceName, s.PodRef.Namespace)
}

func fillContainerResource(dest pcommon.Resource, sPod stats.PodStats, sContainer stats.ContainerStats, metadata Metadata) error {
	labels := map[string]string{
		conventions.AttributeK8SPodUID:        sPod.PodRef.UID,
		conventions.AttributeK8SPodName:       sPod.PodRef.Name,
		conventions.AttributeK8SNamespaceName: sPod.PodRef.Namespace,
		conventions.AttributeK8SContainerName: sContainer.Name,
	}
	if err := metadata.setExtraLabels(labels, sPod.PodRef, MetadataLabelContainerID, sContainer.Name); err != nil {
		return fmt.Errorf("failed to set extra labels from metadata: %w", err)
	}
	for k, v := range labels {
		dest.Attributes().UpsertString(k, v)
	}
	return nil
}

func fillVolumeResource(dest pcommon.Resource, sPod stats.PodStats, vs stats.VolumeStats, metadata Metadata) error {
	labels := map[string]string{
		conventions.AttributeK8SPodUID:        sPod.PodRef.UID,
		conventions.AttributeK8SPodName:       sPod.PodRef.Name,
		conventions.AttributeK8SNamespaceName: sPod.PodRef.Namespace,
		labelVolumeName:                       vs.Name,
	}

	if err := metadata.setExtraLabels(labels, sPod.PodRef, MetadataLabelVolumeType, vs.Name); err != nil {
		return fmt.Errorf("failed to set extra labels from metadata: %w", err)
	}

	for k, v := range labels {
		dest.Attributes().UpsertString(k, v)
	}
	return nil
}
