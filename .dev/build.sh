#!/bin/bash
buildArgs=('build' '-v' '-gcflags' 'all=-N -l' '-o');
go "${buildArgs[@]}" .dev/tmp/http ./cmd/http
