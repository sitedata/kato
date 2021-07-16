#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /run/kato-init-probe version
else
    exec /run/kato-init-probe $@
fi