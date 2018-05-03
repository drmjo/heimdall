#!/bin/base

docker run -it --rm \
  -v `pwd`:/go/src/github.com/drmjo/heimdall \
  -v `pwd`/definitions:/go/definitions \
  -p 8844:80 \
  golang:1.8 bash
