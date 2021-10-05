# Makefile to create go binaries for 3 target platforms
# License type: MIT
# Copyright (c) 2020 by Christian Engel (engel-ch@outlook.com)
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
# associated documentation files (the "Software"), to deal in the Software without restriction,
# including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all copies or substantial
# portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
# LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
# WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
# SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# 5.1.0
# - bumpversion simplified and call corrected for new version of bumpversion (implying main.go)
# - test and test-coverage statements included
#
# 5.0.0
# - merged w/ version from gomakefile project

# Sub-Makefiles have to implement the targets
# 	all
# 	upx

# binary name is the name of this directory, unless the name of this
# directory is src. In this case, use the parent directory name of src
# as the name for the binary
APP 			:= $(shell pwd | xargs basename)
ifeq ($(APP), src)
	APP := $(shell pwd | xargs dirname | xargs basename)
endif

# SRC_VERSION contains the version number of the source code. If a binary with a version-number x.y.z already exists
# then the compilation would stop with an error. You would either have to increase the version-number or you would
# first have to remove the last version with `make lastclean`

SRC_VERSION=$(shell bash version.sh)

# Include all sub-Makefile for different platforms. The filename is of the form Makefile.xxxxx if
# it does not end in ~, .old, or .template.
PLATFORM_MAKEFILE := $(shell ls Makefile.* | grep -v ~$ | grep -v .old$ | grep -v .template$ )

# You can set an environment variable called `ENV` to define whether a debug or a release
# version of the software is to be created. The default version is debug.
# Further stripped down release versions can be created using the `upx` binary. 
# You can create such a version by calling `make upx`
#So far, `upx` worked well for all Linux platforms but the `OSX` alias Darwin version created
# binaries that always crashed immediately after start.
ifeq ($(ENV), release)
	LDFLAGS="-ldflags -w"
	ENVIRONMENT = release
else
	LDFLAGS=
	ENVIRONMENT = debug
endif

BINDIR	:=	./build/

all xxupx: versionChk go.mod
	@for i in $(PLATFORM_MAKEFILE) ; do \
		echo APP=$(APP) SRC_VERSION=$(SRC_VERSION) BINDIR=$(BINDIR)$(ENVIRONMENT)/ LDFLAGS=$(LDFLAGS) make -f $$i $@; \
		APP=$(APP) SRC_VERSION=$(SRC_VERSION) BINDIR=$(BINDIR)$(ENVIRONMENT)/ LDFLAGS=$(LDFLAGS) make -f $$i $@; \
	done

rel: release

release:
	ENV=release make all

release-upx upx:
	ENV=release make xxupx

versionChk:
	@if [ `echo $(SRC_VERSION) | egrep  "[0-9]+\.[0-9]+\.[0-9]+" | wc -c` -le 1 ] ; then echo no correct version detected ; exit 1; fi

# support multiple spellings as previous versions
test-html: test.html

# go modules support
# LOCAL PACKAGES are to be listed in GOPRIVATE; otherwise an error is returned
go.mod:
	GOPRIVATE="miniWebSrv/Utils" go mod init
	GOPRIVATE="miniWebSrv/Utils" go mod vendor

# ----------------------------- bumpversion ----------------------------------------------------


p: patch

patch:
	bumpversion --current-version `./version.sh` --allow-dirty patch main.go

ma: major

major:
	bumpversion --current-version `./version.sh` --allow-dirty major main.go

mi: minor

minor:
	bumpversion --current-version `./version.sh` --allow-dirty minor main.go

# ----------------------------- testing ----------------------------------------------------
# testing on special platform is something not yet implemented yet

# run unit-tests
t: test

test:
	go test -v

# determine test coverage
tc: test-coverage

test-coverage:
	go test -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

tch: test-coverage-html

test-coverage-html:
	go test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out
	
# ----------------------------- cleaning ----------------------------------------------------

# delete the latest binary (specified in main.go) only
# todo: does not fix the s-link in these target directory to the last binary
lastclean:
	@echo removing the following binaries:
	@find $(BINDIR) -name $(APP)-$(SRC_VERSION) -print -exec  /bin/rm -f {} \; | sed 's/^/   /'
	@echo S-links from $(APP) might now point into the void, but it is ok for recompiling.

clean:
	find . -name coverage.out -exec /bin/rm -f {} \;
	find . -name '*~' -exec /bin/rm -f {} \;
	find . -name '*.bak' -exec /bin/rm -f {} \;
	find . -name '*.bup' -exec /bin/rm -f {} \;

distclean: clean
	/bin/rm -fr $(BINDIR)

# EOF