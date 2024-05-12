package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"k8s-patch-app/k8s"
	"os"
	"os/signal"
	"path/filepath"

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
	// Has to be done here as the defer statement wouldn't get executed in the
	// main function.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	var (
		kubeconfigPath *string
		name           *string
		patch          *string
	)
	// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go#L44
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = flag.String("k", filepath.Join(home, ".kube", "config"), "(optional) Absolute path to the kubeconfig file")
	} else {
		kubeconfigPath = flag.String("k", "", "Absolute path to the kubeconfig file")
	}
	name = flag.String("n", "", "Deployment name to patch")
	patch = flag.String("p", "", "JSON patch")

	flag.Parse()

	if *name == "" {
		return errors.New("no deployment name provided")
	}
	if *patch == "" {
		return errors.New("no JSON patch to apply")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to build configs from %s: %w", *kubeconfigPath, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create a new Clientset for the kube config: %w", err)
	}

	jsonPatcher := k8s.NewJSONPatcher(clientset)
	return jsonPatcher.PatchDeployment(ctx, "default", *name, []byte(*patch))
}
