# =========================================================================
# Variable Definitions
# =========================================================================

HOST_PWD ?= $(PWD)
SERVICE_NAME = disk-manager-exercise
VERSION=$(shell git rev-parse HEAD || echo none)
RPM_BUILD_ROOT ?= $(HOST_PWD)/build/rpmbuild
GO_SOURCE_FILES = $(shell find -name "*.go")
DOCKER_VERSION = $(shell docker version -f '{{.Server.Version}}' | cut -d"." -f1)
MIN_DOCKER_VERSION = 18
STRATO_GO_SWAGGER = docker run --rm \
	-e GOPATH=/go \
	-v $(HOST_PWD):/go/src/github.com/Stratoscale/disk-manager-exercise \
	-w /go/src/github.com/Stratoscale/disk-manager-exercise \
	-u $(shell id -u):$(shell id -g) \
	stratoscale/swagger:v1.0.13

# =========================================================================

.PHONY: build generate-by-swagger generate-client generate-server go-generate clean image format lint test subsystem

all: lint test subsystem rpm

build: build/disk-manager-exercise build/disk-manager-exercise-client/dist/disk-manager-exercise-client-*.tar.gz

build/disk-manager-exercise: $(GO_SOURCE_FILES)
	# Build service binary
	CGO_ENABLED=0 GOOS=linux go build -o $@ ./main.go

image: build
	# Build service image
	skipper build $(name)

lint:
	# Run static code analysis
	#golangci-lint run --enable goimports --enable gocyclo --tests
	golangci-lint run --enable goimports --enable gocyclo --enable interfacer --enable unconvert --tests

test: flake8 pylint pytest
	# Run go unit tests
	go test -race ./...

flake8:
	python -m flake8 --config=setup.cfg lib/pytools --exclude=test_*

pylint:
	PYLINTHOME=$(PYLINTHOME) pylint -r n lib/pytools --disable=missing-docstring --max-line-length=145

pytest:
	nose2 --start-dir lib/

format:
	# Auto format go code
	goimports -w $(shell find -maxdepth 1 -name "*.go" -or -type d -not -name vendor -not -name .)

check_docker_version:
	# Ensure docker version is recent enough
	if [ $(DOCKER_VERSION) -lt $(MIN_DOCKER_VERSION) ] ; then \
		echo "Docker version must be at least $(MIN_DOCKER_VERSION)"; \
		exit 1; \
	fi

generate-by-swagger: check_docker_version generate-server generate-client go-generate

generate-server:
	# Cleanup old generated code
	-find models restapi -nowarn -name "*.go" -not -name "_mock.go" -delete

	# Generate server code based on swagger file
	$(STRATO_GO_SWAGGER) generate server

generate-client:
	# Cleanup old generated code
	-find diskmanagerexerciseclient -nowarn -name "*.go" -not -name "_mock.go" -delete
	# Generate server code based on swagger file
	$(STRATO_GO_SWAGGER) generate client --client-package diskmanagerexerciseclient

go-generate:
	# Generate code based on annotations (e.g mocks)
	go generate ./...

clean:
	# Cleanup artifacts
	rm -rf build

dep-ensure:
	# Ensure a dependency is safely vendored in the project
	dep ensure

rpm: $(shell find deploy -type f)
	rpmbuild -bb -vv --define "_srcdir $(HOST_PWD)" --define "_topdir $(RPM_BUILD_ROOT)" deploy/install.spec

build/disk-manager-exercise-client/dist/disk-manager-exercise-client-*.tar.gz: swagger.yaml
	# Generate python client package based on the swagger file
	mkdir -p build
	# Use swagger-codegen container to generate python client code
	@echo '{"packageName" : "disk_manager_exercise_client", "packageVersion": "$(VERSION)"}' > build/code-gen-config.json
	docker run --rm \
        -u $(shell id -u $(USER)) \
        -v $(HOST_PWD)/build:/swagger-api/out \
        -v $(HOST_PWD)/swagger.yaml:/swagger.yaml:ro \
        -v $(HOST_PWD)/build/code-gen-config.json:/config.json:ro \
        jimschubert/swagger-codegen-cli generate \
        --lang python \
        --config /config.json \
        --output ./disk-manager-exercise-client/ \
        --input-spec /swagger.yaml
	# Create the client source distribution
	cd build/disk-manager-exercise-client/ && python setup.py sdist

subsystem: image
	# Cleanup old subsystem logs
	rm -rf subsystem/logs && mkdir -p subsystem/logs

	# Run service subsystem tests
	DTT_COLLECT_STATS=$(BENCHMARK) \
	PYTHONPATH=PYTHONPATH:./build/disk-manager-exercise-client \
	nose2 --config=subsystem/nose2.cfg --verbose -F --project-directory . $(TEST)

ifdef BENCHMARK
	# Create a benchmark report
	ARTIFACT_DIR=reports/benchmark \
	report subsystem/logs/stats/diskmanagerexercise_disk-manager-exercise_1.json
endif

package: image build
	packager pack artifacts.yaml --auto-push

generate_upgrade_manifest:
	python -m packaging_tools.upgrade_manifest_generator artifacts.yaml --file build/update_manifest.json --dirty

deploy: package generate_upgrade_manifest
	upgrade -v $(IP) build/update_manifest.json
