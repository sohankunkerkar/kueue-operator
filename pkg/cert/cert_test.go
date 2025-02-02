package cert

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInjectCertAnnotation(t *testing.T) {
	testcases := map[string]struct {
		annotations         map[string]string
		expectedAnnotations map[string]string
	}{
		"nil annotations": {
			annotations:         nil,
			expectedAnnotations: map[string]string{"cert-manager.io/inject-ca-from": fmt.Sprintf("%s/webhook-cert", "test")},
		},
		"existingAnnotations": {
			annotations: map[string]string{"hello": "world"},
			expectedAnnotations: map[string]string{
				"hello":                          "world",
				"cert-manager.io/inject-ca-from": fmt.Sprintf("%s/webhook-cert", "test"),
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			got := InjectCertAnnotation(tc.annotations, "test")
			if diff := cmp.Diff(got, tc.expectedAnnotations); len(diff) != 0 {
				t.Errorf("Unexpected buckets (-want,+got):\n%s", diff)
			}
		})

	}
}
