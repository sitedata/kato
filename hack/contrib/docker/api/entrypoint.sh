#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /run/kato-api version
else
    exec /run/kato-api --start=true $@
fi