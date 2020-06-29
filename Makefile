# Makefile to create go binaries for 3 target platforms
# License type: MIT
# Copyright (c) 2018 by Christian Engel (engel-ch@outlook.com)
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
# version: 3.2.1
# 3.2.1
# - change sequence of tasks: first: patch
# 3.1.1
# - new targets majord minord patchd to allow for allow-dirty
# 3.1.0
# - bumpversion support
# - targets: major minor patch
# 3.0.0
# - semantic versioning
# - binaries created under build/release/<architecture>
# - separate targets: osx, amd64, arm, armel
# 2.5
# - version detection of software redirected towards a version.sh script
# 2.3
# - appVersion in SRC_VERSION detection
# - glide.lock deleted by distclean target
# 2.2
# - vendor directory removed in distclean target

# binary name is the name of this directory, unless the name of this
# directory is src. In this case, use the parent directory name of src
# as the name for the binary
APP 			:= $(shell pwd | xargs basename)
ifeq ($(APP), src)
	APP := $(shell pwd | xargs dirname | xargs basename)
endif
SRC 			:= $(wildcard *.go)
# if you do not want to have version numbers in the binaries, just say
SRC_VERSION=-$(shell bash version.sh)

DEPS			:= $(SRC) .touchfile
SRCIMPORT 		:= $(foreach file,$(SRC),../../../$(file))
#OS 				:= $(shell uname | tr "A-Z" "a-z")

ifeq ($(ENV), release)
	ENVIRONMENT = release
else
	ENVIRONMENT = debug
endif
BINDIR		= build/${ENVIRONMENT}/
GOOS1 		:= darwin
GOARCH1 		:= amd64
GOARM1		:=
GOOS2 		:= linux
GOARCH2 		:= arm
GOARM2		:=
GOOS3 		:= linux
GOARCH3 		:= amd64
GOARM3		:=
GOOS4 		:= linux
GOARCH4 		:= arm
GOARM4	   := 5

# ARG			:= $(GOOS1)-$(GOARCH1)/$(APP)$(SRC_VERSION))

# targets which are not a file
.PHONY: clean test test.html test-html ndebug nodebug upx touch force deploy version chkVersion major minor patch

patch:
	bumpversion patch
	make all

patchd:
	bumpversion --allow-dirty patch
	make all

major:
	bumpversion major
	make all

majord:
	bumpversion --allow-dirty major
	make all

minor:
	bumpversion minor
	make all

minord:
	bumpversion --allow-dirty minor
	make all


# pl4 currently not called
all: chkVersion pl1 pl2 pl3 pl4 deploy

osx: chkVersion pl1

osx-deploy: osx deploy

osx-release:
	make LDFLAGS="-ldflags -w" osx

arm: chkVersion pl2

amd64: chkVersion pl3
	@echo environment is $(ENVIRONMENT)

amd64-deploy: amd64 deploy

amd64-release:
	make LDFLAGS="-ldflags -w" ENV="release" amd64
	[ ! -f build/release/$(GOOS3)-$(GOARCH3)$(GOARM2)/$(APP)$(SRC_VERSION).upx ] && \
	[ -f build/release/$(GOOS3)-$(GOARCH3)$(GOARM2)/$(APP)$(SRC_VERSION) ] && \
	cd build/release/$(GOOS3)-$(GOARCH3)$(GOARM2) && cp $(APP)$(SRC_VERSION) $(APP)$(SRC_VERSION).upx && upx $(APP)$(SRC_VERSION).upx && ln -fs $(APP)$(SRC_VERSION).upx $(APP).upx


armel: chkVersion pl4

version:
	@if [ ! -f version.sh ] ; then  echo version.sh not found ; exit 1 ; fi
	@echo $(SRC_VERSION) | sed 's/-//'

chkVersion:
ifeq ("$(SRC_VERSION)", "-")
	@echo Version of source code could not be detected. Please adapt the SRC_VERSION definition of the Makefile
	@echo or the source-code accordingly.
	@exit 1
else
	@echo  version detected to be $(SRC_VERSION) | sed 's/-//'
endif

# too call platforms independently from each other by a simple name plX

pl1: $(BINDIR)$(GOOS1)-$(GOARCH1)$(GOARM1)/$(APP)$(SRC_VERSION)

pl2: $(BINDIR)$(GOOS2)-$(GOARCH2)$(GOARM2)/$(APP)$(SRC_VERSION)

pl3: $(BINDIR)$(GOOS3)-$(GOARCH3)$(GOARM3)/$(APP)$(SRC_VERSION)

pl4: $(BINDIR)$(GOOS4)-$(GOARCH4)$(GOARM4)/$(APP)$(SRC_VERSION)

deploy:
	@test -f deploy.sh && echo deploying… && bash deploy.sh $(APP)$(SRC_VERSION); exit 0

# compile without debugging information (upx still has a positive effect on it.)

nodebug: ndebug

ndebug:
	touch .touchfile
	make LDFLAGS="-ldflags -w" all

force:
	touch .touchfile
	make all

upx:
	make nodebug all
	# test -d $(GOOS1)-$(GOARCH1) && (cd $(GOOS1)-$(GOARCH1) ; upx $(APP)$(SRC_VERSION) ; ln -s $(APP)$(SRC_VERSION) $(APP)$(SRC_VERSION).upx)
	@echo upx not working for OSX, skipping...
	test -d $(BINDIR)$(GOOS3)-$(GOARCH3)$(GOARM2) && (cd $(BINDIR)$(GOOS3)-$(GOARCH3)$(GOARM2) ; upx $(APP)$(SRC_VERSION) ; ln -s $(APP)$(SRC_VERSION) $(APP)$(SRC_VERSION).upx)
	test -d $(BINDIR)$(GOOS2)-$(GOARCH2)$(GOARM3) && (cd $(BINDIR)$(GOOS2)-$(GOARCH2)$(GOARM3) ; upx $(APP)$(SRC_VERSION) ; ln -s $(APP)$(SRC_VERSION) $(APP)$(SRC_VERSION).upx)
	test -d $(BINDIR)$(GOOS4)-$(GOARCH4)$(GOARM4) && (cd $(BINDIR)$(GOOS4)-$(GOARCH4)$(GOARM4) ; upx $(APP)$(SRC_VERSION) ; ln -s $(APP)$(SRC_VERSION) $(APP)$(SRC_VERSION).upx)
	make deploy

$(BINDIR)$(GOOS1)-$(GOARCH1)$(GOARM1)/$(APP)$(SRC_VERSION): $(DEPS)
	test -d  $(dir $@) || mkdir -p $(dir $@)
	@if [ -f $(dir $@)/$(APP)$(SRC_VERSION) ] ; then  echo version $(dir $@)/$(APP)$(SRC_VERSION) already existing ; exit 1 ; fi
ifeq ($(GOARM1), "")
	cd $(dir $@) ; env GOARCH=$(GOARCH1) GOOS=$(GOOS1) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
else
	cd $(dir $@) ; env GOARCH=$(GOARCH1) GOOS=$(GOOS1) GOARM=$(GOARM1) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
endif

$(BINDIR)$(GOOS2)-$(GOARCH2)$(GOARM2)/$(APP)$(SRC_VERSION): $(DEPS)
	test -d  $(dir $@) || mkdir -p $(dir $@)
	@if [ -f $(dir $@)/$(APP)$(SRC_VERSION) ] ; then  echo version $(dir $@)/$(APP)$(SRC_VERSION) already existing ; exit 1 ; fi
ifeq ($(GOARM2), "")
	cd $(dir $@) ; env GOARM=7 GOARCH=$(GOARCH2) GOOS=$(GOOS2) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
else
	cd $(dir $@) ; env GOARM=7 GOARCH=$(GOARCH2) GOOS=$(GOOS2) GOARM=$(GOARM2) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
endif

$(BINDIR)$(GOOS3)-$(GOARCH3)$(GOARM3)/$(APP)$(SRC_VERSION): $(DEPS)
	test -d  $(dir $@) || mkdir -p $(dir $@)
	@if [ -f $(dir $@)/$(APP)$(SRC_VERSION) ] ; then  echo version $(dir $@)/$(APP)$(SRC_VERSION) already existing ; exit 1 ; fi
ifeq ($(GOARM3), "")
	cd $(dir $@) ; env GOARCH=$(GOARCH3) GOOS=$(GOOS3) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
else
	cd $(dir $@) ; env GOARCH=$(GOARCH3) GOOS=$(GOOS3) GOARM=$(GOARM3) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
endif

$(BINDIR)$(GOOS4)-$(GOARCH4)$(GOARM4)/$(APP)$(SRC_VERSION): $(DEPS)
	@test -d  $(dir $@) || mkdir -p $(dir $@) && echo directory $(dir $@) created
	@if [ -f $(dir $@)/$(APP)$(SRC_VERSION) ] ; then  echo version $(dir $@)/$(APP)$(SRC_VERSION) already existing ; exit 1 ; fi
ifeq ($(GOARM4), "")
	cd $(dir $@) ; env GOARCH=$(GOARCH4) GOOS=$(GOOS4) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
else
	cd $(dir $@) ; GOOS=$(GOOS4) GOARCH=$(GOARCH4) GOARM=$(GOARM4) go build $(LDFLAGS) -o $(APP)$(SRC_VERSION) $(SRCIMPORT) && ln -fs $(APP)$(SRC_VERSION) $(APP)
endif

clean:
	rm -fr build coverage.out

distclean: clean
	find . -name '*.~' -exec /bin/rm -f {} \;
	find . -name '*~' -exec /bin/rm -f {} \;
	find . -name '*.bak' -exec /bin/rm -f {} \;
	#/bin/rm -f glide.lock

# delete the latest binary only
lastclean:
	for dir in $(BINDIR)$(GOOS1)-$(GOARCH1)$(GOARM1) $(BINDIR)$(GOOS2)-$(GOARCH2)$(GOARM2) $(BINDIR)$(GOOS3)-$(GOARCH3)$(GOARM3) $(BINDIR)$(GOOS4)-$(GOARCH4)$(GOARM4) ; do /bin/rm -f $$dir/$(APP)$(SRC_VERSION)* ; done

test:
	go test -v -coverprofile=coverage.out

test-html: test.html

test.html:
	go tool cover -html=coverage.out

.touchfile:
	touch .touchfile
