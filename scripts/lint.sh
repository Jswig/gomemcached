#!/bin/sh
set -e

repo_root=$(git rev-parse --show-toplevel)

echo "Running go vet"
go vet "${repo_root}/..."

echo "Running staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck "${repo_root}/..."
