#!/usr/bin/env bash

# This script builds the executor binary.

cd "$(dirname "${BASH_SOURCE[0]}")"/../../..
set -eu

OUTPUT=$(mktemp -d -t sgbuild_XXXXXXX)
cleanup() {
  rm -rf "$OUTPUT"
}
trap cleanup EXIT

mkdir -p "${OUTPUT}/$(git rev-parse HEAD)/linux-amd64"

# Environment for building linux binaries
export GO111MODULE=on
export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export VERSION

echo "--- go build"
pushd ./enterprise/cmd/executor 1>/dev/null
pkg="github.com/sourcegraph/sourcegraph/enterprise/cmd/executor"
bin_name="$OUTPUT/$(git rev-parse HEAD)/linux-amd64/$(basename $pkg)"
go build -trimpath -ldflags "-X github.com/sourcegraph/sourcegraph/internal/version.version=$VERSION -X github.com/sourcegraph/sourcegraph/internal/version.timestamp=$(date +%s)" -buildmode exe -tags dist -o "$bin_name" "$pkg"
popd 1>/dev/null

echo "--- create binary artifacts"
pushd "${OUTPUT}/$(git rev-parse HEAD)" 1>/dev/null

echo "executor built from https://github.com/sourcegraph/sourcegraph" >info.txt
echo >>info.txt
git log -n1 >>info.txt
sha256sum linux-amd64/executor >>linux-amd64/executor_SHA256SUM
popd 1>/dev/null

# Upload the new release folder
echo "--- upload binary artifacts"
gsutil cp -r "${OUTPUT}" gs://sourcegraph-artifacts
gsutil iam ch allUsers:objectViewer gs://sourcegraph-artifacts
