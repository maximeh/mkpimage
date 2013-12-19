#!/bin/bash

if [ "$UID" -ne "0" ] ; then
  echo "You must be root to do that!"
  exit 1
fi

export GOPATH=$(pwd)
cd $(pwd)
go get github.com/docopt/docopt.go &> /dev/null
go build
install mkpimage /usr/local/bin/
