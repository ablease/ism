
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
CLI_NAME = bin/ism
GINKGO_ARGS = -r -p

all: clean test manager cli

# Run tests
test: generate fmt vet manifests unit-tests integration-tests acceptance-tests

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/pivotal-cf/ism/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

### CUSTOM MAKE RULES ###

cli:
	go build -o ${CLI_NAME} cmd/ism/main.go

clean:
	rm -f ${CLI_NAME}

acceptance-tests:
	ginkgo ${GINKGO_ARGS} acceptance

# skip integration/acceptance tests
unit-tests:
	ginkgo ${GINKGO_ARGS} -skipPackage acceptance,kube,pkg/controller,pkg/api,pkg/internal/repositories

integration-tests: cli-integration-tests kube-integration-tests

cli-integration-tests:
	ginkgo ${GINKGO_ARGS} kube

kube-integration-tests:
	ginkgo ${GINKGO_ARGS} pkg/controller pkg/api pkg/internal/repositories
