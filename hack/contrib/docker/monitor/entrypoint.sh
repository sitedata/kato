#!/bin/ash
if [ "$1" = "bash" ];then
    exec /bin/ash
elif [ "$1" = "version" ];then
    /run/kato-monitor version
else
    exec /run/kato-monitor $@
fi