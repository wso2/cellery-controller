/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package resources

import (
	"fmt"
	"strings"

	"cellery.io/cellery-controller/pkg/crypto"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeDeployment(gateway *v1alpha2.Gateway, cfg config.Interface) (*appsv1.Deployment, error) {
	// if gateway.Spec.Type == v1alpha1.GatewayTypeMicroGateway {
	// 	return createMicroGatewayDeployment(gateway, gatewayConfig), nil
	// } else {
	return createEnvoyGatewayDeployment(gateway, cfg)
	// }
}

// func createMicroGatewayDeployment(gateway *v1alpha1.Gateway, gatewayConfig config.Gateway) *appsv1.Deployment {
// 	podTemplateAnnotations := map[string]string{}
// 	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "true"
// 	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
// 	var cellName string
// 	cellName, ok := gateway.Labels[mesh.CellLabelKey]
// 	if !ok {
// 		cellName = gateway.Name
// 	}

// 	one := int32(1)
// 	return &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      GatewayDeploymentName(gateway),
// 			Namespace: gateway.Namespace,
// 			Labels:    createGatewayLabels(gateway),
// 			OwnerReferences: []metav1.OwnerReference{
// 				*controller.CreateGatewayOwnerRef(gateway),
// 			},
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: &one,
// 			Selector: createGatewaySelector(gateway),
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels:      createGatewayLabels(gateway),
// 					Annotations: podTemplateAnnotations,
// 				},
// 				Spec: corev1.PodSpec{
// 					InitContainers: []corev1.Container{
// 						{
// 							Name:  "cell-gateway-init",
// 							Image: gatewayConfig.InitImage,
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      configVolumeName,
// 									MountPath: configMountPath,
// 									ReadOnly:  true,
// 								},
// 								{
// 									Name:      setupConfigVolumeName,
// 									MountPath: setupConfigMountPath,
// 									ReadOnly:  true,
// 								},
// 								{
// 									Name:      gatewayBuildVolumeName,
// 									MountPath: gatewayBuildMountPath,
// 								},
// 							},
// 						},
// 					},
// 					Containers: []corev1.Container{
// 						{
// 							Name:  "cell-gateway",
// 							Image: gatewayConfig.Image,
// 							Env: []corev1.EnvVar{
// 								{
// 									Name:  "CELL_NAME",
// 									Value: cellName,
// 								},
// 							},
// 							Ports: []corev1.ContainerPort{{
// 								ContainerPort: gatewayContainerPort,
// 							}},
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      gatewayBuildVolumeName,
// 									MountPath: gatewayBuildMountPath,
// 								},
// 							},
// 						},
// 					},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: configVolumeName,
// 							VolumeSource: corev1.VolumeSource{
// 								ConfigMap: &corev1.ConfigMapVolumeSource{
// 									LocalObjectReference: corev1.LocalObjectReference{
// 										Name: ApiPublisherConfigMap(gateway),
// 									},
// 									Items: []corev1.KeyToPath{
// 										{
// 											Key:  apiConfigKey,
// 											Path: apiConfigFile,
// 										},
// 										{
// 											Key:  gatewayConfigKey,
// 											Path: gatewayConfigFile,
// 										},
// 									},
// 								},
// 							},
// 						},
// 						{
// 							Name: setupConfigVolumeName,
// 							VolumeSource: corev1.VolumeSource{
// 								ConfigMap: &corev1.ConfigMapVolumeSource{
// 									LocalObjectReference: corev1.LocalObjectReference{
// 										Name: ApiPublisherConfigMap(gateway),
// 									},
// 									Items: []corev1.KeyToPath{
// 										{
// 											Key:  gatewaySetupConfigKey,
// 											Path: gatewaySetupConfigFile,
// 										},
// 									},
// 								},
// 							},
// 						},
// 						{
// 							Name: gatewayBuildVolumeName,
// 							VolumeSource: corev1.VolumeSource{
// 								EmptyDir: &corev1.EmptyDirVolumeSource{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

func createEnvoyGatewayDeployment(gateway *v1alpha2.Gateway, cfg config.Interface) (*appsv1.Deployment, error) {
	// podTemplateAnnotations := map[string]string{}
	// podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "false"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	// var cellName string
	// cellName, ok := gateway.Labels[mesh.CellLabelKey]
	// if !ok {
	// 	cellName = gateway.Name
	// }

	var containers []corev1.Container
	var volumes []corev1.Volume

	containers = append(containers, *makeIstioProxyContainer(gateway, cfg))
	volumes = append(volumes, corev1.Volume{
		Name: "istio-certs",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "istio.default",
			},
		},
	})

	// add oidc filter container to the gateway if oidc is enabled
	if gateway.Spec.Ingress.IngressExtensions.HasOidc() {

		oidcContainer, err := makeOidcContainer(gateway, cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot create oidc container: %v", err)
		}
		containers = append(containers, *oidcContainer)
		volumes = append(volumes, corev1.Volume{
			Name: "oidc-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: gateway.Spec.Ingress.IngressExtensions.OidcConfig.SecretName,
				},
			},
		})
	}
	// 	containers = append(containers, corev1.Container{
	// 		Name:  "envoy-oidc-filter",
	// 		Image: gatewayConfig.OidcFilterImage,
	// 		Env: []corev1.EnvVar{
	// 			{
	// 				Name:  "PROVIDER_URL",
	// 				Value: oidc.ProviderUrl,
	// 			},
	// 			{
	// 				Name:  "CLIENT_ID",
	// 				Value: oidc.ClientId,
	// 			},
	// 			{
	// 				Name:  "CLIENT_SECRET",
	// 				Value: string(clientSecretBytes),
	// 			},
	// 			{
	// 				Name:  "DCR_ENDPOINT",
	// 				Value: oidc.DcrUrl,
	// 			},
	// 			{
	// 				Name:  "DCR_USER",
	// 				Value: oidc.DcrUser,
	// 			},
	// 			{
	// 				Name:  "DCR_PASSWORD",
	// 				Value: oidc.DcrPassword,
	// 			},
	// 			{
	// 				Name:  "REDIRECT_URL",
	// 				Value: oidc.RedirectUrl,
	// 			},
	// 			{
	// 				Name:  "APP_BASE_URL",
	// 				Value: oidc.BaseUrl,
	// 			},
	// 			{
	// 				Name:  "PRIVATE_KEY_FILE",
	// 				Value: "/etc/certs/key.pem",
	// 			},
	// 			{
	// 				Name:  "CERTIFICATE_FILE",
	// 				Value: "/etc/certs/cert.pem",
	// 			},
	// 			{
	// 				Name:  "JWT_ISSUER",
	// 				Value: cellName + "--gateway",
	// 			},
	// 			{
	// 				Name:  "JWT_AUDIENCE",
	// 				Value: cellName,
	// 			},
	// 			{
	// 				Name:  "SUBJECT_CLAIM",
	// 				Value: oidc.SubjectClaim,
	// 			},
	// 			{
	// 				Name:  "NON_SECURE_PATHS",
	// 				Value: strings.Join(oidc.NonSecurePaths, ","),
	// 			},
	// 			{
	// 				Name:  "SECURE_PATHS",
	// 				Value: strings.Join(oidc.SecurePaths, ","),
	// 			},
	// 			{
	// 				Name:  "SKIP_DISCOVERY_URL_CERT_VERIFY",
	// 				Value: gatewayConfig.SkipTlsVerify,
	// 			},
	// 		},
	// 		Ports: []corev1.ContainerPort{
	// 			{
	// 				ContainerPort: 15800,
	// 				Protocol:      corev1.ProtocolTCP,
	// 			},
	// 			{
	// 				ContainerPort: 15810,
	// 				Protocol:      corev1.ProtocolTCP,
	// 			},
	// 		},
	// 		VolumeMounts: []corev1.VolumeMount{
	// 			{
	// 				Name:      "cell-certs",
	// 				MountPath: "/etc/certs",
	// 				ReadOnly:  true,
	// 			},
	// 		},
	// 	})
	// }

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: gateway.MinReplicas(),
			Selector: makeSelector(gateway),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      makeLabels(gateway),
					Annotations: makePodAnnotations(gateway),
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    volumes,
				},
			},
		},
	}, nil
}

func makeIstioProxyContainer(gateway *v1alpha2.Gateway, cfg config.Interface) *corev1.Container {

	envVars := []corev1.EnvVar{
		{
			Name: "NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "spec.nodeName",
				},
			},
		},
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.name",
				},
			},
		},
		{
			Name: "POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		},
		{
			Name: "INSTANCE_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "status.podIP",
				},
			},
		},
		{
			Name: "HOST_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "status.hostIP",
				},
			},
		},
		{
			Name: "ISTIO_META_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.name",
				},
			},
		},
		{
			Name: "ISTIO_META_CONFIG_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		},
	}
	args := []string{
		"proxy",
		"router",
		"--domain",
		"$(POD_NAMESPACE).svc.cluster.local",
		"--drainDuration",
		"45s",
		"--parentShutdownDuration",
		"1m0s",
		"--connectTimeout",
		"10s",
		"--serviceCluster",
		gateway.Name + ".$(POD_NAMESPACE)",
		"--zipkinAddress",
		cfg.StringValue(config.ConfigMapKeyZipkinAddress),
		"--proxyAdminPort",
		"15000",
		"--statusPort",
		"15020",
		"--controlPlaneAuthPolicy",
		"NONE",
		"--discoveryAddress",
		"istio-pilot.istio-system:15010",
	}

	return &corev1.Container{
		Name:  "envoy-gateway",
		Image: fmt.Sprintf("docker.io/istio/proxyv2:%s", cfg.StringValue(config.ConfigMapKeyIstioVersion)),
		Env:   envVars,
		Args:  args,
		// Ports: []corev1.ContainerPort{
		// 	{
		// 		ContainerPort: 80,
		// 		Protocol:      corev1.ProtocolTCP,
		// 	},
		// 	{
		// 		ContainerPort: 443,
		// 		Protocol:      corev1.ProtocolTCP,
		// 	},
		// },
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "istio-certs",
				MountPath: "/etc/certs",
			},
		},
	}
}

func makeOidcContainer(gateway *v1alpha2.Gateway, cfg config.Interface) (*corev1.Container, error) {

	oidc := gateway.Spec.Ingress.IngressExtensions.OidcConfig

	key, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}

	clientSecretBytes, err := crypto.TryDecrypt(oidc.ClientSecret, key)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt the clientSecret: %v", err)
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "CELL_NAMESPACE",
			Value: gateway.Namespace,
		},
		{
			Name:  "PROVIDER_URL",
			Value: oidc.ProviderUrl,
		},
		{
			Name:  "CLIENT_ID",
			Value: oidc.ClientId,
		},
		{
			Name:  "CLIENT_SECRET",
			Value: string(clientSecretBytes),
		},
		{
			Name:  "DCR_ENDPOINT",
			Value: oidc.DcrUrl,
		},
		{
			Name:  "DCR_USER",
			Value: oidc.DcrUser,
		},
		{
			Name:  "DCR_PASSWORD",
			Value: oidc.DcrPassword,
		},
		{
			Name:  "REDIRECT_URL",
			Value: oidc.RedirectUrl,
		},
		{
			Name:  "APP_BASE_URL",
			Value: oidc.BaseUrl,
		},
		{
			Name:  "PRIVATE_KEY_FILE",
			Value: "/etc/certs/key.pem",
		},
		{
			Name:  "CERTIFICATE_FILE",
			Value: "/etc/certs/cert.pem",
		},
		{
			Name:  "JWT_ISSUER",
			Value: oidc.JwtIssuer,
		},
		{
			Name:  "JWT_AUDIENCE",
			Value: oidc.JwtAudience,
		},
		{
			Name:  "SUBJECT_CLAIM",
			Value: oidc.SubjectClaim,
		},
		{
			Name:  "NON_SECURE_PATHS",
			Value: strings.Join(oidc.NonSecurePaths, ","),
		},
		{
			Name:  "SECURE_PATHS",
			Value: strings.Join(oidc.SecurePaths, ","),
		},
		{
			Name:  "SKIP_DISCOVERY_URL_CERT_VERIFY",
			Value: cfg.StringValue(config.ConfigMapKeySkipTlsVerification),
		},
	}

	return &corev1.Container{
		Name:  "envoy-oidc-filter",
		Image: cfg.StringValue(config.ConfigMapKeyOidcImage),
		Env:   envVars,
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: 15800,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				ContainerPort: 15810,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "oidc-certs",
				MountPath: "/etc/certs",
				ReadOnly:  true,
			},
		},
	}, err
}

func RequireDeployment(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.HasRoutes()
}

func RequireDeploymentUpdate(gateway *v1alpha2.Gateway, deployment *appsv1.Deployment) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		deployment.Generation != gateway.Status.DeploymentGeneration
}

func CopyDeployment(source, destination *appsv1.Deployment) {
	destination.Spec.Template = source.Spec.Template
	destination.Spec.Selector = source.Spec.Selector
	destination.Spec.Replicas = source.Spec.Replicas
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromDeployment(gateway *v1alpha2.Gateway, deployment *appsv1.Deployment) {
	gateway.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	gateway.Status.DeploymentGeneration = deployment.Generation
	if deployment.Status.AvailableReplicas > 0 {
		gateway.Status.Status = v1alpha2.GatewayCurrentStatusReady
	} else {
		gateway.Status.Status = v1alpha2.GatewayCurrentStatusNotReady
	}
}
