package k8s_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
)

var _ = Describe("Mapper", func() {

	var lrp *opi.LRP

	BeforeEach(func() {
		statefulset := appsv1.StatefulSet{
			ObjectMeta: meta.ObjectMeta{
				Name: "baldur",
				Annotations: map[string]string{
					cf.ProcessGUID:   "Baldur-guid",
					cf.LastUpdated:   "last-updated-some-time-ago",
					cf.VcapAppUris:   "my.example.route",
					cf.VcapAppID:     "guid_1234",
					cf.VcapVersion:   "version_1234",
					cf.VcapAppName:   "Baldur",
					cf.VcapSpaceName: "space-foo",
				},
			},
			Spec: appsv1.StatefulSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: "busybox",
								Command: []string{
									"/bin/sh",
									"-c",
									"while true; do echo hello; sleep 10;done",
								},
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 8888,
									},
									{
										ContainerPort: 9999,
									},
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceMemory: *resource.NewScaledQuantity(1024, resource.Mega),
									},
									Limits: corev1.ResourceList{
										corev1.ResourceEphemeralStorage: *resource.NewScaledQuantity(2048, resource.Mega),
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "some-claim",
										MountPath: "/some/path",
									},
								},
							},
						},
					},
				},
			},
			Status: appsv1.StatefulSetStatus{
				ReadyReplicas: 3,
			},
		}
		lrp = StatefulSetToLRP(statefulset)
	})

	It("should set the correct LRP identifier", func() {
		Expect(lrp.LRPIdentifier.GUID).To(Equal("guid_1234"))
		Expect(lrp.LRPIdentifier.Version).To(Equal("version_1234"))
	})

	It("should set the correct LRP app name", func() {
		Expect(lrp.AppName).To(Equal("Baldur"))
	})

	It("should set the correct LRP space name", func() {
		Expect(lrp.SpaceName).To(Equal("space-foo"))
	})

	It("should set the correct LRP image", func() {
		Expect(lrp.Image).To(Equal("busybox"))
	})

	It("should set the correct LRP command", func() {
		Expect(lrp.Command).To(Equal([]string{"/bin/sh", "-c", "while true; do echo hello; sleep 10;done"}))
	})

	It("should set the correct LRP running instances", func() {
		Expect(lrp.RunningInstances).To(Equal(3))
	})

	It("should set the correct LRP ports", func() {
		Expect(lrp.Ports).To(Equal([]int32{8888, 9999}))
	})

	It("should translate statefulset annotations to LRP metadata", func() {
		Expect(lrp.Metadata).To(Equal(map[string]string{
			cf.ProcessGUID: "Baldur-guid",
			cf.LastUpdated: "last-updated-some-time-ago",
			cf.VcapAppUris: "my.example.route",
			cf.VcapAppID:   "guid_1234",
			cf.VcapVersion: "version_1234",
			cf.VcapAppName: "Baldur",
		}))
	})

	It("should set the correct LRP disk and memory usage", func() {
		Expect(lrp.MemoryMB).To(Equal(int64(1024)))
		Expect(lrp.DiskMB).To(Equal(int64(2048)))
	})

	It("should set the correct LRP volume mounts", func() {
		Expect(lrp.VolumeMounts).To(Equal([]opi.VolumeMount{
			{
				ClaimName: "some-claim",
				MountPath: "/some/path",
			},
		}))
	})
})
