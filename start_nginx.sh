#!/bin/bash

scriptPath=$(realpath $0)
scriptDir=$(dirname $scriptPath)

docker run \
  --rm \
  -it \
  --network dn42 \
  --name pagerank-api \
  -v $scriptDir/nginx/nginx.conf:/etc/nginx/nginx.conf \
  -v $scriptDir/nginx/conf.d:/etc/nginx/conf.d \
  -v $scriptDir/data:/usr/share/nginx/html \
  nginx:1.27.3
