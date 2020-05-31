#!/bin/bash

echo "[$(date)] Waiting for Redis..." \
&& wait-for-it $REDIS_ADDRESS -t 0 -s \
&& /bin/bash -c "$@"
