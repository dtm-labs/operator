package controllers

import (
	dtmappv1 "github.com/dtm-labs/operator/api/v1"
	"gopkg.in/yaml.v2"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	AppName       = "dtm"
	ConfigMapName = "dtm-conf"
	ServiceName   = "dtm-svc"
)

var AppProtocolHTTP = "HTTP"
var AppProtocolGRPC = "GRPC"

type (
	Store struct {
		Driver string `yaml:"Driver"`
	}
)

type AppConfig struct {
	Store Store `yaml:"Store"`
}

func (r *DtmReconciler) GetDtmDeployment(dtm *dtmappv1.Dtm) *appv1.Deployment {
	replicas := dtm.Spec.Replicas
	version := dtm.Spec.Version

	deploy := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppName,
			Namespace: dtm.Namespace,
			Labels: map[string]string{
				"app":     AppName,
				"version": version,
			},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     AppName,
					"version": version,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     AppName,
						"version": version,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            AppName,
						Image:           "yedf/dtm:" + version,
						ImagePullPolicy: corev1.PullIfNotPresent,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "config",
							MountPath: "/app/dtm/configs",
						}},
						Args: []string{
							"-c=/app/dtm/configs/config.yaml",
						},
						Ports: []corev1.ContainerPort{{
							Name:          "http",
							ContainerPort: 36789,
							Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "grpc",
							ContainerPort: 36790,
							Protocol:      corev1.ProtocolTCP,
						}},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/api/ping",
									Port:   intstr.FromInt(36789),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/api/ping",
									Port:   intstr.FromInt(36789),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
						},
					}},
					Volumes: []corev1.Volume{{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: ConfigMapName,
								},
							},
						},
					}},
				},
			},
		},
	}

	_ = controllerutil.SetControllerReference(dtm, deploy, r.Scheme)

	return deploy
}

func (r *DtmReconciler) GetDtmConfigMap(dtm *dtmappv1.Dtm) *corev1.ConfigMap {

	// todo
	appConfig := AppConfig{
		Store: Store{
			Driver: "boltdb",
		},
	}

	data, _ := yaml.Marshal(appConfig)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: dtm.Namespace,
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Data: map[string]string{
			"config.yaml": string(data),
		},
	}
	_ = controllerutil.SetControllerReference(dtm, cm, r.Scheme)
	return cm
}

func (r *DtmReconciler) GetDtmService(dtm *dtmappv1.Dtm) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName,
			Namespace: dtm.Namespace,
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:        "http",
				Protocol:    corev1.ProtocolTCP,
				AppProtocol: &AppProtocolHTTP,
				Port:        36789,
				TargetPort:  intstr.FromInt(36789),
			}, {
				Name:        "grpc",
				Protocol:    corev1.ProtocolTCP,
				AppProtocol: &AppProtocolGRPC,
				Port:        36790,
				TargetPort:  intstr.FromInt(36790),
			}},
			Selector: map[string]string{
				"app": AppName,
			},
		},
	}

	_ = controllerutil.SetControllerReference(dtm, svc, r.Scheme)

	return svc
}
