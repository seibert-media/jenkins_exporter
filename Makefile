EXECUTABLE ?= jenkins_exporter
IMAGE ?= exporters/jenkins

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst master,latest,$(DRONE_BRANCH))
	else
		VERSION ?= latest
	endif
endif

ifndef SHA
	SHA := $(shell git rev-parse --short HEAD)
endif

ifndef DATE
	DATE := $(shell date -u '+%FT%T%z')
endif

LDFLAGS += -X "main.Version=$(VERSION)"
LDFLAGS += -X "main.Revision=$(SHA)"
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES := $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /_tools/)
SOURCES := $(shell find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./_tools/*")

.PHONY: all
all: build

.PHONY: update
update:
	retool do dep ensure -update

.PHONY: sync
sync:
	retool do dep ensure

.PHONY: graph
graph:
	mkdir -p docs/
	retool do dep status -dot | dot -T png -o docs/deps.png

.PHONY: clean
clean:
	go clean -i ./...
	rm -rf dist/

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: megacheck
megacheck:
	retool do megacheck $(PACKAGES)

.PHONY: lint
lint:
	for PKG in $(PACKAGES); do retool do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	retool do goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	CGO_ENABLED=0 go build -i -v -ldflags '-w $(LDFLAGS)'

.PHONY: install
install: $(SOURCES)
	CGO_ENABLED=0 go install -v -ldflags '-w $(LDFLAGS)'

.PHONY: release
release:
	CGO_ENABLED=0 retool do gox -verbose -ldflags '-w $(LDFLAGS)' -output="dist/$(EXECUTABLE)-$(VERSION)-{{.OS}}-{{.Arch}}"

HAS_RETOOL := $(shell command -v retool)

.PHONY: retool
retool:
ifndef HAS_RETOOL
	go get -u github.com/twitchtv/retool
endif
	retool sync
	retool build
