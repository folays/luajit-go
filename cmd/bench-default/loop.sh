#!/bin/sh

color=41

while true; do
  go run ../test/ 2>&1 | cat > /tmp/OUT

  clear
  #sleep 0.05

  cat /tmp/OUT
  echo ; echo ; echo

  printf "\e[%dm%s\e[49m\n" $color "$(date)"
  color=$((41 + ($color - 40 + 1) % 5))

  sleep 1
done
