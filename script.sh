#!/bin/sh
export GO111MODULE=on
cd $TRAVIS_BUILD_DIR/story/cmd
go build 
