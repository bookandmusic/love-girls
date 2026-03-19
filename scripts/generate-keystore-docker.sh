#!/bin/bash

docker run --rm -it \
    -v "$(pwd):/output" \
    -w /output \
    eclipse-temurin:17-jdk \
    bash scripts/generate-android-keystore.sh