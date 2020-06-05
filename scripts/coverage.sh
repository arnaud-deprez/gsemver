#!/usr/bin/env bash
set -euo pipefail

covermode=${COVERMODE:-atomic}
coverdir=build/coverage
profile="${coverdir}/cover.out"

mkdir -p $coverdir

go test -race -coverprofile="${profile}" -covermode="$covermode" ./...

go tool cover -func "${profile}"
go tool cover -html "${profile}"