# Make an artifacts directory
# clear the artifacts directory before building
if [ -d ./artifacts ]; then
    rm -rf ./artifacts
fi
mkdir -p artifacts

# run the build for each supported OS
for os in "linux" "darwin" "windows"; do
    echo "Building for $os..."
    GOOS=$os CGO_ENABLED=0 go build -a -ldflags="-X main.version=${{ steps.get_tag.outputs.SOURCE_TAG }}" -o ./artifacts/aws-secret-reader_${os} .
    
    # If its windows we need to rename it to have .exe at the end.
    if [ $os == "windows" ]; then
        mv ./artifacts/aws-secret-reader_$os ./artifacts/aws-secret-reader_$os.exe
    fi
done
# Make an Arm bin for linux also
for arch in arm64 arm; do
    GOOS=linux GOARCH=$arch CGO_ENABLED=0 go build -a -ldflags="-X main.version=${{ steps.get_tag.outputs.SOURCE_TAG }}" -o ./artifacts/aws-secret-reader_${os}_$arch .
done