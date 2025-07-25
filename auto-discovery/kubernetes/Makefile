# SPDX-FileCopyrightText: the secureCodeBox authors
#
# SPDX-License-Identifier: Apache-2.0

include ../../prerequisites.mk

IMG_NS ?= securecodebox

# Image URL to use all building/pushing image targets
IMG ?= auto-discovery-kubernetes

# Tag used for the image
IMG_TAG ?= sha-$$(git rev-parse --short HEAD)

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.33.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set).
# IMPORTANT: The body of conditionals MUST not be indented! Indentation result in
# errors on macOS/FreeBSD because the line wil be interpreted as command which must
# inside a recipe (target). (see https://github.com/secureCodeBox/secureCodeBox/issues/1353)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=auto-discovery,headerFile="hack/boilerplate.yaml.txt" crd paths="./..." output:crd:artifacts:config=crds  output:rbac:artifacts:config=templates/rbac

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin

.PHONY: test
test: manifests fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out
	
.PHONY: ginko
ginkgo: manifests fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" ginko test ./... -coverprofile cover.out -v --progress

## Helm unit tests
.PHONY: helm-unit-tests
helm-unit-tests:
	echo "Running helm unit tests for operator"; \
	helm unittest . \

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: manifests fmt vet ## Run a controller from your host.
	go run ./main.go --config ./auto-discovery-config.yaml

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	@echo ".: ⚙️ Build Container Images"
	docker build -t ${IMG_NS}/${IMG}:${IMG_TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG_NS}/${IMG}:${IMG_TAG}

.PHONY: docker-export
docker-export:
	@echo ".: 💾 Export Container Images"
	docker save $(IMG_NS)/$(IMG):$(IMG_TAG) > $(IMG).tar

##@ Deployment

.PHONY: kind-import
kind-import:
	@echo ".: 💾 Importing the image archive to local kind cluster."
	kind load image-archive ./$(IMG).tar

.PHONY: helm-deploy
helm-deploy:
	@echo ".: ⚙️ Deploying auto-discovery with the Image tag '$(IMG_TAG)' into kind."
	helm upgrade --install auto-discovery ./ --wait \
		--set="image.repository=docker.io/$(IMG_NS)/$(IMG)" \
		--set="image.tag=$(IMG_TAG)" \
		--set="image.pullPolicy=IfNotPresent" \

# IMPORTANT: The body of conditionals MUST not be indented! Indentation result in
# errors on macOS/FreeBSD because the line wil be interpreted as command which must
# inside a recipe (target). (see https://github.com/secureCodeBox/secureCodeBox/issues/1353)
ifndef ignore-not-found
ignore-not-found = false
endif

.PHONY: install
install: manifests ## Install CRDs into the K8s cluster specified in ~/.kube/config.

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen-$(CONTROLLER_TOOLS_VERSION)
ENVTEST ?= $(LOCALBIN)/setup-envtest-$(ENVTEST_VERSION)

## Tool Versions
CONTROLLER_TOOLS_VERSION ?= v0.18.0
ENVTEST_VERSION ?= release-0.21

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
