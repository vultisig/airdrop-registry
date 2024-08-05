#!/bin/sh

if [ "$1" = "server" ]; then
    exec ./airdrop-server
elif [ "$1" = "worker" ]; then
    exec ./airdrop-worker
else
    echo "Please specify either 'server' or 'worker'"
    exit 1
fi
