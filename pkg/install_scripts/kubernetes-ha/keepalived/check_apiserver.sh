#!/bin/bash

# if check error then repeat check for 12 times, else exit
err=0
for k in $(seq 1 12)
do
    nc -zv localhost 6443
    check_code=$?
    if [[ $check_code == "0" ]]; then
        err=0
        break
    else
        err=$(expr $err + 1)
        sleep 5
        continue
    fi
done

if [[ $err != "0" ]]; then
    # if apiserver is down send SIG=1
    echo '[ERROR] apiserver error'
    exit 1
else
    # if apiserver is up send SIG=0
    echo '[INFO] apiserver ok'
    exit 0
fi
