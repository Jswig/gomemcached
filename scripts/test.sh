#!/bin/sh
set -e

repo_root=$(git rev-parse --show-toplevel)

echo "Running tests"
go test "${repo_root}/..."