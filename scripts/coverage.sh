#!/usr/bin/env bash

set -euo pipefail

covermode=${COVERMODE:-atomic}
coverdir=build/coverage
profile="${coverdir}/cover.out"

hash goveralls 2>/dev/null || go get github.com/mattn/goveralls
hash godir 2>/dev/null || go get github.com/Masterminds/godir

generate_cover_data() {
  for d in $(godir) ; do
    (
      local output="${coverdir}/${d//\//-}.cover"
      echo "$coverdir -> $output"
      go test -race -coverprofile="${output}" -covermode="$covermode" "$d"
    )
  done

  echo "mode: $covermode" >"$profile"
  grep -h -v "^mode:" "$coverdir"/*.cover >>"$profile"
}

push_to_coveralls() {
  goveralls -coverprofile="${profile}" -service=travis-ci
}

push_to_codecov() {
  bash <(curl -s https://codecov.io/bash)
}

mkdir -p $coverdir

generate_cover_data
go tool cover -func "${profile}"

case "${1-}" in
  --html)
    go tool cover -html "${profile}"
    ;;
  --codecov)
    push_to_codecov
    ;;
  --coveralls)
    push_to_coveralls
    ;;
esac