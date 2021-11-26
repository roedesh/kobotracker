#!/bin/sh

export DBUS_SESSION_BUS_ADDRESS NICKEL_HOME WIFI_MODULE LANG WIFI_MODULE_PATH INTERFACE
sync
killall -TERM nickel hindenburg sickel fickel 2>/dev/null
/mnt/onboard/.adds/cryptokobo/cryptokobo
exec /mnt/onboard/.adds/cryptokobo/nickel.sh