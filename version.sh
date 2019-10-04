#!/bin/bash

appDir=$(dirname $0)
cd "$appDir"

START="egrep -ih 'app(\.){0,1}Version.=' *.go | egrep '\d+\.\d+\.\d+' | tail -n1"

if [ "$1" = -d ] ; then
   eval $START
else
   appVersion=$(eval $START | sed 's/^.*=//' | sed 's/\"//g' | sed 's,//.*,,' | sed 's/[[:space:]]//g')
   echo $appVersion
fi
