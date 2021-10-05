#!/usr/bin/env bash

# -d as debug option to show if the correct line is identified
# without -d, the line is reduced to the pure version number.

appDir=$(dirname $0)
cd "$appDir"

START="egrep -i$_egrepFlag 'app(\.)?Version[[:space:]]*=' *.go | egrep '[0-9]+\.[0-9]+\.[0-9]+' | tail -n1"

if [ "$1" = -d ] ; then
   _egrepFlag=-H
   eval $START
else
   _egrepFlag=-h
   appVersion=$(eval $START | sed 's/^.*=//' | sed 's/\"//g' | sed 's,//.*,,' | sed 's/[[:space:]]//g' | sed 's/^-//')
   [ -z "$appVersion" ] && echo could not determine app version 1>&2 && exit 1
   echo $appVersion
fi
