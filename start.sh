#!/bin/bash

scriptPath=$(realpath $0)
scriptDir=$(dirname $scriptPath)

mkdir -p $scriptDir/data

docker run \
  --rm \
  -it \
  --network dn42 \
  --dns fdda:8ca4:1556:3000:ffff::b \
  --name pagerank \
  -v $scriptDir:/root/work \
  -v $scriptDir/data:/root/work/data \
  --entrypoint /root/work/entrypoint.sh \
  -w /root/work \
    test:latest
