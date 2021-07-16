#!/bin/bash
if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    /run/kato-mq version
else
    exec /run/kato-mq $@
fi