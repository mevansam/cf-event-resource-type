#!/bin/bash

ROOT_DIR=$(cd $(dirname $0) && pwd)
pushd $ROOT_DIR
set -x

go get -u github.com/kardianos/govendor
govendor sync
govendor test +local

cd $ROOT_DIR/cmd/check
GOOS=linux GOARCH=amd64 go build

cd $ROOT_DIR/cmd/in
GOOS=linux GOARCH=amd64 go build

cd $ROOT_DIR/cmd/out
GOOS=linux GOARCH=amd64 go build

TAG="$(git tag -l --points-at HEAD)"
if [[ "$1" == "release" ]] && [[ -n "$TAG" ]] ; then
    docker build . -t cf-event --squash
    docker tag cf-event $2cf-event/$TAG
    # docker push $2cf-event
fi

set +x
popd
