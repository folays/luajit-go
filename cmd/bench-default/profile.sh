#!/bin/sh
set -ex

GODEBUG=cgocheck=0 LUAJIT_PROFILE= SILENCE= go run -gcflags=all='-N -l' -tags profile -v ./ -- "$@"

go tool pprof -http :8080 /tmp/luajit.pprof
