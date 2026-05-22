package k8s

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	admission "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type MaxRequests struct {
	Memory       int64  // ex. 32*1024*1024
	MemoryString string // ex. 32Mi
}

// Mutate takes AdminssionReview and changes pod resources requests to
// maxRequests if they're greater.
func Mutate(body []byte, maxRequests MaxRequests) ([]byte, error) {
	var admissionReview admission.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("unmarshal request as admission review: %w", err)
	}

	if admissionReview.Request == nil {
		return nil, errors.New("received nil admission review request")
	}

	var pod core.Pod
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("unmarshal request as pod: %v", err)
	}

	patches := []map[string]string{}
	for i, c := range pod.Spec.Containers {
		if c.Resources.Requests.Memory().Value() > maxRequests.Memory {
			patch := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/containers/%d/resources/requests/memory", i),
				"value": maxRequests.MemoryString,
			}
			patches = append(patches, patch)
		}
	}

	var patch []byte
	if len(patches) > 0 {
		var err error
		patch, err = json.Marshal(patches)
		if err != nil {
			return nil, fmt.Errorf("marshal patch: %v", err)
		}
	}

	admissionReview.Response = getAdmissionResponse(admissionReview.Request.UID, patch)
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		return nil, fmt.Errorf("marshal response body as admission review: %w", err)
	}
	return responseBody, nil
}

func getAdmissionResponse(admissionRequestUID types.UID, patch []byte) *admission.AdmissionResponse {
	patchTypeJson := admission.PatchTypeJSONPatch
	var admissionResponse = &admission.AdmissionResponse{
		Allowed: true,
		UID:     admissionRequestUID,
		Patch:   patch,
		Result: &meta.Status{
			Status: "Success",
		},
	}

	// patch type can be set only if there is actual patch
	if patch != nil {
		slog.Info(fmt.Sprintf("patch: %+v", string(patch)))
		admissionResponse.AuditAnnotations = map[string]string{"mutated": "mutating-wh"}
		admissionResponse.PatchType = &patchTypeJson
	}
	return admissionResponse
}
