#!/usr/bin/env bash
# creates a ton of images in a specific registry, best be run with GNU parallel:
# parallel -j 5 'bash populate_registry.sh {}' ::: {500..1600}

temp_dir=$(mktemp -d)

cp Dockerfile "$temp_dir"

cd "$temp_dir"

sed -i "s/5/$RANDOM/g" Dockerfile

docker build -t eu.gcr.io/faulty-crane-testing/faulty-crane-test:v$1 .
docker push eu.gcr.io/faulty-crane-testing/faulty-crane-test:v$1

# cleanup
docker rmi eu.gcr.io/faulty-crane-testing/faulty-crane-test:v$1
rm -rf "$temp_dir"
