APP=protodep
BASE_PACKAGE=github.com/stormcat24/$(APP)
SERIAL_PACKAGES= \
		 cmd \
		 dependency \
		 helper \
		 repository \
		 service
TARGET_SERIAL_PACKAGES=$(addprefix test-,$(SERIAL_PACKAGES))

deps-build:
		go get -u github.com/golang/dep/cmd/dep
		go get github.com/golang/lint/golint

deps: deps-build
		dep ensure

deps-update: deps-build
		rm -rf ./vendor
		rm -rf Gopkg.lock
		dep ensure -update

define build-artifact
		GOOS=$(1) GOARCH=$(2) go build -o artifacts/$(APP)
		cd artifacts && tar cvzf $(APP)_$(1)_$(2).tar.gz $(APP)
		rm ./artifacts/$(APP)
		@echo [INFO]build success: $(1)_$(2)
endef

build-all:
		rm -rf ./artifacts/*
		$(call build-artifact,linux,amd64)
		$(call build-artifact,darwin,amd64)

build:
		go build -ldflags="-w -s" -o bin/protodep main.go

test: $(TARGET_SERIAL_PACKAGES)

$(TARGET_SERIAL_PACKAGES): test-%:
		go test $(BASE_PACKAGE)/$(*)

mock:
	go get github.com/golang/mock/mockgen
	mockgen -source helper/auth.go -package helper -destination helper/auth_mock.go

release:
	./release.sh

release-npm:
	npm publish
