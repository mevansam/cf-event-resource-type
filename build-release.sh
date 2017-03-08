#!/bin/bash

function usage {
    echo -e "\nUSAGE: ./build-release.sh [release DOCKER_REPO DOCKER_USER DOCKER_PASSWORD]"
    echo -e "\n    release          Specifying \"release\" instructs the script to build the Docker image and push it to a remote repository"
    echo -e "    DOCKER_REPO      The name of the repository in Docker hub or address of private repository server"
    echo -e "    DOCKER_USER      User name to login as to the remote repository"
    echo -e "    DOCKER_PASSWORD  The users password\n"
    exit 1
}

case $1 in
help)
    usage
    ;;
release)
    [[ $# -eq 4 ]] || usage
    ;;
*)
    [[ $# -eq 0 ]] || usage
    ;;
esac

TAG="$(git tag -l --points-at HEAD)"
if [[ "$1" == "release" ]] && [[ -z "$TAG" ]] ; then
    echo "To build and push the release image their must be a version tag at the head of the branch."
    exit 1
fi

ROOT_DIR=$(cd $(dirname $0) && pwd)
pushd $ROOT_DIR
set -x

go get -u github.com/kardianos/govendor
govendor sync
govendor test +local

pushd $ROOT_DIR/cmd/check
GOOS=linux GOARCH=amd64 go build
popd

pushd $ROOT_DIR/cmd/in
GOOS=linux GOARCH=amd64 go build
popd

pushd $ROOT_DIR/cmd/out
GOOS=linux GOARCH=amd64 go build
popd

TAG="$(git tag -l --points-at HEAD)"
if [[ "$1" == "release" ]]; then
    docker build . -t cf-event-release --squash

    if [[ -n $3 ]] && [[ -n $4 ]]; then
        docker login -u $3 -p $4

        docker tag cf-event-release $2/cf-event:latest
        docker tag cf-event-release $2/cf-event:$TAG
        docker push $2/cf-event
    fi
fi

set +x
popd
