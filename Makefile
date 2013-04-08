#!/usr/bin/env make -f

PROGRAM := relpdump
VERSION := 0.0.1

tempdir        := $(shell mktemp -d)
controldir     := $(tempdir)/DEBIAN
installpath    := $(tempdir)/usr/bin
buildpath      := .build
buildpackpath  := $(buildpath)/pack
buildpackcache := $(buildpath)/cache

define DEB_CONTROL
Package: $(PROGRAM)
Version: $(VERSION)
Architecture: amd64
Maintainer: "Dan Peterson" <dpiddy@gmail.com>
Section: admin
Priority: optional
Description: Receive, acknowledge, and dump RELP messages.
endef
export DEB_CONTROL

deb: build
	mkdir -p -m 0755 $(controldir)
	echo "$$DEB_CONTROL" > $(controldir)/control
	mkdir -p $(installpath)
	install bin/$(PROGRAM) $(installpath)/$(PROGRAM)
	fakeroot dpkg-deb --build $(tempdir) .
	rm -rf $(tempdir)

clean:
	rm -rf $(buildpath)
	rm -f $(PROGRAM)*.deb

build: $(buildpackpath)/bin
	$(buildpackpath)/bin/compile . $(buildpackcache)

$(buildpackcache):
	mkdir -p $(buildpath)
	mkdir -p $(buildpackcache)
	curl -o $(buildpath)/go-git-only.tgz http://codon-buildpacks.s3.amazonaws.com/buildpacks/fabiokung/go-git-only.tgz

$(buildpackpath)/bin: $(buildpackcache)
	mkdir -p $(buildpackpath)
	tar -C $(buildpackpath) -zxf $(buildpath)/go-git-only.tgz
