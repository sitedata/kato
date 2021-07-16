#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /run/kato-chaos version
else
    exec /run/kato-chaos $@
fi