#!/bin/bash
ps -A | grep '[b]uild.sh' | awk '{print $1}' | xargs kill -9 $1;
ps -A | grep '[o]gma-.*'  | awk '{print $1}' | xargs kill -9 $1
ps -A | grep '[d]lv' | awk '{print $1}' | xargs kill -9 $1
source ~/.bashrc;
echo 'use debug =' $USE_DEBUG;
if [ "$USE_DEBUG" = "true" ]; then
	echo 'debug mode';
	parallel --linebuffer --tagstring [{1}] --colsep ' ' -a ./.dev/run.txt /go/bin/dlv --continue --listen=:{2} --headless=true --api-version=2 --accept-multiclient exec {3} -- {4} {5} {6} {7};
else
	echo 'normal mode';
	parallel --linebuffer --tagstring [{1}] --colsep ' ' -a ./.dev/run.txt {3} {4} {5} {6} {7};
fi
