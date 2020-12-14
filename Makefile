VERSION := $(shell jq -r .version <package.json)

APP      := $(shell jq -r .name <package.json)
PACKAGES := $(shell go list -f {{.Dir}} ./...)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: clean release README.md

clean:
	rm -rf binaries/
	rm -rf release/

release-exists:
	hub release show "$(VERSION)"

release: zip
	go mod tidy
	git add go.mod go.sum vendor
	git add package.json
	git add Makefile
	git commit -am "Release $(VERSION)" || true
	git push
	hub release create $(VERSION) -m "$(VERSION)" -a release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_64.zip -a release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_32.zip -a release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -a release/$(APP)_$(VERSION)_linux_arm64.tar.gz

README.md:
	go get github.com/keilerkonzept/$(APP) && <README.template.md subst \
		EXAMPLES_MAIN_TF="$$(cat examples/main.tf)"\
		EXAMPLE_PRETTY="$$($(APP) check examples)"\
		EXAMPLE_LIST="$$($(APP) list -o json examples | jq .)"\
		EXAMPLE_LIST_PRETTY="$$($(APP) list examples)"\
		EXAMPLE_UPDATES="$$($(APP) check -o json examples | jq .)"\
		EXAMPLE_UPDATES_PRETTY="$$($(APP) check examples)"\
		EXAMPLE_UPDATES_ALL_PRETTY="$$($(APP) check -all examples)"\
		EXAMPLE_UPDATES_SINGLE="$$($(APP) check -o json -module=consul_github_https -module=consul_github_ssh examples | jq .)"\
		EXAMPLE_UPDATES_SINGLE_ALL_PRETTY="$$($(APP) check -all -module=consul_github_https -module=consul_github_ssh examples)"\
		EXAMPLE_UPDATES_SINGLE_PRETTY="$$($(APP) check -module=consul_github_https -module=consul_github_ssh examples)"\
		USAGE="$$($(APP) -h 2>&1)"\
		USAGE_LIST="$$($(APP) list -h 2>&1)"\
		USAGE_CHECK="$$($(APP) check -h 2>&1)"\
		VERSION="$(VERSION)" APP="$(APP)"> README.md

zip: release/$(APP)_$(VERSION)_osx_x86_64.tar.gz release/$(APP)_$(VERSION)_windows_x86_64.zip release/$(APP)_$(VERSION)_linux_x86_64.tar.gz release/$(APP)_$(VERSION)_windows_x86_32.zip release/$(APP)_$(VERSION)_linux_x86_32.tar.gz release/$(APP)_$(VERSION)_linux_arm64.tar.gz

binaries: binaries/osx_x86_64/$(APP) binaries/windows_x86_64/$(APP).exe binaries/linux_x86_64/$(APP) binaries/windows_x86_32/$(APP).exe binaries/linux_x86_32/$(APP)

release/$(APP)_$(VERSION)_osx_x86_64.tar.gz: binaries/osx_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -C binaries/osx_x86_64 $(APP)

binaries/osx_x86_64/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/osx_x86_64/$(APP) .

release/$(APP)_$(VERSION)_windows_x86_64.zip: binaries/windows_x86_64/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_64.zip $(APP).exe

binaries/windows_x86_64/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/windows_x86_64/$(APP).exe .

release/$(APP)_$(VERSION)_linux_x86_64.tar.gz: binaries/linux_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -C binaries/linux_x86_64 $(APP)

binaries/linux_x86_64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/linux_x86_64/$(APP) .

release/$(APP)_$(VERSION)_windows_x86_32.zip: binaries/windows_x86_32/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_32 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_32.zip $(APP).exe

binaries/windows_x86_32/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/windows_x86_32/$(APP).exe .

release/$(APP)_$(VERSION)_linux_x86_32.tar.gz: binaries/linux_x86_32/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -C binaries/linux_x86_32 $(APP)

binaries/linux_x86_32/$(APP): $(GOFILES)
	GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/linux_x86_32/$(APP) .

release/$(APP)_$(VERSION)_linux_arm64.tar.gz: binaries/linux_arm64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_arm64.tar.gz -C binaries/linux_arm64 $(APP)

binaries/linux_arm64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION) -X main.app=$(APP)" -o binaries/linux_arm64/$(APP) .
