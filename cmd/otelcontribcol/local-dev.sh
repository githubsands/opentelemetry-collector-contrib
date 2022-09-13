#!/bin/zsh

if [[ -z $1 ]];
then 
    echo "must set a first option ie /local-dev.sh receiver"
fi

RECEIVER=${1}
BUILDPATH=$(pwd)

rm -rf opentelemetry-collector-contrib
go work init opentelemetry-collector-contrib

(cd ../../.. && cp -R opentelemetry-collector-contrib ${BUILDPATH})
docker build -t collector-${RECEIVER} -f Dockerfile.local .
