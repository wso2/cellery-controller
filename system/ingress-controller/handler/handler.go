package handler

import (
	log "github.com/golang/glog"
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"github.com/wso2/product-vick/system/ingress-controller/config"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj *v1alpha1.Ingress)
	ObjectDeleted(obj *v1alpha1.Ingress)
	ObjectUpdated(objOld, objNew *v1alpha1.Ingress)
}

var ingressRuleCreator *IngressRuleCreator

// IngressHandler is a sample implementation of Handler
type IngressHandler struct{}

// Init handles any handler initialization
func (t *IngressHandler) Init() error {
	ingressConfig, err := config.GetIngressConfigs("/etc/gw.json")
	if err != nil {
		return err
	}
	ingressRuleCreator = NewIngressRuleCreator(ingressConfig)
	log.Infof("IngressHandler Initialized")
	return nil
}

// ObjectCreated is called when an object is created
func (t *IngressHandler) ObjectCreated(obj *v1alpha1.Ingress) {
	clientId, clientSecret, err := ingressRuleCreator.RegisterClient()
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	err = ingressRuleCreator.CreateApi(token, &obj.Spec)
	if err != nil {
		// log and continue
		log.Error(err)
	}
	log.Infof("Ingress Created %+v\n", obj)
}

// ObjectDeleted is called when an object is deleted
func (t *IngressHandler) ObjectDeleted(obj *v1alpha1.Ingress) {
	clientId, clientSecret, err := ingressRuleCreator.RegisterClient()
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	err = ingressRuleCreator.DeleteApi(token, &obj.Spec)
	if err != nil {
		// log and continue
		log.Error(err)
	}
	log.Infof("Ingress Deleted %+v\n", obj)
}

// ObjectUpdated is called when an object is updated
func (t *IngressHandler) ObjectUpdated(objOld, objNew *v1alpha1.Ingress) {
	// the old and new objects can't differ in name, context and version.
	// Those are not allowed to change in an update.
	// if they need to be changed, that is a new API version
	if (objOld.Spec.Name != objNew.Spec.Name) || (objOld.Spec.Context != objNew.Spec.Context) ||
		(objOld.Spec.Version != objNew.Spec.Version) {
		log.Error("Modifying either name, context or version in an existing Ingress rule is not permitted")
		return
	}

	log.Infof("Ingress Updated %+v\n", objNew)
	//newIngRules := getNewIngressRules(*objOld, *objNew)
	//updatedIngRules := getUpdatedIngressRules(*objOld, *objNew)
	//deletedIngRules := getDeletedIngressRules(*objOld, *objNew)

	clientId, clientSecret, err := ingressRuleCreator.RegisterClient()
	if err != nil {
		log.Error(err)
		return
	}
	token, err := ingressRuleCreator.GenerateAccessToken(clientId, clientSecret)
	if err != nil {
		log.Error(err)
		return
	}
	//// create new ones
	//if newIngRules != nil {
	//	for _, ingRule := range newIngRules {
	//		err := ingressRuleCreator.CreateApi(objNew.Spec.Name, objNew.Spec.Context, objNew.Spec.Version, token, &ingRule)
	//		if err != nil {
	//			// log and continue
	//			log.Error(err)
	//		}
	//	}
	//}

	// update existing ones which changed
	err = ingressRuleCreator.UpdateApi(token, &objNew.Spec)
	if err != nil {
		// log and continue
		log.Error(err)
	}

	//// delete non-existing ones
	//if deletedIngRules != nil {
	//	for _, ingRule := range deletedIngRules {
	//		err := ingressRuleCreator.DeleteApi(objNew.Spec.Name, objNew.Spec.Context, objNew.Spec.Version, token, &ingRule)
	//		if err != nil {
	//			// log and continue
	//			log.Error(err)
	//		}
	//	}
	//}
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
