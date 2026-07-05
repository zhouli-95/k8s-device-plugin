APP_NAME := device-plugin
OUT_DIR := out
IMAGE_TAG := "dev"

VERSION ?= "0.0.1"
DATE=`date -Iseconds`
COMMIT?=`git rev-parse --verify HEAD`
LDFLAGS="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

GOPROXY ?= "https://goproxy.cn,direct"

all: build

build:
	@echo "Building binary"
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR) -ldflags $(LDFLAGS) ./... 

docker:
	@echo "Building docker image"
	@docker build . \
	-f build/package/Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg GIT_COMMIT=$(COMMIT) \
	--build-arg GOPROXY=$(GOPROXY) \
	-t $(APP_NAME):$(IMAGE_TAG)

clean:
	@echo "Cleaning up..."
	@go clean
	rm -rf $(OUT_DIR)
	@echo "Cleanup complete!"

deploy:
	@echo "Deploying to Kubernetes"
	@kubectl apply -k deployments/kustomize/base

undeploy:
	@echo "Undeploying from Kubernetes"
	@kubectl delete -k deployments/kustomize/base

generate-static:
	@echo "Generating static yaml"
	@kubectl kustomize deployments/kustomize/base > deployments/static/device-plugin.yaml

test:
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@go fmt  ./...

cover:
	@echo "go cover"
	@go test -v -coverprofile=$(OUT_DIR)/coverage.out ./...
	@go tool cover -html=$(OUT_DIR)/coverage.out -o $(OUT_DIR)/coverage.html

prepare:
	@echo "prepare"
	@go install github.com/onsi/ginkgo/v2/ginkgo@v2.27.5
	@go install sigs.k8s.io/kind@v0.32.0

e2e:
	@echo "e2e test"
	@BUILD_IMAGE=true CREATE_KIND_CLUSTER=true bash ./scripts/e2e-test.sh

.PHONY: all build clean test fmt cover e2e deploy undeploy
