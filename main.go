package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var (
		kubeconfig *string
		name       *string
		patch      *string
	)
	name = flag.String("n", "", "Deployment name to patch")
	patch = flag.String("p", "", "JSON patch")
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if *name == "" {
		return errors.New("no deployment name provided")
	}
	if *patch == "" {
		return errors.New("no JSON patch to apply")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build configs from %s: %w", *kubeconfig, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create a new Clientset for the kube config: %w", err)
	}

	_, err = clientset.
		AppsV1().
		Deployments("default").
		Patch(ctx, *name, types.JSONPatchType, []byte(*patch), v1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch the deployment: %w", err)
	}

	return nil
}
