package cert

import "fmt"

func InjectCertAnnotation(annotation map[string]string, namespace string) map[string]string {
	newAnnotation := annotation
	if annotation == nil {
		newAnnotation = map[string]string{}
	}
	newAnnotation["cert-manager.io/inject-ca-from"] = fmt.Sprintf("%s/webhook-cert", namespace)
	return newAnnotation
}
