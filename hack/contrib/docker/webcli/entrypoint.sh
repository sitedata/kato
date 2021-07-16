#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /usr/bin/kato-webcli version
else
    exec /usr/bin/kato-webcli $@
fi