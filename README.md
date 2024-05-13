# k8s-patch-app

This Go application provides functionality to patch Kubernetes deployments using JSON patches. It utilizes the Kubernetes client library.

## 🚀 Prerequisites

Ensure you have the following prerequisites installed:

- Go
- Kubernetes cluster access
- kubectl configured (for obtaining kubeconfig)

## 💾 Installation

Clone the repository:

```bash
git clone git@github.com:matthieugusmini/k8s-patch-app.git
cd k8s-patch-app
```

Build the application:

```bash
make build
```

## ✍ Usage

Run the application with the following command:

```bash
./k8s-patch-app -n <deployment-name> -p <json-patch>
```

For more details:
```bash
$ ./k8s-patch-app --help
Usage of ./k8s-patch-app:
  -k string
        (optional) absolute path to the kubeconfig file
  -n string
        Deployment name to patch
  -p string
        JSON patch
```

Example:

```bash
./k8s-patch-app -n=foo -p='[{"op": "replace", "path": "/spec/replicas", "value": 42}]'
```

## 📝 Description

The application reads the Kubernetes configuration from the provided kubeconfig file or uses the default location (`~/.kube/config`). It then patches the specified deployment with the provided JSON patch.
