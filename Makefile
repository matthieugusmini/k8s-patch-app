.PHONY: build
build:
	go build -o k8s-patch-app

.PHONY: test
test:
	go test -race -v ./...
