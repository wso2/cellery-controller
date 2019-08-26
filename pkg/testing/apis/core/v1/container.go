/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

type ContainerOption func(*corev1.Container)

func Container(opt ...ContainerOption) *corev1.Container {
	c := &corev1.Container{}
	for _, opt := range opt {
		opt(c)
	}

	return c
}

func WithContainerName(name string) ContainerOption {
	return func(c *corev1.Container) {
		c.Name = name
	}
}

func WithContainerImage(image string) ContainerOption {
	return func(c *corev1.Container) {
		c.Image = image
	}
}

func WithContainerCommand(command []string) ContainerOption {
	return func(c *corev1.Container) {
		c.Command = command
	}
}

func WithContainerArgs(args []string) ContainerOption {
	return func(c *corev1.Container) {
		c.Args = args
	}
}

func WithContainerPort(port int32) ContainerOption {
	return func(c *corev1.Container) {
		c.Ports = append(c.Ports, corev1.ContainerPort{
			ContainerPort: port,
		})
	}
}

func WithContainerEnv(env corev1.EnvVar) ContainerOption {
	return func(c *corev1.Container) {
		c.Env = append(c.Env, env)
	}
}

func WithContainerEnvFromValue(name string, value string) ContainerOption {
	return func(c *corev1.Container) {
		c.Env = append(c.Env, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
	}
}

func WithContainerEnvFrom(envSource corev1.EnvFromSource) ContainerOption {
	return func(c *corev1.Container) {
		c.EnvFrom = append(c.EnvFrom, envSource)
	}
}

func WithContainerVolumeMounts(name string, readonly bool, mountPath string, subPath string) ContainerOption {
	return func(c *corev1.Container) {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      name,
			ReadOnly:  readonly,
			MountPath: mountPath,
			SubPath:   subPath,
		})
	}
}
