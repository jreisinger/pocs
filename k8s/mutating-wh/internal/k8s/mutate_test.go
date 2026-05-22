package k8s

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var maxRequests = MaxRequests{
	Memory:       32 * 1024 * 1024,
	MemoryString: "32Mi",
}

func TestMutate(t *testing.T) {

	t.Run("when requested memory > 32Mi mutate it to 32Mi", func(t *testing.T) {
		pod64Mi, err := os.ReadFile("./testdata/nginx-pod-64Mi.json")
		require.NoError(t, err)

		b, err := Mutate(getAdmissionReview(t, string(pod64Mi)), maxRequests)
		require.NoError(t, err)

		var response admission.AdmissionReview
		require.NoError(t, json.Unmarshal(b, &response))

		expectedPatch := `[{"op":"replace","path":"/spec/containers/0/resources/requests/memory","value":"32Mi"}]`
		actualPatch := string(response.Response.Patch)
		assert.Equal(t, expectedPatch, actualPatch)
		assert.Equal(t, "mutating-wh", response.Response.AuditAnnotations["mutated"])
	})

	t.Run("when requested memory <= 32Mi don't mutate it", func(t *testing.T) {
		pod16Mi, err := os.ReadFile("./testdata/nginx-pod-16Mi.json")
		require.NoError(t, err)

		b, err := Mutate(getAdmissionReview(t, string(pod16Mi)), maxRequests)
		require.NoError(t, err)

		var response admission.AdmissionReview
		require.NoError(t, json.Unmarshal(b, &response))

		expectedPatch := ``
		actualPatch := string(response.Response.Patch)
		assert.Equal(t, expectedPatch, actualPatch)
		assert.Equal(t, 0, len(response.Response.AuditAnnotations))
	})
}

// --- test helpers ---

func getAdmissionReview(t *testing.T, object string) []byte {
	data := admission.AdmissionReview{
		Request: &admission.AdmissionRequest{
			UID: "abc123",
			Object: runtime.RawExtension{
				Raw: []byte(object),
			},
		},
	}

	out, err := json.Marshal(data)
	require.NoError(t, err)
	return out
}
