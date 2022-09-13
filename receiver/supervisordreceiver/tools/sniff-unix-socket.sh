#!/bin/zsh

PIPOED_SOCKET=/var/run/supervisor.sock
SOCKET=/var/run/supervisor.sock.temp

mv $SOCKET $PIPED_SOCKET
socat -t100 -x -v UNIX-LISTEN:$PIPED_SOCKET,mode=777,reuseaddr,fork UNIX-CONNECT:$SOCKET
