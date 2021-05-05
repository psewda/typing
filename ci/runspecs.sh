#!/bin/bash

set -e

# cleanup the stuff here
cleanup() {
  docker stop typing > /dev/null
  docker rm typing > /dev/null
}

# check if google client id exported
if [ -z "$GOOGLE_CLIENT_ID" ]; then
  echo "error: google client id is NOT exported"
  exit 1
fi

# check if google client secret exported
if [ -z "$GOOGLE_CLIENT_SECRET" ]; then
  echo "error: google client secret is NOT exported"
  exit 1
fi

# check if refresh token exported
if [ -z "$REFRESH_TOKEN" ]; then
  echo "error: refresh token is NOT exported"
  exit 1
fi

# build the docker image
docker build -t typing:latest .

# update the google client cred file
mkdir -p /tmp/typing
cp ./ci/google_client_cred.json /tmp/typing/google_client_cred.json
sed -i "s/{client_id}/$GOOGLE_CLIENT_ID/" /tmp/typing/google_client_cred.json
sed -i "s/{client_secret}/${GOOGLE_CLIENT_SECRET}/" /tmp/typing/google_client_cred.json

# run the docker image in container
docker run \
  --detach \
  --name=typing \
  --publish 127.0.0.1:7070:7070 \
  --volume /tmp/typing/google_client_cred.json:/etc/typing/google_client_cred.json \
  --env REFRESH_TOKEN=$REFRESH_TOKEN \
  typing:latest
docker ps

# run specs
trap cleanup EXIT
npm test
