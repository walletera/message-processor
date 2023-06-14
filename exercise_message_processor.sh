#!/bin/bash

print_message() {
  blue='\033[0;34m'
  nc='\033[0m' # No Color
  echo -e "${blue}$1${nc}"
}

print_message "Exercising message processor"
print_message "----------------------------"

print_message "Starting rabbitmq"
docker-compose up -d

rabbitmq_is_ready() {
  docker-compose exec rabbitmq rabbitmqadmin list queues > /dev/null
}

until rabbitmq_is_ready; do
  sleep 1
  echo "Waiting rabbitmq to start..."
done

print_message "Starting message processor"
go build -o message_processor
./message_processor &> message_processor_output.txt &
MESSAGE_PROCESSOR_PID=$!

print_message "Sending messages to rabbitmq"
for i in {1..10} ; do
    message="Message #$i"
    echo "Publishing message $message"
    docker-compose exec rabbitmq rabbitmqadmin publish routing_key="message-consumer-queue" payload="$message"
done

print_message "Terminating message processor"
kill -n 15 $MESSAGE_PROCESSOR_PID

print_message "Shutting down rabbitmq"
docker-compose down

print_message "Printing message processor output"
cat message_processor_output.txt
