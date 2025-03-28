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

package instanceset

import (
	"fmt"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1"
	"github.com/apecloud/kubeblocks/pkg/controller/builder"
	"github.com/apecloud/kubeblocks/pkg/controller/kubebuilderx"
)

var _ = Describe("replicas alignment reconciler test", func() {
	BeforeEach(func() {
		its = builder.NewInstanceSetBuilder(namespace, name).
			SetReplicas(3).
			SetTemplate(template).
			SetVolumeClaimTemplates(volumeClaimTemplates...).
			SetRoles(roles).
			GetObject()
	})

	Context("PreCondition & Reconcile", func() {
		It("should work well", func() {
			By("PreCondition")
			its.Generation = 1
			tree := kubebuilderx.NewObjectTree()
			tree.SetRoot(its)
			reconciler = NewReplicasAlignmentReconciler()
			Expect(reconciler.PreCondition(tree)).Should(Equal(kubebuilderx.ConditionSatisfied))

			By("prepare current tree")
			// desired: bar-0, bar-1, bar-2, bar-3, bar-foo-0, bar-foo-1, bar-hello-0
			// current: bar-1, bar-foo-0
			replicas := int32(7)
			its.Spec.Replicas = &replicas
			nameHello := "hello"
			instanceHello := workloads.InstanceTemplate{
				Name: nameHello,
			}
			its.Spec.Instances = append(its.Spec.Instances, instanceHello)
			nameFoo := "foo"
			replicasFoo := int32(2)
			instanceFoo := workloads.InstanceTemplate{
				Name:     nameFoo,
				Replicas: &replicasFoo,
			}
			its.Spec.Instances = append(its.Spec.Instances, instanceFoo)
			podFoo0 := builder.NewPodBuilder(namespace, its.Name+"-foo-0").GetObject()
			pvcFoo0 := builder.NewPVCBuilder(namespace, volumeClaimTemplates[0].Name+"-"+podFoo0.Name).GetObject()
			podBar1 := builder.NewPodBuilder(namespace, "bar-1").GetObject()
			pvcBar1 := builder.NewPVCBuilder(namespace, volumeClaimTemplates[0].Name+"-"+podBar1.Name).GetObject()
			Expect(tree.Add(podFoo0, pvcFoo0, podBar1, pvcBar1)).Should(Succeed())

			By("update revisions")
			revisionUpdateReconciler := NewRevisionUpdateReconciler()
			_, err := revisionUpdateReconciler.Reconcile(tree)
			Expect(err).Should(BeNil())

			By("do reconcile with OrderedReady(Serial) policy")
			orderedReadyTree, err := tree.DeepCopy()
			Expect(err).Should(BeNil())
			res, err := reconciler.Reconcile(orderedReadyTree)
			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(kubebuilderx.Continue))
			// desired: bar-0, bar-1, bar-foo-0
			pods := orderedReadyTree.List(&corev1.Pod{})
			Expect(pods).Should(HaveLen(3))
			pvcs := orderedReadyTree.List(&corev1.PersistentVolumeClaim{})
			Expect(pvcs).Should(HaveLen(3))
			podBar0 := builder.NewPodBuilder(namespace, "bar-0").GetObject()
			for _, object := range []client.Object{podFoo0, podBar0, podBar1} {
				Expect(slices.IndexFunc(pods, func(item client.Object) bool {
					return item.GetName() == object.GetName()
				})).Should(BeNumerically(">=", 0))
				Expect(slices.IndexFunc(pvcs, func(item client.Object) bool {
					expectedPVCName := fmt.Sprintf("%s-%s", volumeClaimTemplates[0].Name, object.GetName())
					return expectedPVCName == item.GetName()
				})).Should(BeNumerically(">=", 0))
			}

			By("do reconcile with Parallel policy")
			parallelTree, err := tree.DeepCopy()
			Expect(err).Should(BeNil())
			parallelITS, ok := parallelTree.GetRoot().(*workloads.InstanceSet)
			Expect(ok).Should(BeTrue())
			parallelITS.Spec.PodManagementPolicy = appsv1.ParallelPodManagement
			res, err = reconciler.Reconcile(parallelTree)
			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(kubebuilderx.Continue))
			// desired: bar-0, bar-1, bar-2, bar-3, bar-foo-0, bar-foo-1, bar-hello-0
			pods = parallelTree.List(&corev1.Pod{})
			Expect(pods).Should(HaveLen(7))
			pvcs = parallelTree.List(&corev1.PersistentVolumeClaim{})
			Expect(pvcs).Should(HaveLen(7))
			podHello := builder.NewPodBuilder(namespace, its.Name+"-hello-0").GetObject()
			podFoo1 := builder.NewPodBuilder(namespace, its.Name+"-foo-1").GetObject()
			podBar2 := builder.NewPodBuilder(namespace, "bar-2").GetObject()
			podBar3 := builder.NewPodBuilder(namespace, "bar-3").GetObject()
			for _, object := range []client.Object{podHello, podFoo0, podFoo1, podBar0, podBar1, podBar2, podBar3} {
				Expect(slices.IndexFunc(pods, func(item client.Object) bool {
					return item.GetName() == object.GetName()
				})).Should(BeNumerically(">=", 0))
				Expect(slices.IndexFunc(pvcs, func(item client.Object) bool {
					expectedPVCName := fmt.Sprintf("%s-%s", volumeClaimTemplates[0].Name, object.GetName())
					return expectedPVCName == item.GetName()
				})).Should(BeNumerically(">=", 0))
			}

			By("do reconcile with Parallel policy, ParallelPodManagementConcurrency is 50%")
			parallelTree, err = tree.DeepCopy()
			Expect(err).Should(BeNil())
			parallelITS, ok = parallelTree.GetRoot().(*workloads.InstanceSet)
			Expect(ok).Should(BeTrue())
			parallelITS.Spec.PodManagementPolicy = appsv1.ParallelPodManagement
			parallelITS.Spec.ParallelPodManagementConcurrency = &intstr.IntOrString{Type: intstr.String, StrVal: "50%"}
			res, err = reconciler.Reconcile(parallelTree)
			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(kubebuilderx.Continue))
			// replicas is 7, ParallelPodManagementConcurrency is 50%, so concurrency is 4.
			// since the original bar-1 and bar-foo-0 are not ready, only the new instances bar-0 and bar-2 will be added.
			// desired: bar-0, bar-1, bar-2, bar-foo-0
			pods = parallelTree.List(&corev1.Pod{})
			Expect(pods).Should(HaveLen(4))
			pvcs = parallelTree.List(&corev1.PersistentVolumeClaim{})
			Expect(pvcs).Should(HaveLen(4))
			for _, object := range []client.Object{podFoo0, podBar0, podBar1, podBar2} {
				Expect(slices.IndexFunc(pods, func(item client.Object) bool {
					return item.GetName() == object.GetName()
				})).Should(BeNumerically(">=", 0))
				Expect(slices.IndexFunc(pvcs, func(item client.Object) bool {
					expectedPVCName := fmt.Sprintf("%s-%s", volumeClaimTemplates[0].Name, object.GetName())
					return expectedPVCName == item.GetName()
				})).Should(BeNumerically(">=", 0))
			}
		})

		It("handles nodeSelectorOnce Annotation", func() {
			tree := kubebuilderx.NewObjectTree()
			tree.SetRoot(its)
			name := "bar-1"
			node := "test-1"
			Expect(MergeNodeSelectorOnceAnnotation(its, map[string]string{name: node})).To(Succeed())

			res, err := reconciler.Reconcile(tree)
			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(kubebuilderx.Continue))
			pods := tree.List(&corev1.Pod{})
			for _, obj := range pods {
				pod := obj.(*corev1.Pod)
				if pod.Name == name {
					Expect(pod.Spec.NodeSelector).To(Equal(map[string]string{
						corev1.LabelHostname: node,
					}))
				}
			}
		})
	})
})
