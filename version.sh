#!/usr/bin/env bash

appDir=$(dirname $0)
cd "$appDir"

START="egrep -ih 'app(\.)?Version[[:space:]]*=' *.go | egrep '\d+\.\d+\.\d+' | tail -n1"

if [ "$1" = -d ] ; then
   eval $START
else
   appVersion=$(eval $START | sed 's/^.*=//' | sed 's/\"//g' | sed 's,//.*,,' | sed 's/[[:space:]]//g' | sed 's/^-//')
   [ -z "$appVersion" ] && echo could not determine app version && exit 1
   echo $appVersion
fi
