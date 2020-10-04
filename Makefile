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
# 5.0.0
# - merged w/ version from gomakefile project

# Sub-Makefiles have to implement the targets
# 	all
# 	upx
# 	test
# 	test.html

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

all xxupx: go.mod
	@for i in $(PLATFORM_MAKEFILE) ; do \
		APP=$(APP) SRC_VERSION=$(SRC_VERSION) BINDIR=$(BINDIR)$(ENVIRONMENT)/ LDFLAGS=$(LDFLAGS) make -f $$i $@; \
	done

release:
	ENV=release make all

release-upx upx:
	ENV=release make xxupx

# support multiple spellings as previous versions
test-html: test.html

# go modules support
go.mod:
	go mod init
	go mod vendor

# ----------------------------- bumpversion ----------------------------------------------------

pd: patchd

patchd:
	bumpversion --allow-dirty patch

p: patch

patch:
	bumpversion patch

ma: major

major:
	bumpversion major

mad: majord

majord:
	bumpversion --allow-dirty major

mi: minor

minor:
	bumpversion minor

mid: minord

minord:
	bumpversion --allow-dirty minor

# ----------------------------- testing ----------------------------------------------------
# testing on special platform is something not yet implemented yet

test:
	go test -v -coverprofile=coverage.out

test.html: test
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