#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /run/kato-grctl version
elif [ "$1" = "copy" ];then
    cp -a /run/kato-grctl /rootfs/usr/local/bin/
else
    exec /run/kato-grctl "$@"
fi