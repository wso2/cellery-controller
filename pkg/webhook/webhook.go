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

package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	"cellery.io/cellery-controller/pkg/apis"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
)

const (
	pathMutate   = "/mutate"
	pathValidate = "/validate"
)

type ServerOptions struct {
	Namespace             string
	ServerSecretName      string
	RootSecretName        string
	ServiceName           string
	DeploymentName        string
	MutatingWebhookName   string
	ValidatingWebhookName string
	Port                  int
}
type server struct {
	kubeClient kubernetes.Interface
	options    *ServerOptions
	logger     *zap.SugaredLogger
	validators map[schema.GroupVersionKind]func() apis.Validatable
}

func NewServer(kubeClient kubernetes.Interface, opt ServerOptions, logger *zap.SugaredLogger) *server {
	return &server{
		kubeClient: kubeClient,
		options:    &opt,
		logger:     logger.Named("admission-webhook"),
		validators: map[schema.GroupVersionKind]func() apis.Validatable{
			v1alpha2.SchemeGroupVersion.WithKind("Component"):    func() apis.Validatable { return &v1alpha2.Component{} },
			v1alpha2.SchemeGroupVersion.WithKind("Gateway"):      func() apis.Validatable { return &v1alpha2.Gateway{} },
			v1alpha2.SchemeGroupVersion.WithKind("TokenService"): func() apis.Validatable { return &v1alpha2.TokenService{} },
			v1alpha2.SchemeGroupVersion.WithKind("Cell"):         func() apis.Validatable { return &v1alpha2.Cell{} },
			v1alpha2.SchemeGroupVersion.WithKind("Composite"):    func() apis.Validatable { return &v1alpha2.Composite{} },
		},
	}
}

func (s *server) Run(stopCh <-chan struct{}) error {
	s.logger.Info("Configuring webhook tls certificates...")
	tlsCfg, caCertPEM, err := s.configureTls()
	if err != nil {
		s.logger.Errorf("Cannot configure webhook tls certificates: %v", err)
		return fmt.Errorf("tls configuration failed: %v", err)
	}
	s.logger.Info("Registering admission webhooks...")
	if err := s.registerWebhooks(caCertPEM); err != nil {
		s.logger.Errorf("Cannot register admission webhooks: %v", err)
		return err
	}
	addr := fmt.Sprintf(":%d", s.options.Port)
	srv := http.Server{
		Addr:      addr,
		Handler:   s,
		TLSConfig: tlsCfg,
	}

	go func() {
		s.logger.Infof("Serving admission webhook on %s", addr)
		if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			s.logger.Errorf("Admission webhook ListenAndServeTLS error: %v", err)
		}
	}()
	<-stopCh
	s.logger.Info("Shutting down admission webhook...")
	return srv.Shutdown(context.Background())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var admit func(*admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse
	if r.URL.Path == pathMutate {
		admit = s.mutate
	} else if r.URL.Path == pathValidate {
		admit = s.validate
	} else {
		http.NotFound(w, r)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, fmt.Sprintf("invalid Content-Type, got %q want %q", contentType, "application/json"), http.StatusUnsupportedMediaType)
		return
	}

	reviewRequest := admissionv1beta1.AdmissionReview{}
	if err := json.NewDecoder(r.Body).Decode(&reviewRequest); err != nil {
		http.Error(w, fmt.Sprintf("could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	reviewResponse := admissionv1beta1.AdmissionReview{}
	reviewResponse.Response = admit(reviewRequest.Request)
	if reviewRequest.Request != nil {
		reviewResponse.Response.UID = reviewRequest.Request.UID
	}

	if err := json.NewEncoder(w).Encode(reviewResponse); err != nil {
		http.Error(w, fmt.Sprintf("could encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *server) mutate(req *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	// todo: implement defaulting mutations
	return &admissionv1beta1.AdmissionResponse{Allowed: true}
}

func (s *server) validate(req *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	logger := s.makeLogger("validating", req)

	logger.Info("Validating object")

	gvk := SchemaGroupVersionKind(req.Kind)
	fn, ok := s.validators[gvk]
	if !ok {
		logger.Errorf("Unknown kind %v", gvk)
		return makeErrorResponse("unknown kind %v", gvk)
	}
	obj := fn()
	if err := json.Unmarshal(req.Object.Raw, &obj); err != nil {
		logger.Errorf("Cannot not unmarshal raw object: %v", err)
		return makeErrorResponse("cannot not unmarshal raw object: %v", err)
	}

	if err := obj.Validate(); err != nil {
		logger.Errorf("Validation failed: %v", err)
		return makeErrorResponse("validation failed: %v", err)
	}

	return &admissionv1beta1.AdmissionResponse{Allowed: true}
}

func (s *server) makeLogger(name string, req *admissionv1beta1.AdmissionRequest) *zap.SugaredLogger {
	return s.logger.Named(name).With(
		zap.String("uid", fmt.Sprint(req.UID)),
		zap.String("kind", fmt.Sprint(req.Kind)),
		zap.String("resource", fmt.Sprint(req.Resource)),
		zap.String("subresource", fmt.Sprint(req.SubResource)),
		zap.String("name", fmt.Sprint(req.Name)),
		zap.String("namespace", fmt.Sprint(req.Namespace)),
		zap.String("operation", fmt.Sprint(req.Operation)),
		zap.String("userinfo", fmt.Sprint(req.UserInfo)),
	)
}

func SchemaGroupVersionKind(gvk metav1.GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
}

func makeErrorResponse(reason string, args ...interface{}) *admissionv1beta1.AdmissionResponse {
	result := apierrors.NewBadRequest(fmt.Sprintf(reason, args...)).Status()
	return &admissionv1beta1.AdmissionResponse{
		Result:  &result,
		Allowed: false,
	}
}
