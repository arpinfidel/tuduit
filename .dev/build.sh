#!/bin/bash
buildArgs=('build' '-v' '-gcflags' 'all=-N -l' '-o');
go "${buildArgs[@]}" .dev/tuduit ./cmd