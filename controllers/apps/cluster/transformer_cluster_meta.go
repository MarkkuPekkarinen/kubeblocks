/*
Copyright (C) 2022-2025 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package cluster

import (
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
)

type clusterMetaTransformer struct{}

var _ graph.Transformer = &clusterMetaTransformer{}

func (t *clusterMetaTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*clusterTransformContext)
	cluster := transCtx.Cluster

	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !controllerutil.ContainsFinalizer(cluster, constant.DBClusterFinalizerName) {
		controllerutil.AddFinalizer(cluster, constant.DBClusterFinalizerName)
	}

	// patch the label to prevent the label from being modified by the user.
	labels := cluster.Labels
	if labels == nil {
		labels = map[string]string{}
	}
	cdLabelName := labels[constant.ClusterDefLabelKey]
	cdName := cluster.Spec.ClusterDef
	if cdLabelName == cdName {
		return nil
	}
	labels[constant.ClusterDefLabelKey] = cdName
	cluster.Labels = labels

	return nil
}
