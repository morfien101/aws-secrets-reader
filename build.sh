#!/bin/bash

VERSION="v0.0.3"

# clear the artifacts directory before building
if [ -d ./artifacts ]; then
    rm -rf ./artifacts
fi

mkdir artifacts

# test the code. Failing tests means no build. Soz, fix the tests...
if go test ./...; then
    # build for all the major OSs
    for os in "linux" "darwin" "windows"; do
        echo "Building for $os..."
        GOOS=$os CGO_ENABLED=0 go build -a -ldflags="-X 'main.version=${VERSION}'" -o ./artifacts/aws-secret-reader_$os .
        
        # If its windows we need to rename it to have .exe at the end.
        if [ $os == "windows" ]; then
            mv ./artifacts/aws-secret-reader_$os ./artifacts/aws-secret-reader_$os.exe
        fi
    done
fi
