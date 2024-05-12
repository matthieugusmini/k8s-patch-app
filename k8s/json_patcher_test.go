package k8s_test

import (
	"context"
	"k8s-patch-app/k8s"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJSONPatcher_PatchDeployment(t *testing.T) {
	testCases := []struct {
		name               string
		deployment         runtime.Object
		patch              []byte
		expectedReplicasNb int
	}{
		{
			name: "scale deployment",
			deployment: &appsv1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-deploy",
					Namespace: "my-ns",
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: ptr(int32(1)),
				},
			},
			patch:              []byte(`[{"op": "replace", "path": "/spec/replicas", "value": 2}]`),
			expectedReplicasNb: 2,
		},
		{
			name: "downscale deployment",
			deployment: &appsv1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-deploy",
					Namespace: "my-ns",
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: ptr(int32(2)),
				},
			},
			patch:              []byte(`[{"op": "replace", "path": "/spec/replicas", "value": 1}]`),
			expectedReplicasNb: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			clientset := fake.NewSimpleClientset(tc.deployment)

			jsonPatcher := k8s.NewJSONPatcher(clientset)
			err := jsonPatcher.PatchDeployment(ctx, "my-ns", "my-deploy", tc.patch)
			if err != nil {
				t.Fatalf("failed to patch deployment: %v", err)
			}

			deploy, err := clientset.
				AppsV1().
				Deployments("my-ns").
				Get(ctx, "my-deploy", v1.GetOptions{})
			if err != nil {
				t.Fatalf("failed to get deployment: %v", err)
			}
			if *deploy.Spec.Replicas != int32(tc.expectedReplicasNb) {
				t.Fatalf(
					"deployment replicas should be scaled to %d but is equal to %d",
					tc.expectedReplicasNb,
					*deploy.Spec.Replicas,
				)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
