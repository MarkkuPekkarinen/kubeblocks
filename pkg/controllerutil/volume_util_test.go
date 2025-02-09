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

package controllerutil

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	kbappsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	cfgutil "github.com/apecloud/kubeblocks/pkg/configuration/util"
	"github.com/apecloud/kubeblocks/pkg/constant"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

var _ = Describe("lifecycle_utils", func() {

	Context("has the checkAndUpdatePodVolumes function which generates Pod Volumes for mounting ConfigMap objects", func() {
		var sts appsv1.StatefulSet
		var volumes map[string]kbappsv1.ComponentTemplateSpec
		BeforeEach(func() {
			sts = appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "data",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
							},
							Containers: []corev1.Container{
								{
									Name:            "mysql",
									Image:           "docker.io/apecloud/apecloud-mysql-server:latest",
									ImagePullPolicy: "IfNotPresent",
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "data",
											MountPath: "/data",
										},
									},
								},
							},
						},
					},
				},
			}
			volumes = make(map[string]kbappsv1.ComponentTemplateSpec)

		})

		It("should succeed in corner case where input volumes is nil, which means no volume is added", func() {
			ps := &sts.Spec.Template.Spec
			err := CreateOrUpdatePodVolumes(ps, volumes, nil)
			Expect(err).Should(BeNil())
			Expect(len(ps.Volumes)).To(Equal(1))
		})

		It("should succeed in normal test case, where one volume is added", func() {
			volumes["my_config"] = kbappsv1.ComponentTemplateSpec{
				Name:        "myConfig",
				TemplateRef: "myConfig",
				VolumeName:  "myConfigVolume",
			}
			ps := &sts.Spec.Template.Spec
			err := CreateOrUpdatePodVolumes(ps, volumes, nil)
			Expect(err).Should(BeNil())
			Expect(len(ps.Volumes)).To(Equal(2))
		})

		It("should succeed in normal test case, where two volumes are added", func() {
			volumes["my_config"] = kbappsv1.ComponentTemplateSpec{
				Name:        "myConfig",
				TemplateRef: "myConfig",
				VolumeName:  "myConfigVolume",
			}
			volumes["my_config1"] = kbappsv1.ComponentTemplateSpec{
				Name:        "myConfig",
				TemplateRef: "myConfig",
				VolumeName:  "myConfigVolume2",
			}
			ps := &sts.Spec.Template.Spec
			err := CreateOrUpdatePodVolumes(ps, volumes, nil)
			Expect(err).Should(BeNil())
			Expect(len(ps.Volumes)).To(Equal(3))
		})

		It("should fail if updated volume doesn't contain ConfigMap", func() {
			const (
				cmName            = "my_config_for_test"
				replicaVolumeName = "mytest-cm-volume_for_test"
			)
			sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes,
				corev1.Volume{
					Name: replicaVolumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				})
			volumes[cmName] = kbappsv1.ComponentTemplateSpec{
				Name:        "configTplName",
				TemplateRef: "configTplName",
				VolumeName:  replicaVolumeName,
			}
			ps := &sts.Spec.Template.Spec
			Expect(CreateOrUpdatePodVolumes(ps, volumes, nil)).ShouldNot(Succeed())
		})

		It("should succeed if updated volume contains ConfigMap", func() {
			const (
				cmName            = "my_config_for_isv"
				replicaVolumeName = "mytest-cm-volume_for_isv"
			)

			// mock clusterdefinition has volume
			sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes,
				corev1.Volume{
					Name: replicaVolumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "anything"},
						},
					},
				})

			volumes[cmName] = kbappsv1.ComponentTemplateSpec{
				Name:        "configTplName",
				TemplateRef: "configTplName",
				VolumeName:  replicaVolumeName,
			}
			ps := &sts.Spec.Template.Spec
			err := CreateOrUpdatePodVolumes(ps, volumes, nil)
			Expect(err).Should(BeNil())
			Expect(len(sts.Spec.Template.Spec.Volumes)).To(Equal(2))
			volume := GetVolumeMountName(sts.Spec.Template.Spec.Volumes, cmName)
			Expect(volume).ShouldNot(BeNil())
			Expect(volume.ConfigMap).ShouldNot(BeNil())
			Expect(volume.ConfigMap.Name).Should(BeEquivalentTo(cmName))
			Expect(volume.Name).Should(BeEquivalentTo(replicaVolumeName))
		})

	})
})

func Test_buildVolumeMode(t *testing.T) {
	type args struct {
		configs    []string
		configSpec kbappsv1.ComponentTemplateSpec
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{{
		name: "config_test",
		args: args{
			configs: []string{"config1", "config2"},
			configSpec: kbappsv1.ComponentTemplateSpec{
				Name:        "config1",
				DefaultMode: cfgutil.ToPointer(int32(0777)),
			},
		},
		want: cfgutil.ToPointer(int32(0777)),
	}, {
		name: "config_test2",
		args: args{
			configs: []string{"config1", "config2"},
			configSpec: kbappsv1.ComponentTemplateSpec{
				Name: "config1",
			},
		},
		want: cfgutil.ToPointer(configsDefaultMode),
	}, {
		name: "script_test",
		args: args{
			configs: []string{"config1", "config2"},
			configSpec: kbappsv1.ComponentTemplateSpec{
				Name:        "script",
				DefaultMode: cfgutil.ToPointer(int32(0777)),
			},
		},
		want: cfgutil.ToPointer(int32(0777)),
	}, {
		name: "script_test2",
		args: args{
			configs: []string{"config1", "config2"},
			configSpec: kbappsv1.ComponentTemplateSpec{
				Name: "script",
			},
		},
		want: cfgutil.ToPointer(scriptsDefaultMode),
	}}
	viper.Set(constant.FeatureGateIgnoreConfigTemplateDefaultMode, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildVolumeMode(tt.args.configs, tt.args.configSpec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildVolumeMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
