#!/bin/bash

server_name="server"
network_name="tp0_testing_net"
port_default=12345
message_default="Testing server using netcat"

# If first parameter (port) is empty or "-", use port_default
if [ -z "$1" ] || [ "$1" == "-" ]; then
  port=$port_default
else
  port=$1
fi

message=${2:-$message_default}

result=$(docker run --rm --network $network_name busybox:latest sh -c "echo $message | nc -w 1 $server_name $port")

if [ "$result" == "$message" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
