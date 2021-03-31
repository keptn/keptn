#!/bin/bash

function spin {
  sleep 0   # For deferring
  message=$1
  length=${#message}
  blocks=$(seq -s '\b' $((length+3))|tr -d '[:digit:]')
  i=1
  sp="🌑🌒🌓🌔🌕🌖🌗🌘"
  echo -n ' '
  while true
  do
    printf "${blocks}%s" "${sp:i++%${#sp}:1} ${message}"
    sleep 0.1
  done
}

function cleanup {
  kill -9 "$_main_pid" 2>/dev/null
  kill -9 "$_sp_pid" 2>/dev/null
}

message=$1
command=$2

trap "exit" INT TERM ERR
trap cleanup EXIT

printf "\n"
eval "$command" 1>/dev/null &
export _main_pid=$!
spin "$message" &
export _sp_pid=$!

wait $_main_pid
