package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// JSONPatcher is a wrapper around a kubernetes.Interface that apply
// JSON patch to k8s resources.
type JSONPatcher struct {
	client kubernetes.Interface
}

// NewJSONPatcher returns a newly instanciated JSONPatcher.
func NewJSONPatcher(client kubernetes.Interface) *JSONPatcher {
	return &JSONPatcher{
		client: client,
	}
}

// PatchDeploymment apply the given patch to the deployment referenced by name in namespace.
func (p *JSONPatcher) PatchDeployment(ctx context.Context, namespace string, name string, patch []byte) error {
	_, err := p.client.
		AppsV1().
		Deployments(namespace).
		Patch(ctx, name, types.JSONPatchType, patch, v1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch the deployment: %w", err)
	}
	return nil
}
