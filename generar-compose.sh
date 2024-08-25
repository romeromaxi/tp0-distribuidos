#!/bin/bash
message_usage="Usage: $0 <output_filename> <clients_number>"

if [ "$#" -ne 2 ]; then
    echo $message_usage
    exit 1
fi

output_filename=$1
clients_number=$2

regex='^[0-9]+$'
if ! [[ $clients_number =~ $regex ]] ; then
    echo "Error: <clients_number> must be a number"
    echo $message_usage
    exit 1
fi

compose_name="tp0"

server_name="server"
server_image="server:latest"
server_entrypoint="python3 /main.py"
server_log_level="DEBUG"

client_image="client:latest"
client_entrypoint="/client"
client_log_level="DEBUG"

network_name="testing_net"
network_subnet="172.25.125.0/24"

echo "Generating configuration file $output_filename with $clients_number client(s)..."

cat <<EOL > $output_filename
name: $compose_name
services:
  $server_name:
    container_name: $server_name
    image: $server_image
    entrypoint: $server_entrypoint
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=$server_log_level
    networks:
      - $network_name

EOL

for i in $(seq 1 $clients_number)
do
cat <<EOL >> $output_filename
  client$i:
    container_name: client$i
    image: $client_image
    entrypoint: $client_entrypoint
    environment:
      - CLI_ID=$i
      - CLI_LOG_LEVEL=$client_log_level
    networks:
      - $network_name
    depends_on:
      - $server_name

EOL
done

cat <<EOL >> $output_filename
networks:
  $network_name:
    ipam:
      driver: default
      config:
        - subnet: $network_subnet
EOL

echo "The file $output_filename with $clients_number client(s) was generated successfully."