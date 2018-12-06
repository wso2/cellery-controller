package handler

import (
	"crypto/tls"
	log "github.com/golang/glog"
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"github.com/wso2/product-vick/system/ingress-controller/config"
	"net/http"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj *v1alpha1.Ingress)
	ObjectDeleted(obj *v1alpha1.Ingress)
	ObjectUpdated(objOld, objNew *v1alpha1.Ingress)
}

var ingressRuleCreator *IngressRuleCreator
var httpClient *http.Client

// IngressHandler is a sample implementation of Handler
type IngressHandler struct{}

// Init handles any handler initialization
func (t *IngressHandler) Init() error {
	ingressConfig, err := config.GetIngressConfigs("/etc/gw.json")
	if err != nil {
		return err
	}
	ingressRuleCreator = NewIngressRuleCreator(ingressConfig)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: transport}
	log.Infof("IngressHandler Initialized")
	return nil
}

// ObjectCreated is called when an object is created
func (t *IngressHandler) ObjectCreated(obj *v1alpha1.Ingress) {
	regCaller, err := RegisterClientHttpCaller(ingressRuleCreator, httpClient)
	if err != nil {
		log.Error(err)
		return
	}
	clientId, clientSecret, err := ingressRuleCreator.RegisterClient(regCaller)
	if err != nil {
		log.Error(err)
		return
	}
	tokCaller, err := GenerateTokenHttpCaller(ingressRuleCreator, httpClient, clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(tokCaller)
	if err != nil {
		log.Error(err)
		return
	}
	createApiCaller, err := CreateApiHttpCaller(ingressRuleCreator, httpClient, &obj.Spec, token)
	if err != nil {
		log.Error(err)
		return
	}
	apiId, err := ingressRuleCreator.CreateApi(createApiCaller, &obj.Spec)
	if err != nil {
		// log and continu
		log.Error(err)
		return
	}
	log.Infof("Ingress successfully created %+v\n", obj)
	publishApiCaller, err := PublishApiHttpCaller(ingressRuleCreator, httpClient, apiId, token)
	if err != nil {
		log.Error(err)
		return
	}
	err = ingressRuleCreator.PublishApi(publishApiCaller, apiId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Ingress successfully published %+v\n", obj)
}

// ObjectDeleted is called when an object is deleted
func (t *IngressHandler) ObjectDeleted(obj *v1alpha1.Ingress) {
	regCaller, err := RegisterClientHttpCaller(ingressRuleCreator, httpClient)
	if err != nil {
		log.Error(err)
		return
	}
	clientId, clientSecret, err := ingressRuleCreator.RegisterClient(regCaller)
	if err != nil {
		log.Error(err)
		return
	}
	tokCaller, err := GenerateTokenHttpCaller(ingressRuleCreator, httpClient, clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(tokCaller)
	if err != nil {
		log.Error(err)
		return
	}
	getApiIdCaller, err := GetApiIdHttpCaller(ingressRuleCreator, obj.Spec.Context, httpClient, token)
	if err != nil {
		log.Error(err)
		return
	}
	apiId, err := ingressRuleCreator.GetApiId(getApiIdCaller, obj.Spec.Context, obj.Spec.Version)
	if err != nil {
		log.Error(err)
		return
	}
	deleteApiCaller, err :=  DeleteApiHttpCaller(ingressRuleCreator, httpClient, apiId, token)
	if err != nil {
		log.Error(err)
		return
	}
	err = ingressRuleCreator.DeleteApi(deleteApiCaller, apiId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Ingress successfully deleted %+v\n", obj)
}

// ObjectUpdated is called when an object is updated
func (t *IngressHandler) ObjectUpdated(objOld, objNew *v1alpha1.Ingress) {
	// the old and new objects can't differ in name, context and version.
	// Those are not allowed to change in an update.
	// if they need to be changed, that is a new API version
	// TODO: handle this validation in a K8s admission validation controller
	if (objOld.Spec.Name != objNew.Spec.Name) || (objOld.Spec.Context != objNew.Spec.Context) ||
		(objOld.Spec.Version != objNew.Spec.Version) {
		log.Error("Modifying either name, context or version in an existing Ingress rule is not permitted")
		return
	}

	regCaller, err := RegisterClientHttpCaller(ingressRuleCreator, httpClient)
	if err != nil {
		log.Error(err)
		return
	}
	clientId, clientSecret, err := ingressRuleCreator.RegisterClient(regCaller)
	if err != nil {
		log.Error(err)
		return
	}
	tokCaller, err := GenerateTokenHttpCaller(ingressRuleCreator, httpClient, clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(tokCaller)
	if err != nil {
		log.Error(err)
		return
	}
	getApiIdCaller, err := GetApiIdHttpCaller(ingressRuleCreator, objNew.Spec.Context, httpClient, token)
	if err != nil {
		log.Error(err)
		return
	}
	apiId, err := ingressRuleCreator.GetApiId(getApiIdCaller, objNew.Spec.Context, objNew.Spec.Version)
	if err != nil {
		log.Error(err)
		return
	}
	updateCaller, err := UpdateApiHttpCaller(ingressRuleCreator, &objNew.Spec, httpClient, apiId, token)
	if err != nil {
		log.Error(err)
		return
	}
	apiId, err = ingressRuleCreator.UpdateApi(updateCaller, &objNew.Spec)
	if err != nil {
		log.Error(err)
		return
	}
	publishApiCaller, err := PublishApiHttpCaller(ingressRuleCreator, httpClient, apiId, token)
	if err != nil {
		log.Error(err)
		return
	}
	err = ingressRuleCreator.PublishApi(publishApiCaller, apiId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Ingress successfully updated %+v\n", objNew)
}

//func getNewIngressRules (objOld, objNew v1.Ingress) ([]v1.IngressRule) {
//	var newIngRules []v1.IngressRule
//	var matchfound bool
//	for _, newRule := range objNew.Spec.Rules {
//		matchfound = false
//		for _, oldRule := range objOld.Spec.Rules {
//			if newRule.IsEqual(oldRule) {
//				// match found, not a new rule
//				matchfound = true
//				break
//			}
//		}
//		if !matchfound {
//			// new rule
//			newIngRules = append(newIngRules, newRule)
//		}
//	}
//	return newIngRules
//}
//
//func getUpdatedIngressRules (objOld, objNew v1.Ingress) ([]v1.IngressRule) {
//	var updatedIngRules []v1.IngressRule
//	for _, newRule := range objNew.Spec.Rules {
//		for _, oldRule := range objOld.Spec.Rules {
//			if newRule.IsEqual(oldRule) {
//				// match found, modified rule
//				updatedIngRules = append(updatedIngRules, newRule)
//				break
//			}
//		}
//	}
//	return updatedIngRules
//}
//
//func getDeletedIngressRules (objOld, objNew v1.Ingress) ([]v1.IngressRule) {
//	var deletedIngRules []v1.IngressRule
//	var matchfound bool
//	for _, oldRule := range objOld.Spec.Rules {
//		matchfound = false
//		for _, newRule := range objNew.Spec.Rules {
//			if oldRule.IsEqual(newRule) {
//				// match found, not a new rule
//				matchfound = true
//				break
//			}
//		}
//		if !matchfound {
//			// new rule
//			deletedIngRules = append(deletedIngRules, oldRule)
//		}
//	}
//	return deletedIngRules
//}
