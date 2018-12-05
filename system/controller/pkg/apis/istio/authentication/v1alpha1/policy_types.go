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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Taken from: https://github.com/istio/istio/blob/1.0.2/vendor/istio.io/api/networking/v1alpha3/envoy_filter.pb.go#L180

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PolicySpec `json:"spec"`
}

type PolicySpec struct {
	// List rules to select destinations that the policy should be applied on.
	// If empty, policy will be used on all destinations in the same namespace.
	Targets []*TargetSelector `json:"targets,omitempty"`
	// List of authentication methods that can be used for peer authentication.
	// They will be evaluated in order; the first validate one will be used to
	// set peer identity (source.user) and other peer attributes. If none of
	// these methods pass, and peer_is_optional flag is false (see below),
	// request will be rejected with authentication failed error (401).
	// Leave the list empty if peer authentication is not required
	Peers []*PeerAuthenticationMethod `json:"peers,omitempty"`
	// Set this flag to true to accept request (for peer authentication perspective),
	// even when none of the peer authentication methods defined above satisfied.
	// Typically, this is used to delay the rejection decision to next layer (e.g
	// authorization).
	// This flag is ignored if no authentication defined for peer (peers field is empty).
	PeerIsOptional bool `json:"peerIsOptional,omitempty"`
	// List of authentication methods that can be used for origin authentication.
	// Similar to peers, these will be evaluated in order; the first validate one
	// will be used to set origin identity and attributes (i.e request.auth.user,
	// request.auth.issuer etc). If none of these methods pass, and origin_is_optional
	// is false (see below), request will be rejected with authentication failed
	// error (401).
	// Leave the list empty if origin authentication is not required.
	Origins []*OriginAuthenticationMethod `json:"origins,omitempty"`
	// Set this flag to true to accept request (for origin authentication perspective),
	// even when none of the origin authentication methods defined above satisfied.
	// Typically, this is used to delay the rejection decision to next layer (e.g
	// authorization).
	// This flag is ignored if no authentication defined for origin (origins field is empty).
	OriginIsOptional bool `json:"originIsOptional,omitempty"`
	// Define whether peer or origin identity should be use for principal. Default
	// value is USE_PEER.
	// If peer (or orgin) identity is not available, either because of peer/origin
	// authentication is not defined, or failed, principal will be left unset.
	// In other words, binding rule does not affect the decision to accept or
	// reject request.
	PrincipalBinding string `json:"principalBinding,omitempty"`
}


type TargetSelector struct {
	// REQUIRED. The name must be a short name from the service registry. The
	// fully qualified domain name will be resolved in a platform specific manner.
	Name string `json:"name,omitempty"`
	// Specifies the ports on the destination. Leave empty to match all ports
	// that are exposed.
	Ports []*PortSelector `json:"ports,omitempty"`
}

type PortSelector struct {
	Number int32 `json:"number,omitempty"`
}

type PeerAuthenticationMethod struct {
	Mtls string `json:"mtls"`
}

// OriginAuthenticationMethod defines authentication method/params for origin
// authentication. Origin could be end-user, device, delegate service etc.
// Currently, only JWT is supported for origin authentication.
type OriginAuthenticationMethod struct {
	// Jwt params for the method.
	Jwt *Jwt `json:"jwt,omitempty"`
}

// JSON Web Token (JWT) token format for authentication as defined by
// https://tools.ietf.org/html/rfc7519. See [OAuth
// 2.0](https://tools.ietf.org/html/rfc6749) and [OIDC
// 1.0](http://openid.net/connect) for how this is used in the whole
// authentication flow.
//
// Example,
//
// ```yaml
// issuer: https://example.com
// audiences:
// - bookstore_android.apps.googleusercontent.com
//   bookstore_web.apps.googleusercontent.com
// jwksUri: https://example.com/.well-known/jwks.json
// ```
type Jwt struct {
	// Identifies the issuer that issued the JWT. See
	// [issuer](https://tools.ietf.org/html/rfc7519#section-4.1.1)
	// Usually a URL or an email address.
	//
	// Example: https://securetoken.google.com
	// Example: 1234567-compute@developer.gserviceaccount.com
	Issuer string `json:"issuer,omitempty"`
	// The list of JWT
	// [audiences](https://tools.ietf.org/html/rfc7519#section-4.1.3).
	// that are allowed to access. A JWT containing any of these
	// audiences will be accepted.
	//
	// The service name will be accepted if audiences is empty.
	//
	// Example:
	//
	// ```yaml
	// audiences:
	// - bookstore_android.apps.googleusercontent.com
	//   bookstore_web.apps.googleusercontent.com
	// ```
	Audiences []string `json:"audiences,omitempty"`
	// URL of the provider's public key set to validate signature of the
	// JWT. See [OpenID
	// Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata).
	//
	// Optional if the key set document can either (a) be retrieved from
	// [OpenID
	// Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html) of
	// the issuer or (b) inferred from the email domain of the issuer (e.g. a
	// Google service account).
	//
	// Example: https://www.googleapis.com/oauth2/v1/certs
	JwksUri string `json:"jwksUri,omitempty"`
	// JWT is sent in a request header. `header` represents the
	// header name.
	//
	// For example, if `header=x-goog-iap-jwt-assertion`, the header
	// format will be x-goog-iap-jwt-assertion: <JWT>.
	JwtHeaders []string `json:"jwtHeaders,omitempty"`
	// JWT is sent in a query parameter. `query` represents the
	// query parameter name.
	//
	// For example, `query=jwt_token`.
	JwtParams []string `json:"jwtParams,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Policy `json:"items"`
}
